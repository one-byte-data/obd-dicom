package network

import (
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// UIDitem - UIDitem
type UIDitem struct {
	ItemType  byte
	Reserved1 byte
	Length    uint16
	UIDName   string
}

// Size - Size
func (uid *UIDitem) Size() uint16 {
	return uid.Length + 4
}

// NewUIDitem - NewUIDitem
func NewUIDitem(UIDName string, ItemType byte) *UIDitem {
	return &UIDitem{
		ItemType: ItemType,
		UIDName:  UIDName,
		Length:   uint16(len(UIDName)),
	}
}

func (uid *UIDitem) Write(conn net.Conn) bool {
	var bd media.BufData
	bd.BigEndian = true
	bd.WriteByte(uid.ItemType)
	bd.WriteByte(uid.Reserved1)
	bd.WriteUint16(uid.Length)
	bd.WriteString(uid.UIDName)
	return bd.Send(conn)
}

func (uid *UIDitem) Read(conn net.Conn) bool {
	uid.ItemType = ReadByte(conn)
	return uid.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (uid *UIDitem) ReadDynamic(conn net.Conn) bool {
	uid.Reserved1 = ReadByte(conn)
	uid.Length = ReadUint16(conn)
	buffer := make([]byte, uid.Length)
	conn.Read(buffer)
	uid.UIDName = string(buffer)
	return true
}
