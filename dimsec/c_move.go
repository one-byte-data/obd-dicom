package dimsec

import (
	"errors"

	"github.com/one-byte-data/obd-dicom/dictionary/tags"
	"github.com/one-byte-data/obd-dicom/media"
	"github.com/one-byte-data/obd-dicom/network"
	"github.com/one-byte-data/obd-dicom/network/dicomcommand"
	"github.com/one-byte-data/obd-dicom/network/dicomstatus"
	"github.com/one-byte-data/obd-dicom/network/priority"
)

// CMoveReadRQ CMove request read
func CMoveReadRQ(pdu network.PDUService) (media.DcmObj, error) {
	return pdu.NextPDU()
}

// CMoveWriteRQ CMove request write
func CMoveWriteRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string, AETDest string) error {
	DCO := media.NewEmptyDCMObj()
	var size uint32
	var valor, largo uint16

	largo = uint16(len(AETDest))
	if largo%2 == 1 {
		largo++
	}

	valor = uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + largo + 8 + 2 + 8 + 2)

	DCO.WriteUint32(tags.CommandGroupLength, size)
	DCO.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
	DCO.WriteUint16(tags.CommandField, dicomcommand.CMoveRequest)
	DCO.WriteUint16(tags.MessageID, network.Uniq16odd())
	DCO.WriteString(tags.MoveDestination, AETDest)
	DCO.WriteUint16(tags.Priority, priority.Medium)
	DCO.WriteUint16(tags.CommandDataSetType, 0x0102)

	err := pdu.Write(DCO, SOPClassUID, 0x01)
	if err != nil {
		return err
	}
	return pdu.Write(DDO, SOPClassUID, 0x00)
}

// CMoveReadRSP CMove response read
func CMoveReadRSP(pdu network.PDUService, pending *int) (media.DcmObj, uint16, error) {
	status := dicomstatus.FailureUnableToProcess
	dco, err := pdu.NextPDU()
	if err != nil {
		return nil, dicomstatus.FailureUnableToProcess, err
	}

	if dco.GetUShort(tags.CommandField) == dicomcommand.CMoveResponse {
		if dco.GetUShort(tags.CommandDataSetType) != 0x0101 {
			ddo, err := pdu.NextPDU()
			if err != nil {
				return nil, dicomstatus.FailureUnableToProcess, err
			}
			status = dco.GetUShort(tags.Status)
			*pending = int(dco.GetUShort(tags.NumberOfRemainingSuboperations))
			return ddo, status, nil
		}
		status = dco.GetUShort(tags.Status)
		*pending = -1
	}

	return nil, status, nil
}

// CMoveWriteRSP CMove response write
func CMoveWriteRSP(pdu network.PDUService, DCO media.DcmObj, status uint16, pending uint16) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32

	DCOR.SetTransferSyntax(DCO.GetTransferSyntax())

	SOPClassUID := DCO.GetString(tags.AffectedSOPClassUID)
	sopclasslength := uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2 + 8 + 2 + 8 + 2)

		DCOR.WriteUint32(tags.CommandGroupLength, size)
		DCOR.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
		DCOR.WriteUint16(tags.CommandField, dicomcommand.CMoveResponse)
		valor := DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(tags.MessageIDBeingRespondedTo, valor)
		DCOR.WriteUint16(tags.CommandDataSetType, 0x101)
		DCOR.WriteUint16(tags.Status, status)
		DCOR.WriteUint16(tags.NumberOfRemainingSuboperations, pending)

		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CMoveWriteRSP, unknown error")
}
