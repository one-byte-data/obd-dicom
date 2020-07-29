package services

import (
	"errors"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/dimsec"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

// SCU - inteface to a scu
type SCU interface {
	EchoSCU(timeout int) error
	FindSCU(Query media.DcmObj, Results *[]media.DcmObj, timeout int) (int, error)
	MoveSCU(destAET string, Query media.DcmObj, timeout int) (int, error)
	StoreSCU(FileName string, timeout int) error
	openAssociation(pdu network.PDUService, AbstractSyntax string, timeout int) error
	writeStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) (int, error)
}

type scu struct {
	destination *network.Destination
}

// NewSCU - Creates an interface to scu
func NewSCU(destination *network.Destination) SCU {
	return &scu{
		destination: destination,
	}
}

func (d *scu) EchoSCU(timeout int) error {
	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, "1.2.840.10008.1.1", timeout)
	if err != nil {
		return err
	}
	err = dimsec.CEchoWriteRQ(pdu, "1.2.840.10008.1.1")
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

func (d *scu) FindSCU(Query media.DcmObj, Results *[]media.DcmObj, timeout int) (int, error) {
	status := 1
	SOPClassUID := "1.2.840.10008.5.1.4.1.2.2.1"

	pdu := network.NewPDUService()
	err := d.openAssociation(pdu, SOPClassUID, timeout)
	if err != nil {
		return -1, err
	}
	err = dimsec.CFindWriteRQ(pdu, Query, SOPClassUID)
	if err != nil {
		return -1, err
	}
	for (status != -1) && (status != 0) {
		DDO := media.NewEmptyDCMObj()
		status, err = dimsec.CFindReadRSP(pdu, DDO)
		if err != nil {
			return status, err
		}
		if (status == 0xFF00) || (status == 0xFF01) {
			*Results = append(*Results, DDO)
		}
	}

	pdu.Close()
	return status, nil
}

func (d *scu) MoveSCU(destAET string, Query media.DcmObj, timeout int) (int, error) {
	var pending int
	status := 0xFF00
	SOPClassUID := "1.2.840.10008.5.1.4.1.2.2.2"

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
		DDO := media.NewEmptyDCMObj()
		status, err = dimsec.CMoveReadRSP(pdu, DDO, &pending)
		if err != nil {
			return -1, err
		}
		DDO.DumpTags()
	}

	pdu.Close()
	return status, nil
}

func (d *scu) StoreSCU(FileName string, timeout int) error {
	DDO, err := media.NewDCMObjFromFile(FileName)
	if err != nil {
		return err
	}

	SOPClassUID := DDO.GetString(0x08, 0x16)
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

func (d *scu) openAssociation(pdu network.PDUService, AbstractSyntax string, timeout int) error {
	pdu.SetRQCallingAE(d.destination.CallingAE)
	pdu.SetRQCalledAE(d.destination.CalledAE)
	pdu.SetTimeout(timeout)

	network.Resetuniq()
	PresContext := network.NewPresentationContext()
	PresContext.SetAbstractSyntax(AbstractSyntax)
	PresContext.AddTransferSyntax("1.2.840.10008.1.2")
	PresContext.AddTransferSyntax("1.2.840.10008.1.2.1")
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

	if TrnSyntOUT == DDO.GetTransferSynxtax() {
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return status, err
		}
		status = 0
	} else {
		DDO.SetTransferSyntax(TrnSyntOUT)
		DDO.SetExplicitVR(true)
		DDO.SetBigEndian(false)
		if TrnSyntOUT == "1.2.840.10008.1.2" {
			DDO.SetExplicitVR(false)
		}
		if TrnSyntOUT == "1.2.840.10008.1.2.2" {
			DDO.SetBigEndian(true)
		}
		err := dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)
		if err != nil {
			return -1, err
		}
		status = 0
	}

	return status, nil
}
