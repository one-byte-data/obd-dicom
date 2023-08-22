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

// CStoreReadRQ CStore request read
func CStoreReadRQ(pdu network.PDUService, command media.DcmObj) (media.DcmObj, error) {
	return pdu.NextPDU()
}

// CStoreWriteRQ CStore request write
func CStoreWriteRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) error {
	DCO := media.NewEmptyDCMObj()
	var size uint32

	valor := uint16(len(SOPClassUID))
	if valor%2 == 1 {
		valor++
	}

	size = uint32(8 + valor + 8 + 2 + 8 + 2 + 8 + 2)

	SOPInstance := DDO.GetString(tags.SOPInstanceUID)
	length := uint32(len(SOPInstance))
	if length%2 == 1 {
		length++
		size = size + 8 + length
	}

	DCO.WriteUint32(tags.CommandGroupLength, size)
	DCO.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
	DCO.WriteUint16(tags.CommandField, dicomcommand.CStoreRequest)
	DCO.WriteUint16(tags.MessageID, network.Uniq16odd())
	DCO.WriteUint16(tags.Priority, priority.Medium)
	DCO.WriteUint16(tags.CommandDataSetType, 0x0102)

	if length > 0 {
		DCO.WriteString(tags.AffectedSOPInstanceUID, SOPInstance)
	}

	err := pdu.Write(DCO, SOPClassUID, 0x01)
	if err != nil {
		return err
	}
	return pdu.Write(DDO, SOPClassUID, 0x00)
}

// CStoreReadRSP CStore response read
func CStoreReadRSP(pdu network.PDUService) (uint16, error) {
	dco, err := pdu.NextPDU()
	if err != nil {
		return dicomstatus.FailureUnableToProcess, err
	}
	// Is this a C-Store RSP?
	if dco.GetUShort(tags.CommandField) == dicomcommand.CStoreResponse {
		return dco.GetUShort(tags.Status), nil
	}
	return dicomstatus.FailureUnableToProcess, errors.New("ERROR, CStoreReadRSP, unknown error")
}

// CStoreWriteRSP CStore response write
func CStoreWriteRSP(pdu network.PDUService, DCO media.DcmObj, status uint16) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32
	var sopclasslength, sopinstancelength uint16

	DCOR.SetTransferSyntax(DCO.GetTransferSyntax())
	SOPClassUID := DCO.GetString(tags.AffectedSOPClassUID)
	sopclasslength = uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		SOPInstance := DCO.GetString(tags.AffectedSOPInstanceUID)
		sopinstancelength = uint16(len(SOPClassUID))
		if sopinstancelength > 0 {
			if sopinstancelength%2 == 1 {
				sopinstancelength++
			}

			size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2 + 8 + sopinstancelength)

			DCOR.WriteUint32(tags.CommandGroupLength, size)
			DCOR.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
			DCOR.WriteUint16(tags.CommandField, dicomcommand.CStoreResponse)
			valor := DCO.GetUShort(tags.MessageID)
			DCOR.WriteUint16(tags.MessageIDBeingRespondedTo, valor)
			DCOR.WriteUint16(tags.CommandDataSetType, 0x0101)
			DCOR.WriteUint16(tags.Status, status)
			DCOR.WriteString(tags.AffectedSOPInstanceUID, SOPInstance)
			return pdu.Write(DCOR, SOPClassUID, 0x01)
		}
	}
	return errors.New("ERROR, CStoreWriteRSP, unknown error")
}
