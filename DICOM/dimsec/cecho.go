package dimsec

import (
	"rafael/DICOM/media"
	"rafael/DICOM/network"
)

func CEchoReadRQ(pdu network.PDUService, DCO media.DcmObj) bool {
	return DCO.GetUShort(0x00, 0x0100) == 0x30
}

func CEchoWriteRQ(pdu network.PDUService, SOP string) bool {
	var DCO media.DcmObj
	var size uint32
	var valor uint16
	flag := false

	valor = uint16(len(SOP))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	DCO.WriteUint32(0x00, 0x00, "UL", size)          // Length
	DCO.WriteString(0x0000, 0x0002, "UI", SOP)       //SOP Class UID
	DCO.WriteUint16(0x00, 0x0100, "US", 0x30)        //Command Field
	DCO.WriteUint16(0x00, 0x0110, "US", network.Uniq16odd()) //Message ID
	DCO.WriteUint16(0x00, 0x0800, "US", 0x0101)      //Data Set type

	flag = pdu.Write(DCO, SOP, 0x01)
	return flag
}

func CEchoReadRSP(pdu network.PDUService) bool {
	flag := false
	var DCO media.DcmObj

	if pdu.Read(&DCO) == false {
		return false
	}
	if DCO.GetUShort(0x00, 0x0100) == 0x8030 {
		flag = DCO.GetUShort(0x00, 0x0900) == 0x00
	}
	return flag
}

func CEchoWriteRSP(pdu network.PDUService, DCO media.DcmObj) bool{
	return false
}