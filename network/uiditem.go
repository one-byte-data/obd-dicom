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

func (uid *UIDitem) Write(conn net.Conn) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(uid.ItemType)
	bd.WriteByte(uid.Reserved1)
	bd.WriteUint16(uid.Length)
	bd.WriteString(uid.UIDName)

	return bd.Send(conn)
}

func (uid *UIDitem) Read(conn net.Conn) (err error) {
	uid.ItemType, err = ReadByte(conn)
	if err != nil {
		return
	}
	return uid.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (uid *UIDitem) ReadDynamic(conn net.Conn) (err error) {
	uid.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	uid.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}

	buffer := make([]byte, uid.Length)
	_, err = conn.Read(buffer)
	if err != nil {
		return
	}

	uid.UIDName = string(buffer)
	return
}
