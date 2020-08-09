package services

import (
	"errors"
	"log"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/dimsec"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/uid"
)

// SCU - inteface to a scu
type SCU interface {
	EchoSCU(timeout int) error
	FindSCU(Query media.DcmObj, timeout int) (int, int, error)
	MoveSCU(destAET string, Query media.DcmObj, timeout int) (int, error)
	StoreSCU(FileName string, timeout int) error
	SetOnCFindResult(f func(result media.DcmObj))
	SetOnCMoveResult(f func(result media.DcmObj))
	openAssociation(pdu network.PDUService, AbstractSyntax string, timeout int) error
	writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (int, error)
}

type scu struct {
	destination   *network.Destination
	onCFindResult func(result media.DcmObj)
	onCMoveResult func(result media.DcmObj)
}

// NewSCU - Creates an interface to scu
func NewSCU(destination *network.Destination) SCU {
	return &scu{
		destination: destination,
	}
}

func (d *scu) EchoSCU(timeout int) error {
	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, uid.VerificationSOPClass, timeout)
	if err != nil {
		return err
	}
	err = dimsec.CEchoWriteRQ(pdu, uid.VerificationSOPClass)
	if err != nil {
		return err
	}
	err = dimsec.CEchoReadRSP(pdu)
	if err != nil {
		return err
	}
	pdu.Close()
	return nil
}

func (d *scu) FindSCU(Query media.DcmObj, timeout int) (int, int, error) {
	results := 0
	status := 1
	SOPClassUID := uid.StudyRootQueryRetrieveInformationModelFIND

	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, SOPClassUID, timeout)
	if err != nil {
		return results, status, err
	}
	err = dimsec.CFindWriteRQ(pdu, Query, SOPClassUID)
	if err != nil {
		return results, status, err
	}
	for (status != -1) && (status != 0) {
		ddo, s, err := dimsec.CFindReadRSP(pdu)
		status = s
		if err != nil {
			return results, status, err
		}
		if (status == 0xFF00) || (status == 0xFF01) {
			results++
			if d.onCFindResult != nil {
				d.onCFindResult(ddo)
			} else {
				log.Println("No onCFindResult event found")
			}
		}
	}

	pdu.Close()
	return results, status, nil
}

func (d *scu) MoveSCU(destAET string, Query media.DcmObj, timeout int) (int, error) {
	var pending int
	status := 0xFF00
	SOPClassUID := uid.StudyRootQueryRetrieveInformationModelMOVE

	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, SOPClassUID, timeout)
	if err != nil {
		return -1, err
	}
	err = dimsec.CMoveWriteRQ(pdu, Query, SOPClassUID, destAET)
	if err != nil {
		return -1, err
	}

	for status == 0xFF00 {
		ddo, s, err := dimsec.CMoveReadRSP(pdu, &pending)
		status = s
		if err != nil {
			return -1, err
		}
		if d.onCMoveResult != nil {
			d.onCMoveResult(ddo)
		} else {
			log.Println("No onCMoveResult event found")
		}
	}

	pdu.Close()
	return status, nil
}

func (d *scu) StoreSCU(FileName string, timeout int) error {
	DDO, err := media.NewDCMObjFromFile(FileName)
	if err != nil {
		return err
	}

	SOPClassUID := DDO.GetString(tags.SOPClassUID)
	if len(SOPClassUID) > 0 {
		pdu := network.NewPDUService()
		err := d.openAssociation(pdu, SOPClassUID, timeout)
		if err != nil {
			return err
		}
		r, err := d.writeStoreRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return err
		}
		if r != 0x00 {
			return errors.New("ERROR, serviceuser::StoreSCU, dimsec.CStoreReadRSP failed")
		}
		c, err := dimsec.CStoreReadRSP(pdu)
		if err != nil {
			return err
		}
		if c != 0x00 {
			return errors.New("ERROR, serviceuser::StoreSCU, dimsec.CStoreReadRSP failed")
		}

		pdu.Close()
		return nil
	}
	return errors.New("ERROR, serviceuser::StoreSCU, OpenAssociation failed, RAET: " + d.destination.CalledAE)
}

func (d *scu) SetOnCFindResult(f func(result media.DcmObj)) {
	d.onCFindResult = f
}

func (d *scu) SetOnCMoveResult(f func(result media.DcmObj)) {
	d.onCMoveResult = f
}

func (d *scu) openAssociation(pdu network.PDUService, AbstractSyntax string, timeout int) error {
	pdu.SetRQCallingAE(d.destination.CallingAE)
	pdu.SetRQCalledAE(d.destination.CalledAE)
	pdu.SetTimeout(timeout)

	network.Resetuniq()
	PresContext := network.NewPresentationContext()
	PresContext.SetAbstractSyntax(AbstractSyntax)
	PresContext.AddTransferSyntax(uid.ImplicitVRLittleEndian)
	PresContext.AddTransferSyntax(uid.ExplicitVRLittleEndian)
	pdu.AddPresContexts(PresContext)

	return pdu.Connect(d.destination.HostName, strconv.Itoa(d.destination.Port))
}

func (d *scu) writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (int, error) {
	status := -1

	PCID := pdu.GetPresentationContextID()
	if PCID == 0 {
		return -1, errors.New("ERROR, serviceuser::WriteStoreRQ, PCID==0")
	}
	TrnSyntOUT := pdu.GetTransferSyntaxUID(PCID)

	if len(TrnSyntOUT) == 0 {
		return -1, errors.New("ERROR, serviceuser::WriteStoreRQ, TrnSyntOut is empty")
	}

	if TrnSyntOUT == DDO.GetTransferSyntax() {
		DDO.DumpTags()
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return status, err
		}
		status = 0
	} else {
		DDO.SetTransferSyntax(TrnSyntOUT)
		DDO.SetExplicitVR(true)
		DDO.SetBigEndian(false)
		if TrnSyntOUT == uid.ImplicitVRLittleEndian {
			DDO.SetExplicitVR(false)
		}
		if TrnSyntOUT == uid.ExplicitVRBigEndian {
			DDO.SetBigEndian(true)
		}
		DDO.DumpTags()
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return -1, err
		}
		status = 0
	}

	return status, nil
}
