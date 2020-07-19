package network

import (
	"errors"
	"net"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// MaximumSubLength - MaximumSubLength
type MaximumSubLength interface {
	GetMaximumLength() uint32
	SetMaximumLength(length uint32)
	Size() uint16
	Write(conn net.Conn) bool
	Read(conn net.Conn) (err error)
	ReadDynamic(conn net.Conn) (err error)
}

type maximumSubLength struct {
	ItemType      byte //0x51
	Reserved1     byte
	Length        uint16
	MaximumLength uint32
}

// NewMaximumSubLength - NewMaximumSubLength
func NewMaximumSubLength() MaximumSubLength {
	return &maximumSubLength{
		ItemType: 0x51,
		Length:   4,
	}
}

func (maxim *maximumSubLength) GetMaximumLength() uint32 {
	return maxim.MaximumLength
}

func (maxim *maximumSubLength) SetMaximumLength(length uint32) {
	maxim.MaximumLength = length
}

func (maxim *maximumSubLength) Size() uint16 {
	return maxim.Length + 4
}

func (maxim *maximumSubLength) Write(conn net.Conn) bool {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(maxim.ItemType)
	bd.WriteByte(maxim.Reserved1)
	bd.WriteUint16(maxim.Length)
	bd.WriteUint32(maxim.MaximumLength)

	if err := bd.Send(conn); err != nil {
		return false
	}
	return true
}

func (maxim *maximumSubLength) Read(conn net.Conn) (err error) {
	maxim.ItemType, err = ReadByte(conn)
	if err != nil {
		return err
	}
	return maxim.ReadDynamic(conn)
}

func (maxim *maximumSubLength) ReadDynamic(conn net.Conn) (err error) {
	maxim.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	maxim.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}
	maxim.MaximumLength, err = ReadUint32(conn)
	return
}

// AsyncOperationWindow - AsyncOperationWindow
type AsyncOperationWindow interface {
	Size() uint16
	Read(conn net.Conn) (err error)
	ReadDynamic(conn net.Conn) (err error)
}

type asyncOperationWindow struct {
	ItemType                     byte //0x53
	Reserved1                    byte
	Length                       uint16
	MaxNumberOperationsInvoked   uint16
	MaxNumberOperationsPerformed uint16
}

// NewAsyncOperationWindow - NewAsyncOperationWindow
func NewAsyncOperationWindow() AsyncOperationWindow {
	return &asyncOperationWindow{
		ItemType: 0x53,
	}
}

func (async *asyncOperationWindow) Size() uint16 {
	return async.Length + 4
}

func (async *asyncOperationWindow) Read(conn net.Conn) (err error) {
	async.ItemType, err = ReadByte(conn)
	if err != nil {
		return
	}
	return async.ReadDynamic(conn)
}

func (async *asyncOperationWindow) ReadDynamic(conn net.Conn) (err error) {
	async.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	async.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}
	async.MaxNumberOperationsInvoked, err = ReadUint16(conn)
	if err != nil {
		return
	}
	async.MaxNumberOperationsPerformed, err = ReadUint16(conn)
	return
}

// RoleSelect - RoleSelect
type RoleSelect interface {
	Size() uint16
	Write(conn net.Conn) bool
	Read(conn net.Conn) (err error)
	ReadDynamic(conn net.Conn) (err error)
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

func (scpscu *roleSelect) Write(conn net.Conn) bool {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(scpscu.ItemType)
	bd.WriteByte(scpscu.Reserved1)
	bd.WriteUint16(scpscu.Length)
	bd.WriteUint16(uint16(len(scpscu.uid)))
	bd.Write([]byte(scpscu.uid), len(scpscu.uid))
	bd.WriteByte(scpscu.SCURole)
	bd.WriteByte(scpscu.SCPRole)

	if err := bd.Send(conn); err != nil {
		return false
	}
	return true
}

func (scpscu *roleSelect) Read(conn net.Conn) (err error) {
	scpscu.ItemType, err = ReadByte(conn)
	if err != nil {
		return
	}
	return scpscu.ReadDynamic(conn)
}

func (scpscu *roleSelect) ReadDynamic(conn net.Conn) (err error) {
	scpscu.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	scpscu.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}
	tl, err := ReadUint16(conn)
	if err != nil {
		return
	}

	tuid := make([]byte, tl)
	_, err = conn.Read(tuid)
	if err != nil {
		return
	}

	scpscu.uid = string(tuid)
	scpscu.SCURole, err = ReadByte(conn)
	if err != nil {
		return
	}
	scpscu.SCPRole, err = ReadByte(conn)
	return
}

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
	Write(conn net.Conn) (err error)
	Read(conn net.Conn) (err error)
	ReadDynamic(conn net.Conn) (err error)
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

func (ui *userInformation) Write(conn net.Conn) (err error) {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	ui.Size()
	bd.WriteByte(ui.ItemType)
	bd.WriteByte(ui.Reserved1)
	bd.WriteUint16(ui.Length)

	if err = bd.Send(conn); err == nil {
		ui.MaxSubLength.Write(conn)
		ui.ImpClass.Write(conn)
		ui.ImpVersion.Write(conn)
	}

	return
}

func (ui *userInformation) Read(conn net.Conn) (err error) {
	ui.ItemType, err = ReadByte(conn)
	if err != nil {
		return
	}
	return ui.ReadDynamic(conn)
}

func (ui *userInformation) ReadDynamic(conn net.Conn) (err error) {
	ui.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	ui.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}

	Count := int(ui.Length)
	for Count > 0 {
		TempByte, err := ReadByte(conn)
		if err != nil {
			return err
		}

		switch TempByte {
		case 0x51:
			ui.MaxSubLength.ReadDynamic(conn)
			Count = Count - int(ui.MaxSubLength.Size())
			break
		case 0x52:
			ui.ImpClass.ReadDynamic(conn)
			Count = Count - int(ui.ImpClass.Size())
			break
		case 0x53:
			ui.AsyncOpWindow.ReadDynamic(conn)
			Count = Count - int(ui.AsyncOpWindow.Size())
			break
		case 0x54:
			ui.SCPSCURole.ReadDynamic(conn)
			Count = Count - int(ui.SCPSCURole.Size())
			ui.UserInfoBaggage += uint32(ui.SCPSCURole.Size())
			break
		case 0x55:
			ui.ImpVersion.ReadDynamic(conn)
			Count = Count - int(ui.ImpVersion.Size())
			break
		default:
			conn.Close()
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
