package services

import (
	"errors"
	"log"
	"strconv"

	"git.onebytedata.com/odb/go-dicom/dimsec"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/network/dicomstatus"
	"git.onebytedata.com/odb/go-dicom/tags"
	"git.onebytedata.com/odb/go-dicom/uid"
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
	err := d.openAssociation(pdu, uid.VerificationSOPClass.UID, timeout)
	if err != nil {
		return err
	}
	err = dimsec.CEchoWriteRQ(pdu, uid.VerificationSOPClass.UID)
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
	err := d.openAssociation(pdu, SOPClassUID.UID, timeout)
	if err != nil {
		return results, status, err
	}
	err = dimsec.CFindWriteRQ(pdu, Query, SOPClassUID.UID)
	if err != nil {
		return results, status, err
	}
	for (status != -1) && (status != 0) {
		ddo, s, err := dimsec.CFindReadRSP(pdu)
		status = s
		if err != nil {
			return results, status, err
		}
		if (status == dicomstatus.Pending) || (status == dicomstatus.PendingWithWarnings) {
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
	status := dicomstatus.Pending
	SOPClassUID := uid.StudyRootQueryRetrieveInformationModelMOVE

	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, SOPClassUID.UID, timeout)
	if err != nil {
		return -1, err
	}
	err = dimsec.CMoveWriteRQ(pdu, Query, SOPClassUID.UID, destAET)
	if err != nil {
		return -1, err
	}

	for status == dicomstatus.Pending {
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
		if r != dicomstatus.Success {
			return errors.New("ERROR, serviceuser::StoreSCU, dimsec.CStoreReadRSP failed")
		}
		c, err := dimsec.CStoreReadRSP(pdu)
		if err != nil {
			return err
		}
		if c != dicomstatus.Success {
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
	pdu.SetCallingAE(d.destination.CallingAE)
	pdu.SetCalledAE(d.destination.CalledAE)
	pdu.SetTimeout(timeout)

	network.Resetuniq()
	PresContext := network.NewPresentationContext()
	PresContext.SetAbstractSyntax(AbstractSyntax)
	PresContext.AddTransferSyntax(uid.ImplicitVRLittleEndian.UID)
	PresContext.AddTransferSyntax(uid.ExplicitVRLittleEndian.UID)
	pdu.AddPresContexts(PresContext)

	return pdu.Connect(d.destination.HostName, strconv.Itoa(d.destination.Port))
}

func (d *scu) writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (int, error) {
	status := -1

	PCID := pdu.GetPresentationContextID()
	if PCID == 0 {
		return -1, errors.New("ERROR, serviceuser::WriteStoreRQ, PCID==0")
	}
	TrnSyntOUT := pdu.GetTransferSyntax(PCID)

	if TrnSyntOUT == nil {
		return -1, errors.New("ERROR, serviceuser::WriteStoreRQ, TrnSyntOut is empty")
	}

	if TrnSyntOUT == DDO.GetTransferSyntax() {
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return status, err
		}
		status = dicomstatus.Success
	} else {
		DDO.SetTransferSyntax(TrnSyntOUT)
		DDO.SetExplicitVR(true)
		DDO.SetBigEndian(false)
		if TrnSyntOUT.UID == uid.ImplicitVRLittleEndian.UID {
			DDO.SetExplicitVR(false)
		}
		if TrnSyntOUT.UID == uid.ExplicitVRBigEndian.UID {
			DDO.SetBigEndian(true)
		}
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return -1, err
		}
		status = dicomstatus.Success
	}

	return status, nil
}
