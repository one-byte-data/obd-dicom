package dimsec

import (
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

func CFindReadRQ(pdu network.PDUService, DCO media.DcmObj, DDO *media.DcmObj) bool {
	if DCO.TagCount() != 0 {
		// Is this a C-Find?
		if DCO.GetUShort(0x00, 0x100) == 0x20 {
			// Does it have data?
			if DCO.GetUShort(0x00, 0x0800) != 0x0101 {
				return pdu.Read(DDO)
			}
		}
	}
	return false
}

func CFindWriteRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) bool {
	var DCO media.DcmObj
	var size uint32
	var valor uint16

	valor = uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	DCO.WriteUint32(0x00, 0x00, "UL", size)                  // Length
	DCO.WriteString(0x0000, 0x0002, "UI", SOPClassUID)       //SOP Class UID
	DCO.WriteUint16(0x00, 0x0100, "US", 0x20)                //Command Field
	DCO.WriteUint16(0x00, 0x0110, "US", network.Uniq16odd()) //Message ID
	DCO.WriteUint16(0x00, 0x0700, "US", 0x00)                //Data Set type
	DCO.WriteUint16(0x00, 0x0800, "US", 0x0102)              //Data Set type

	if pdu.Write(DCO, SOPClassUID, 0x01) {
		return pdu.Write(DDO, SOPClassUID, 0x00)
	}
	return false
}

func CFindReadRSP(pdu network.PDUService, DDO *media.DcmObj) int {
	var DCO media.DcmObj
	status := -1

	if pdu.Read(&DCO) == false {
		return status
	}
	// Is this a C-Find RSP?
	if DCO.GetUShort(0x00, 0x0100) == 0x8020 {
		if DCO.GetUShort(0x00, 0x0800) != 0x0101 {
			if pdu.Read(DDO) {
				status = int(DCO.GetUShort(0x00, 0x0900)) // Return Status
			}
		}
	}
	return status
}

func CFindWriteRSP(pdu network.PDUService, DCO media.DcmObj, DDO media.DcmObj, status uint16) bool {
	var DCOR media.DcmObj
	var size uint32
	var sopclasslength, leDSType uint16
	flag := false

	DCOR.TransferSyntax = DCO.TransferSyntax

	if DDO.TagCount() > 0 {
		leDSType = 0x0102
	} else {
		leDSType = 0x0101
	}
	SOPClassUID := DCO.GetString(0x00, 0x02)
	sopclasslength = uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2)

		DCOR.WriteUint32(0x00, 0x00, "UL", size)        // Length
		DCOR.WriteString(0x00, 0x02, "UI", SOPClassUID) //SOP Class UID
		DCOR.WriteUint16(0x00, 0x0100, "US", 0x8020)    //Command Field
		valor := DCO.GetUShort(0x00, 0x0110)
		DCOR.WriteUint16(0x00, 0x0120, "US", valor)    //Message ID
		DCOR.WriteUint16(0x00, 0x0800, "US", leDSType) //Data Set type
		DCOR.WriteUint16(0x00, 0x0900, "US", status)   //Data Set type
		flag = pdu.Write(DCOR, SOPClassUID, 0x01)
		if flag && (DDO.TagCount() > 0) {
			flag = pdu.Write(DDO, SOPClassUID, 0x00)
		}
	}
	return flag
}
