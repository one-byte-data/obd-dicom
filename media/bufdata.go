package media

import (
	"encoding/binary"
	"net"
)

// BufData buffer manipulation class
type BufData struct {
	BigEndian bool
	Ms        MemoryStream
}

func (bd *BufData) getPosition() int {
	return bd.Ms.Position
}

func (bd *BufData) setPosition(pos int) {
	bd.Ms.Position = pos
}

func (bd *BufData) ReadByte() byte{
	c:= make([]byte, 1)
	bd.Ms.Read(c, 1)
	return c[0]
}

func (bd *BufData) ReadUint16() uint16 {
	var val uint16

	c := make([]byte, 2)
	bd.Ms.Read(c, 2)
	if bd.BigEndian {
		val = binary.BigEndian.Uint16(c)
	} else {
		val = binary.LittleEndian.Uint16(c)
	}
	return val
}

func (bd *BufData) ReadUint32() uint32 {
	var val uint32

	c := make([]byte, 4)
	bd.Ms.Read(c, 4)
	if bd.BigEndian {
		val = binary.BigEndian.Uint32(c)
	} else {
		val = binary.LittleEndian.Uint32(c)
	}
	return val
}

func (bd *BufData) WriteByte(val byte) {
	c := make([]byte, 1)
	c[0] = val
	bd.Ms.Write(c, 1)
}

func (bd *BufData) WriteUint16(val uint16) {
	c := make([]byte, 2)
	if bd.BigEndian {
		binary.BigEndian.PutUint16(c, val)
	} else {
		binary.LittleEndian.PutUint16(c, val)
	}
	bd.Ms.Write(c, 2)
}

func (bd *BufData) WriteUint32(val uint32) {
	c := make([]byte, 4)
	if bd.BigEndian {
		binary.BigEndian.PutUint32(c, val)
	} else {
		binary.LittleEndian.PutUint32(c, val)
	}
	bd.Ms.Write(c, 4)
}

func (bd *BufData) readString(length int) string {
	temp := make([]byte, length)
	bd.Ms.Read(temp, length)
	val := string(temp[:])
	return val
}

func (bd *BufData) WriteString(val string) {
	bd.Ms.Write([]byte(val), len(val))
}

// ReadTag - read a single tag from the Stream
func (bd *BufData) ReadTag(tag *DcmTag, explicitVR bool) bool {
	tag.VR = ""
	internalVR := explicitVR
	tag.Group = bd.ReadUint16()
	tag.Element = bd.ReadUint16()
	if tag.Group == 0x0002 {
		internalVR = true
	}
	if (tag.Group != 0x0000) && (tag.Group != 0xfffe) && (internalVR) {
		tag.VR = bd.readString(2)
		if (tag.VR == "OB") || (tag.VR == "OW") || (tag.VR == "SQ") || (tag.VR == "UN") || (tag.VR == "UT") {
			bd.ReadUint16()
			tag.Length = bd.ReadUint32()
		} else {
			tag.Length = uint32(bd.ReadUint16())
		}
	} else {
		if internalVR == false {
			tag.VR = AddVRData(tag.Group, tag.Element)
		}
		tag.Length = bd.ReadUint32()
	}

	if (tag.Length != 0) && (tag.Length != 0xFFFFFFFF) {
		tag.Data = make([]byte, tag.Length)
		if bd.Ms.Read(tag.Data, int(tag.Length)) != int(tag.Length) {
			return false
		}
	}
	return true
}

// WriteTag - Write a single tag to stream
func (bd *BufData) WriteTag(tag DcmTag, explicitVR bool) {
	bd.WriteUint16(tag.Group)
	bd.WriteUint16(tag.Element)
	if (tag.Group != 0x0000) && (tag.Group != 0xfffe) && (explicitVR) {
		bd.Ms.Write([]byte(tag.VR), 2)
		if (tag.VR == "OB") || (tag.VR == "OW") || (tag.VR == "SQ") || (tag.VR == "UN") || (tag.VR == "UT") {
			bd.WriteUint16(0)
			bd.WriteUint32(tag.Length)
		} else {
			bd.WriteUint16(uint16(tag.Length))
		}
	} else {
		bd.WriteUint32(tag.Length)
	}
	if (tag.Length != 0) && (tag.Length != 0xFFFFFFFF) {
		bd.Ms.Write(tag.Data, int(tag.Length))
	}
}

