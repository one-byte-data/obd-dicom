package dimsec

import (
	"errors"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

// CMoveReadRQ CMove request read
func CMoveReadRQ(pdu network.PDUService, DCO media.DcmObj, DDO media.DcmObj) error {
	if DCO.TagCount() != 0 {
		// Is this a C-Move?
		if DCO.GetUShort(0x00, 0x100) == 0x21 {
			// Does it have data?
			if DCO.GetUShort(0x00, 0x0800) != 0x0101 {
				return pdu.Read(DDO)
			}
		}
	}
	return errors.New("ERROR, CMoveReadRQ, unknown error")
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

	DCO.WriteUint32(0x00, 0x00, "UL", size)                  // Length
	DCO.WriteString(0x0000, 0x0002, "UI", SOPClassUID)       //SOP Class UID
	DCO.WriteUint16(0x00, 0x0100, "US", 0x21)                //Command Field
	DCO.WriteUint16(0x00, 0x0110, "US", network.Uniq16odd()) //Message ID
	DCO.WriteString(0x00, 0x0600, "AE", AETDest)             // Destination AET
	DCO.WriteUint16(0x00, 0x0700, "US", 0x00)                // Priority
	DCO.WriteUint16(0x00, 0x0800, "US", 0x0102)              //Data Set type

	err := pdu.Write(DCO, SOPClassUID, 0x01)
	if err != nil {
		return err
	}
	return pdu.Write(DDO, SOPClassUID, 0x00)
}

// CMoveReadRSP CMove response read
func CMoveReadRSP(pdu network.PDUService, DDO media.DcmObj, pending *int) (int, error) {
	DCO := media.NewEmptyDCMObj()
	status := -1

	if err := pdu.Read(DCO); err != nil {
		return status, err
	}
	// Is this a C-Find RSP?
	if DCO.GetUShort(0x00, 0x0100) == 0x8021 {
		if DCO.GetUShort(0x00, 0x0800) != 0x0101 {
			err := pdu.Read(DDO)
			if err != nil {
				return status, err
			}
			status = int(DCO.GetUShort(0x00, 0x0900))
			*pending = int(DCO.GetUShort(0x00, 0x1020))
		} else {
			status = int(DCO.GetUShort(0x00, 0x0900))
			*pending = -1
		}
	}
	return status, nil
}

// CMoveWriteRSP CMove response write
func CMoveWriteRSP(pdu network.PDUService, DCO media.DcmObj, status uint16, pending uint16) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32

	DCOR.SetTransferSyntax(DCO.GetTransferSynxtax())

	SOPClassUID := DCO.GetString(0x00, 0x02)
	sopclasslength := uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2 + 8 + 2 + 8 + 2)

		DCOR.WriteUint32(0x00, 0x00, "UL", size)        // Length
		DCOR.WriteString(0x00, 0x02, "UI", SOPClassUID) //SOP Class UID
		DCOR.WriteUint16(0x00, 0x0100, "US", 0x8021)    //Command Field
		valor := DCO.GetUShort(0x00, 0x0110)
		DCOR.WriteUint16(0x00, 0x0120, "US", valor)   //Message ID
		DCOR.WriteUint16(0x00, 0x0800, "US", 0x101)   //Data Set type
		DCOR.WriteUint16(0x00, 0x0900, "US", status)  //Status
		DCOR.WriteUint16(0x00, 0x1020, "US", pending) //Pending

		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CMoveWriteRSP, unknown error")
}
