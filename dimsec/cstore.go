package dimsec

import (
	"errors"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

// CStoreReadRQ CStore request read
func CStoreReadRQ(pdu network.PDUService, command media.DcmObj) (media.DcmObj, error) {
	return pdu.NextPDU()
}

// CStoreWriteRQ CStore request write
func CStoreWriteRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) error {
	DCO := media.NewEmptyDCMObj()
	var size uint32
	var valor uint16

	valor = uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	SOPInstance := DDO.GetString(0x08, 0x18)
	length := uint32(len(SOPInstance))
	if length%2 == 1 {
		length++
		size = size + 8 + length
	}

	DCO.WriteUint32(0x00, 0x00, "UL", size)                  // Length
	DCO.WriteString(0x0000, 0x0002, "UI", SOPClassUID)       //SOP Class UID
	DCO.WriteUint16(0x00, 0x0100, "US", 0x01)                //Command Field
	DCO.WriteUint16(0x00, 0x0110, "US", network.Uniq16odd()) //Message ID
	DCO.WriteUint16(0x00, 0x0700, "US", 0x00)                //Priority
	DCO.WriteUint16(0x00, 0x0800, "US", 0x0102)              //Data Set type

	if length > 0 {
		DCO.WriteString(0x00, 0x1000, "UI", SOPInstance)
	}

	err := pdu.Write(DCO, SOPClassUID, 0x01)
	if err != nil {
		return err
	}
	return pdu.Write(DDO, SOPClassUID, 0x00)
}

// CStoreReadRSP CStore response read
func CStoreReadRSP(pdu network.PDUService) (int, error) {
	dco, err := pdu.NextPDU()
	if err != nil {
		return -1, err
	}
	// Is this a C-Store RSP?
	if dco.GetUShort(0x00, 0x0100) == 0x8001 {
		return int(dco.GetUShort(0x00, 0x0900)), nil
	}
	return -1, errors.New("ERROR, CStoreReadRSP, unknown error")
}

// CStoreWriteRSP CStore response write
func CStoreWriteRSP(pdu network.PDUService, DCO media.DcmObj, status uint16) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32
	var sopclasslength, sopinstancelength uint16

	DCOR.SetTransferSyntax(DCO.GetTransferSyntax())
	SOPClassUID := DCO.GetString(0x00, 0x02)
	sopclasslength = uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		SOPInstance := DCO.GetString(0x00, 0x1000)
		sopinstancelength = uint16(len(SOPClassUID))
		if sopinstancelength > 0 {
			if sopinstancelength%2 == 1 {
				sopinstancelength++
			}

			size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2 + 8 + sopinstancelength)

			DCOR.WriteUint32(0x00, 0x00, "UL", size)        // Length
			DCOR.WriteString(0x00, 0x02, "UI", SOPClassUID) //SOP Class UID
			DCOR.WriteUint16(0x00, 0x0100, "US", 0x8001)    //Command Field
			valor := DCO.GetUShort(0x00, 0x0110)
			DCOR.WriteUint16(0x00, 0x0120, "US", valor)  //Message ID
			DCOR.WriteUint16(0x00, 0x0800, "US", 0x0101) //Data Set type
			DCOR.WriteUint16(0x00, 0x0900, "US", status) //Data Set type
			DCOR.WriteString(0x00, 0x1000, "UI", SOPInstance)
			return pdu.Write(DCOR, SOPClassUID, 0x01)
		}
	}
	return errors.New("ERROR, CStoreWriteRSP, unknown error")
}
