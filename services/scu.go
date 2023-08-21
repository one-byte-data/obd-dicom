package services

import (
	"errors"
	"log"
	"strconv"

	"github.com/one-byte-data/obd-dicom/dictionary/sopclass"
	"github.com/one-byte-data/obd-dicom/dictionary/tags"
	"github.com/one-byte-data/obd-dicom/dictionary/transfersyntax"
	"github.com/one-byte-data/obd-dicom/dimsec"
	"github.com/one-byte-data/obd-dicom/media"
	"github.com/one-byte-data/obd-dicom/network"
	"github.com/one-byte-data/obd-dicom/network/dicomstatus"
)

// SCU - inteface to a scu
type SCU interface {
	EchoSCU(timeout int) error
	FindSCU(Query media.DcmObj, timeout int) (int, uint16, error)
	MoveSCU(destAET string, Query media.DcmObj, timeout int) (uint16, error)
	StoreSCU(FileName string, timeout int) error
	SetOnCFindResult(f func(result media.DcmObj))
	SetOnCMoveResult(f func(result media.DcmObj))
	openAssociation(pdu network.PDUService, AbstractSyntax string, timeout int) error
	writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (uint16, error)
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
	if err := d.openAssociation(pdu, sopclass.Verification.UID, timeout); err != nil {
		return err
	}
	if err := dimsec.CEchoWriteRQ(pdu, sopclass.Verification.UID); err != nil {
		return err
	}
	if err := dimsec.CEchoReadRSP(pdu); err != nil {
		return err
	}
	pdu.Close()
	return nil
}

func (d *scu) FindSCU(Query media.DcmObj, timeout int) (int, uint16, error) {
	results := 0
	status := dicomstatus.Warning
	SOPClassUID := sopclass.StudyRootQueryRetrieveInformationModelFind

	pdu := network.NewPDUService()
	if err := d.openAssociation(pdu, SOPClassUID.UID, timeout); err != nil {
		return results, status, err
	}
	if err := dimsec.CFindWriteRQ(pdu, Query, SOPClassUID.UID); err != nil {
		return results, status, err
	}
	for status != dicomstatus.Success {
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

func (d *scu) MoveSCU(destAET string, Query media.DcmObj, timeout int) (uint16, error) {
	var pending int
	status := dicomstatus.Pending
	SOPClassUID := sopclass.StudyRootQueryRetrieveInformationModelMove

	pdu := network.NewPDUService()
	if err := d.openAssociation(pdu, SOPClassUID.UID, timeout); err != nil {
		return dicomstatus.FailureUnableToProcess, err
	}
	if err := dimsec.CMoveWriteRQ(pdu, Query, SOPClassUID.UID, destAET); err != nil {
		return dicomstatus.FailureUnableToProcess, err
	}

	for status == dicomstatus.Pending {
		ddo, s, err := dimsec.CMoveReadRSP(pdu, &pending)
		status = s
		if err != nil {
			return dicomstatus.FailureUnableToProcess, err
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
	PresContext.AddTransferSyntax(transfersyntax.ImplicitVRLittleEndian.UID)
	pdu.AddPresContexts(PresContext)

	return pdu.Connect(d.destination.HostName, strconv.Itoa(d.destination.Port))
}

func (d *scu) writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (uint16, error) {
	status := dicomstatus.FailureUnableToProcess

	PCID := pdu.GetPresentationContextID()
	if PCID == 0 {
		return dicomstatus.FailureUnableToProcess, errors.New("ERROR, serviceuser::WriteStoreRQ, PCID==0")
	}
	TrnSyntOUT := pdu.GetTransferSyntax(PCID)

	if TrnSyntOUT == nil {
		return dicomstatus.FailureUnableToProcess, errors.New("ERROR, serviceuser::WriteStoreRQ, TrnSyntOut is empty")
	}

	if TrnSyntOUT == DDO.GetTransferSyntax() {
		if err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID); err != nil {
			return status, err
		}
		status = dicomstatus.Success
	} else {
		DDO.SetTransferSyntax(TrnSyntOUT)
		DDO.SetExplicitVR(true)
		DDO.SetBigEndian(false)
		if TrnSyntOUT.UID == transfersyntax.ImplicitVRLittleEndian.UID {
			DDO.SetExplicitVR(false)
		}
		if TrnSyntOUT.UID == transfersyntax.ExplicitVRBigEndian.UID {
			DDO.SetBigEndian(true)
		}
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return dicomstatus.FailureUnableToProcess, err
		}
		status = dicomstatus.Success
	}

	return status, nil
}
