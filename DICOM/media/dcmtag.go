package media

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

// DcmTag DICOM tag structure
type DcmTag struct {
	Group     uint16
	Element   uint16
	Length    uint32
	VR        string
	Data      []byte
	BigEndian bool
}

// GetUShort convert tag.Data to uint16
func (tag *DcmTag) GetUShort() uint16 {
	var val uint16
	if tag.Length == 2 {
		if tag.BigEndian {
			val = binary.BigEndian.Uint16(tag.Data)
		} else {
			val = binary.LittleEndian.Uint16(tag.Data)
		}
	} else {
		val = 0
	}
	return val
}

// GetUInt convert tag.Data to uint32
func (tag *DcmTag) GetUInt() uint32 {
	var val uint32
	if tag.Length == 4 {
		if tag.BigEndian {
			val = binary.BigEndian.Uint32(tag.Data)
		} else {
			val = binary.LittleEndian.Uint32(tag.Data)
		}
	} else {
		val = 0
	}
	return val
}

// GetString convert tag.Data to string
func (tag *DcmTag) GetString() string {
	n := bytes.IndexByte(tag.Data, 0)
	if n == -1 {
		n = int(tag.Length)
	}
	val := strings.TrimSpace(string(tag.Data[:n]))
	return val
}

// GetFloat convert tag.Data to float32
func (tag *DcmTag) GetFloat() float32 {
	val := tag.GetString()
	if s, err := strconv.ParseFloat(val, 32); err == nil {
		return float32(s)
	}
	return 0.0
}

// WriteSeq - Create an SQ tag from a DICOM Object
func (tag *DcmTag) WriteSeq(group uint16, element uint16, seq DcmObj) {
	var bufdata BufData

	bufdata.BigEndian = seq.BigEndian
	tag.BigEndian = seq.BigEndian
	tag.Group = group
	tag.Element = element
	if tag.Group == 0xFFFE {
		tag.VR = ""
	} else {
		tag.VR = "SQ"
	}
	for i := 0; i < seq.TagCount(); i++ {
		temptag := seq.GetTag(i)
		bufdata.WriteTag(temptag, seq.ExplicitVR)
	}
	tag.Length = uint32(bufdata.Ms.Size)
	if tag.Length%2 == 1 {
		tag.Length++
		bufdata.Ms.Write([]byte(nil), 1)
	}
	if tag.Length > 0 {
		tag.Data = make([]byte, tag.Length)
		bufdata.setPosition(0)
		bufdata.Ms.Read(tag.Data, int(tag.Length))
	}
}
