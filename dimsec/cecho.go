package dimsec

import (
	"errors"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network/commandtype"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
)

// CEchoReadRQ CEcho request read
func CEchoReadRQ(pdu network.PDUService, DCO media.DcmObj) bool {
	return DCO.GetUShort(tags.CommandField) == 0x30
}

// CEchoWriteRQ CEcho request write
func CEchoWriteRQ(pdu network.PDUService, SOPClassUID string) error {
	DCO := media.NewEmptyDCMObj()
	var size uint32
	var valor uint16

	valor = uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	DCO.WriteUint32(tags.CommandGroupLength, size) 
	DCO.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
	DCO.WriteUint16(tags.CommandField, commandtype.CEcho)
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
	if dco.GetUShort(tags.CommandField) == 0x8030 {
		if dco.GetUShort(tags.Status) == 0x00 {

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
		DCOR.WriteUint16(tags.CommandField, 0x8030)
		valor = DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(tags.MessageID, valor)
		valor = DCO.GetUShort(tags.CommandDataSetType)
		DCOR.WriteUint16(tags.CommandDataSetType, valor)
		DCOR.WriteUint16(tags.Status, 0x00)
		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CEchoReadRSP, unknown error")
}