// WriteStringTag - Writes a String to a DICOM tag
func (bd *BufData) WriteStringTag(group uint16, element uint16, vr string, content string, explicitVR bool) {
	var length uint32

	length = uint32(len(content))
	if length%2 == 1 {
		length++
		if vr == "UI" {
			content = content + string(0)
		} else {
			content = content + " "
		}
	}
	tag := DcmTag{group, element, length, vr, []byte(content), false}
	bd.WriteTag(tag, explicitVR)
}

// ReadMeta - Read Meta Header
func (bd *BufData) ReadMeta() string {
	bs := make([]byte, 4)
	TransferSyntax := ""
	pos := 0

	bd.Ms.Position = 128
	bd.Ms.Read(bs, 4)
	if string(bs[:4]) == "DICM" {
		var tag DcmTag
		fin := false
		for (pos < bd.Ms.Size) && (!fin) {
			pos = bd.getPosition()
			bd.ReadTag(&tag, true)
			if (tag.Group == 0x02) && (tag.Element == 0x010) {
				TransferSyntax = tag.GetString()
			}
			if tag.Group > 0x02 {
				fin = true
			}
		}
	}
	bd.setPosition(pos)
	return TransferSyntax
}

// WriteMeta - Write Meta Header
func (bd *BufData) WriteMeta(SOPClassUID string, SOPInstanceUID string, TransferSyntax string) {
	explicitVR := true
	buffer := make([]byte, 128)
	var largo uint32
	var tag DcmTag

	bd.Ms.Write(buffer, 128)
	bd.Ms.Write([]byte("DICM"), 4)
	tag = DcmTag{0x02, 0x00, 4, "UL", []byte{0, 0, 0, 0}, false}
	bd.WriteTag(tag, explicitVR)
	tag = DcmTag{0x02, 0x01, 2, "OB", []byte{0x00, 0x01}, false}
	bd.WriteTag(tag, explicitVR)

	bd.WriteStringTag(0x02, 0x02, "UI", SOPClassUID, explicitVR)
	bd.WriteStringTag(0x02, 0x03, "UI", SOPInstanceUID, explicitVR)
	bd.WriteStringTag(0x02, 0x10, "UI", TransferSyntax, explicitVR)

	// Implementation Class UID
	bd.WriteStringTag(0x02, 0x12, "UI", "123456", explicitVR)
	// Implementation Version Name
	bd.WriteStringTag(0x02, 0x13, "SH", "OneByteData", explicitVR)

	// calculate group length and go Back to group size tag
	ptr := bd.getPosition()
	largo = uint32(bd.Ms.Size - 12 - 128 - 4)
	binary.LittleEndian.PutUint32(buffer, largo)
	bd.setPosition(128 + 4 + 8)
	bd.Ms.Write(buffer, 4)
	bd.setPosition(ptr)
}

// ReadObj - Read a DICOM Object from a BufData
func (bd *BufData) ReadObj(obj *DcmObj) bool {
	flag := false
	var tag DcmTag

	for bd.Ms.Position < bd.Ms.Size {
		if bd.ReadTag(&tag, obj.ExplicitVR) {
			if obj.ExplicitVR == false {
				tag.VR = AddVRData(tag.Group, tag.Element)
			}
			obj.Add(tag)
		}
		flag = true
	}
	return flag
}

// WriteObj - Write a DICOM Object to a BufData
func (bd *BufData) WriteObj(obj *DcmObj) {
//	bd.BigEndian = BigEndian
// Si lo limpio elimino el meta!!
//	bd.Ms.Clear()
	for i := 0; i < obj.TagCount(); i++ {
		tag := obj.GetTag(i)
		bd.WriteTag(tag, obj.ExplicitVR)
	}
}

func (bd *BufData) Send(conn net.Conn) bool {
	buffer := make([]byte, bd.Ms.Size)
	bd.Ms.Position = 0
	bd.Ms.Read(buffer, bd.Ms.Size)
	bd.Ms.Clear()
	_, err := conn.Write(buffer)
	if err != nil {
		return false
	}
	return true
}
