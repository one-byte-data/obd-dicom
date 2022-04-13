package network

import (
	"bufio"
	"errors"
	"strconv"

	"git.onebytedata.com/odb/go-dicom/media"
)

// UserInformation - UserInformation
type UserInformation interface {
	GetItemType() byte
	SetItemType(t byte)
	GetMaxSubLength() MaximumSubLength
	SetMaxSubLength(length MaximumSubLength)
	Size() uint16
	GetImpClass() UIDitem
	SetImpClassUID(name string)
	GetImpVersion() UIDitem
	SetImpVersionName(name string)
	Write(rw *bufio.ReadWriter) (err error)
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type userInformation struct {
	ItemType        byte //0x50
	Reserved1       byte
	Length          uint16
	UserInfoBaggage uint32
	MaxSubLength    MaximumSubLength
	AsyncOpWindow   AsyncOperationWindow
	SCPSCURole      RoleSelect
	ImpClass        UIDitem
	ImpVersion      UIDitem
}

// NewUserInformation - NewUserInformation
func NewUserInformation() UserInformation {
	return &userInformation{
		ItemType:      0x50,
		MaxSubLength:  NewMaximumSubLength(),
		AsyncOpWindow: NewAsyncOperationWindow(),
		SCPSCURole:    NewRoleSelect(),
	}
}

func (ui *userInformation) GetItemType() byte {
	return ui.ItemType
}

func (ui *userInformation) SetItemType(t byte) {
	ui.ItemType = t
}

func (ui *userInformation) GetMaxSubLength() MaximumSubLength {
	return ui.MaxSubLength
}

func (ui *userInformation) SetMaxSubLength(length MaximumSubLength) {
	ui.MaxSubLength = length
}

func (ui *userInformation) Size() uint16 {
	ui.Length = ui.MaxSubLength.Size()
	ui.Length += ui.ImpClass.Size()
	ui.Length += ui.ImpVersion.Size()
	return ui.Length + 4
}

func (ui *userInformation) GetImpClass() UIDitem {
	return ui.ImpClass
}

func (ui *userInformation) SetImpClassUID(name string) {
	ui.ImpClass.ItemType = 0x52
	ui.ImpClass.Reserved1 = 0x00
	ui.ImpClass.UIDName = name
	ui.ImpClass.Length = uint16(len(name))
}

func (ui *userInformation) GetImpVersion() UIDitem {
	return ui.ImpVersion
}

func (ui *userInformation) SetImpVersionName(name string) {
	ui.ImpVersion.ItemType = 0x55
	ui.ImpVersion.Reserved1 = 0x00
	ui.ImpVersion.UIDName = name
	ui.ImpVersion.Length = uint16(len(name))
}

func (ui *userInformation) Write(rw *bufio.ReadWriter) (err error) {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	ui.Size()
	bd.WriteByte(ui.ItemType)
	bd.WriteByte(ui.Reserved1)
	bd.WriteUint16(ui.Length)

	if err = bd.Send(rw); err == nil {
		ui.MaxSubLength.Write(rw)
		ui.ImpClass.Write(rw)
		ui.ImpVersion.Write(rw)
	}

	return
}

func (ui *userInformation) Read(ms media.MemoryStream) (err error) {
	ui.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return ui.ReadDynamic(ms)
}

func (ui *userInformation) ReadDynamic(ms media.MemoryStream) (err error) {
	ui.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	ui.Length, err = ms.GetUint16()
	if err != nil {
		return err
	}

	Count := int(ui.Length)
	for Count > 0 {
		TempByte, err := ms.GetByte()
		if err != nil {
			return err
		}

		switch TempByte {
		case 0x51:
			ui.MaxSubLength.ReadDynamic(ms)
			Count = Count - int(ui.MaxSubLength.Size())
			break
		case 0x52:
			ui.ImpClass.ReadDynamic(ms)
			Count = Count - int(ui.ImpClass.Size())
			break
		case 0x53:
			ui.AsyncOpWindow.ReadDynamic(ms)
			Count = Count - int(ui.AsyncOpWindow.Size())
			break
		case 0x54:
			ui.SCPSCURole.ReadDynamic(ms)
			Count = Count - int(ui.SCPSCURole.Size())
			ui.UserInfoBaggage += uint32(ui.SCPSCURole.Size())
			break
		case 0x55:
			ui.ImpVersion.ReadDynamic(ms)
			Count = Count - int(ui.ImpVersion.Size())
			break
		default:
			ui.UserInfoBaggage = uint32(Count)
			Count = -1
			return errors.New("ERROR, user::ReadDynamic, unknown TempByte: " + strconv.Itoa(int(TempByte)))
		}
	}

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, user::ReadDynamic, Count is not zero")
}
