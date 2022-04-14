package dimsec

import (
	"errors"

	"git.onebytedata.com/odb/go-dicom/dictionary/tags"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/network/dicomcommand"
	"git.onebytedata.com/odb/go-dicom/network/priority"
)

// CFindReadRQ CFind request read
func CFindReadRQ(pdu network.PDUService) (media.DcmObj, error) {
	return pdu.NextPDU()
}

// CFindWriteRQ CFind request write
func CFindWriteRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) error {
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
	DCO.WriteUint16(tags.CommandField, dicomcommand.CFindRequest)
	DCO.WriteUint16(tags.MessageID, network.Uniq16odd())
	DCO.WriteUint16(tags.Priority, priority.Medium)
	DCO.WriteUint16(tags.CommandDataSetType, 0x0102)

	err := pdu.Write(DCO, SOPClassUID, 0x01)
	if err != nil {
		return err
	}
	return pdu.Write(DDO, SOPClassUID, 0x00)
}

// CFindReadRSP CFind response read
func CFindReadRSP(pdu network.PDUService) (media.DcmObj, int, error) {
	status := -1

	dco, err := pdu.NextPDU()
	if err != nil {
		return nil, status, err
	}

	// Is this a C-Find RSP?
	if dco.GetUShort(tags.CommandField) == dicomcommand.CFindResponse {
		if dco.GetUShort(tags.CommandDataSetType) != 0x0101 {
			ddo, err := pdu.NextPDU()
			if err != nil {
				return nil, status, err
			}
			return ddo, int(dco.GetUShort(tags.Status)), nil
		}
		return nil, int(dco.GetUShort(tags.Status)), nil
	}
	return nil, status, errors.New("ERROR, CFindReadRSP, unknown error")
}

// CFindWriteRSP CFind response write
func CFindWriteRSP(pdu network.PDUService, DCO media.DcmObj, DDO media.DcmObj, status uint16) error {
	DCOR := media.NewEmptyDCMObj()
	var size uint32
	var sopclasslength, leDSType uint16

	DCOR.SetTransferSyntax(DCO.GetTransferSyntax())

	if DDO.TagCount() > 0 {
		leDSType = 0x0102
	} else {
		leDSType = 0x0101
	}
	SOPClassUID := DCO.GetString(tags.AffectedSOPClassUID)
	sopclasslength = uint16(len(SOPClassUID))
	if sopclasslength > 0 {
		if sopclasslength%2 == 1 {
			sopclasslength++
		}

		size = uint32(8 + sopclasslength + 8 + 2 + 8 + 2 + 8 + 2)

		DCOR.WriteUint32(tags.CommandGroupLength, size)
		DCOR.WriteString(tags.AffectedSOPClassUID, SOPClassUID)
		DCOR.WriteUint16(tags.CommandField, dicomcommand.CFindResponse)
		valor := DCO.GetUShort(tags.MessageID)
		DCOR.WriteUint16(tags.MessageIDBeingRespondedTo, valor)
		DCOR.WriteUint16(tags.CommandDataSetType, leDSType)
		DCOR.WriteUint16(tags.Status, status)
		err := pdu.Write(DCOR, SOPClassUID, 0x01)
		if err != nil {
			return err
		}
		if DDO.TagCount() > 0 {
			return pdu.Write(DDO, SOPClassUID, 0x00)
		}
	}
	return errors.New("ERROR, CFindReadRSP, unknown error")
}
