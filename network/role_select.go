package network

import (
	"bufio"

	"git.onebytedata.com/odb/go-dicom/media"
)

// RoleSelect - RoleSelect
type RoleSelect interface {
	Size() uint16
	Write(rw *bufio.ReadWriter) bool
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type roleSelect struct {
	ItemType  byte //0x54
	Reserved1 byte
	Length    uint16
	SCURole   byte
	SCPRole   byte
	uid       string
}

// NewRoleSelect - NewRoleSelect
func NewRoleSelect() RoleSelect {
	return &roleSelect{
		ItemType: 0x54,
	}
}

func (scpscu *roleSelect) Size() uint16 {
	return scpscu.Length + 4
}

func (scpscu *roleSelect) Write(rw *bufio.ReadWriter) bool {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(scpscu.ItemType)
	bd.WriteByte(scpscu.Reserved1)
	bd.WriteUint16(scpscu.Length)
	bd.WriteUint16(uint16(len(scpscu.uid)))
	bd.Write([]byte(scpscu.uid), len(scpscu.uid))
	bd.WriteByte(scpscu.SCURole)
	bd.WriteByte(scpscu.SCPRole)

	if err := bd.Send(rw); err != nil {
		return false
	}
	return true
}

func (scpscu *roleSelect) Read(ms media.MemoryStream) (err error) {
	scpscu.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return scpscu.ReadDynamic(ms)
}

func (scpscu *roleSelect) ReadDynamic(ms media.MemoryStream) (err error) {
	scpscu.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	scpscu.Length, err = ms.GetUint16()
	if err != nil {
		return err
	}
	tl, err := ms.GetUint16()
	if err != nil {
		return err
	}

	tuid := make([]byte, tl)
	ms.ReadData(tuid)

	scpscu.uid = string(tuid)
	scpscu.SCURole, err = ms.GetByte()
	if err != nil {
		return err
	}
	scpscu.SCPRole, err = ms.GetByte()
	return
}
