package network

import (
	"bufio"

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

func (uid *UIDitem) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(uid.ItemType)
	bd.WriteByte(uid.Reserved1)
	bd.WriteUint16(uid.Length)
	bd.WriteString(uid.UIDName)

	return bd.Send(rw)
}

func (uid *UIDitem) Read(rw *bufio.ReadWriter) (err error) {
	uid.ItemType, err = ReadByte(rw)
	if err != nil {
		return
	}
	return uid.ReadDynamic(rw)
}

// ReadDynamic - ReadDynamic
func (uid *UIDitem) ReadDynamic(rw *bufio.ReadWriter) (err error) {
	uid.Reserved1, err = ReadByte(rw)
	if err != nil {
		return
	}
	uid.Length, err = ReadUint16(rw)
	if err != nil {
		return
	}

	buffer := make([]byte, uid.Length)
	_, err = rw.Read(buffer)
	if err != nil {
		return
	}

	uid.UIDName = string(buffer)
	return
}
