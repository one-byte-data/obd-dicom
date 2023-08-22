package dimsec

import (
	"errors"

	"github.com/one-byte-data/obd-dicom/dictionary/tags"
	"github.com/one-byte-data/obd-dicom/media"
	"github.com/one-byte-data/obd-dicom/network"
	"github.com/one-byte-data/obd-dicom/network/dicomcommand"
	"github.com/one-byte-data/obd-dicom/network/dicomstatus"
)

// CEchoReadRQ CEcho request read
func CEchoReadRQ(DCO media.DcmObj) bool {
	return DCO.GetUShort(tags.CommandField) == dicomcommand.CEchoRequest
}

// CEchoWriteRQ CEcho request write
func CEchoWriteRQ(pdu network.PDUService, SOPClassUID string) error {
	DCO := media.NewEmptyDCMObj()
	var size uint32

	valor := uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	DCO.WriteUint32(tags.CommandGroupLength, size)
	DCO.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
	DCO.WriteUint16(tags.CommandField, dicomcommand.CEchoRequest)
	DCO.WriteUint16(tags.MessageID, network.Uniq16odd())
	DCO.WriteUint16(tags.CommandDataSetType, 0x0101)

	return pdu.Write(DCO, SOPClassUID, 0x01)
}

// CEchoReadRSP CEcho response read
func CEchoReadRSP(pdu network.PDUService) error {
	dco, err := pdu.NextPDU()
	if err != nil {
		return errors.New("ERROR, CEchoReadRSP, failed pdu.Read(&DCO)")
	}
	if dco.GetUShort(tags.CommandField) == dicomcommand.CEchoResponse {
		if dco.GetUShort(tags.Status) == dicomstatus.Success {
			return nil
		}
	}
	return nil
}

// CEchoWriteRSP CEcho response write
func CEchoWriteRSP(pdu network.PDUService, DCO media.DcmObj) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32
	var valor uint16

	DCOR.SetTransferSyntax(DCO.GetTransferSyntax())
	SOPClassUID := DCO.GetString(tags.AffectedSOPClassUID)
	valor = uint16(len(SOPClassUID))
	if valor > 0 {
		if valor%2 == 1 {
			valor++
		}

		size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

		DCOR.WriteUint32(tags.CommandGroupLength, size)
		DCOR.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
		DCOR.WriteUint16(tags.CommandField, dicomcommand.CEchoResponse)
		valor = DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(tags.MessageIDBeingRespondedTo, valor)
		valor = DCO.GetUShort(tags.CommandDataSetType)
		DCOR.WriteUint16(tags.CommandDataSetType, valor)
		DCOR.WriteUint16(tags.Status, dicomstatus.Success)
		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CEchoReadRSP, unknown error")
}
