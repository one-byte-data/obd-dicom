package dimsec

import (
	"errors"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
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
func CMoveReadRSP(pdu network.PDUService, pending *int) (media.DcmObj, int, error) {
	status := -1

	dco, err := pdu.NextPDU()
	if err != nil {
		return nil, status, err
	}
	// Is this a C-Find RSP?
	if dco.GetUShort(tags.CommandField) == 0x8021 {
		if dco.GetUShort(tags.CommandDataSetType) != 0x0101 {
			ddo, err := pdu.NextPDU()
			if err != nil {
				return nil, status, err
			}
			status = int(dco.GetUShort(tags.Status))
			*pending = int(dco.GetUShort(tags.NumberOfRemainingSuboperations))
			return ddo, status, nil
		}
		status = int(dco.GetUShort(tags.Status))
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

		DCOR.WriteUint32(0x00, 0x00, "UL", size)        // Length
		DCOR.WriteString(0x00, 0x02, "UI", SOPClassUID) //SOP Class UID
		DCOR.WriteUint16(0x00, 0x0100, "US", 0x8021)    //Command Field
		valor := DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(0x00, 0x0120, "US", valor)   //Message ID
		DCOR.WriteUint16(0x00, 0x0800, "US", 0x101)   //Data Set type
		DCOR.WriteUint16(0x00, 0x0900, "US", status)  //Status
		DCOR.WriteUint16(0x00, 0x1020, "US", pending) //Pending

		return pdu.Write(DCOR, SOPClassUID, 0x01)
	}
	return errors.New("ERROR, CMoveWriteRSP, unknown error")
}
