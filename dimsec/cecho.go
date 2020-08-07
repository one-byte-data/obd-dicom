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

	DCO.WriteUint32(0x00, 0x00, "UL", size)                  // Length
	DCO.WriteString(0x0000, 0x0002, "UI", SOPClassUID)       //SOP Class UID
	DCO.WriteUint16(0x00, 0x0100, "US", commandtype.CEcho)                //Command Field
	DCO.WriteUint16(0x00, 0x0110, "US", network.Uniq16odd()) //Message ID
	DCO.WriteUint16(0x00, 0x0800, "US", 0x0101)              //Data Set type

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

		DCOR.WriteUint32(0x00, 0x00, "UL", size)        // Length
		DCOR.WriteString(0x00, 0x02, "UI", SOPClassUID) //SOP Class UID
		DCOR.WriteUint16(0x00, 0x0100, "US", 0x8030)    //Command Field
		valor = DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(0x00, 0x0110, "US", valor) //Message ID
		valor = DCO.GetUShort(tags.CommandDataSetType)
		DCOR.WriteUint16(0x00, 0x0800, "US", valor) //Data Set type
		DCOR.WriteUint16(0x00, 0x0900, "US", 0x00)  //Data Set type
		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CEchoReadRSP, unknown error")
}
