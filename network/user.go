package network

import (
	"errors"
	"log"
	"net"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// MaximumSubLength - MaximumSubLength
type MaximumSubLength struct {
	ItemType      byte //0x51
	Reserved1     byte
	Length        uint16
	MaximumLength uint32
}

// NewMaximumSubLength - NewMaximumSubLength
func NewMaximumSubLength() *MaximumSubLength {
	return &MaximumSubLength{
		ItemType: 0x51,
		Length:   4,
	}
}

// Size - Size
func (maxim *MaximumSubLength) Size() uint16 {
	return maxim.Length + 4
}

func (maxim *MaximumSubLength) Write(conn net.Conn) bool {
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

func (maxim *MaximumSubLength) Read(conn net.Conn) (bool, error) {
	var err error
	maxim.ItemType, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	return maxim.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (maxim *MaximumSubLength) ReadDynamic(conn net.Conn) (bool, error) {
	var err error
	maxim.Reserved1, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	maxim.Length, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	maxim.MaximumLength, err = ReadUint32(conn)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AsyncOperationWindow - AsyncOperationWindow
type AsyncOperationWindow struct {
	ItemType                     byte //0x53
	Reserved1                    byte
	Length                       uint16
	MaxNumberOperationsInvoked   uint16
	MaxNumberOperationsPerformed uint16
}

// NewAsyncOperationWindow - NewAsyncOperationWindow
func NewAsyncOperationWindow() *AsyncOperationWindow {
	return &AsyncOperationWindow{
		ItemType: 0x53,
	}
}

// Size - Size
func (async *AsyncOperationWindow) Size() uint16 {
	return async.Length + 4
}

func (async *AsyncOperationWindow) Read(conn net.Conn) (bool, error) {
	var err error
	async.ItemType, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	return async.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (async *AsyncOperationWindow) ReadDynamic(conn net.Conn) (bool, error) {
	var err error
	async.Reserved1, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	async.Length, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	async.MaxNumberOperationsInvoked, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	async.MaxNumberOperationsPerformed, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SCPSCURoleSelect - SCPSCURoleSelect
type SCPSCURoleSelect struct {
	ItemType  byte //0x54
	Reserved1 byte
	Length    uint16
	SCURole   byte
	SCPRole   byte
	uid       string
}

// NewSCPSCURoleSelect - NewSCPSCURoleSelect
func NewSCPSCURoleSelect() *SCPSCURoleSelect {
	return &SCPSCURoleSelect{
		ItemType: 0x54,
	}
}

// Size - Size
func (scpscu *SCPSCURoleSelect) Size() uint16 {
	return scpscu.Length + 4
}

func (scpscu *SCPSCURoleSelect) Write(conn net.Conn) bool {
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

func (scpscu *SCPSCURoleSelect) Read(conn net.Conn) (bool, error) {
	var err error
	scpscu.ItemType, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	return scpscu.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (scpscu *SCPSCURoleSelect) ReadDynamic(conn net.Conn) (bool, error) {
	var err error
	scpscu.Reserved1, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	scpscu.Length, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	tl, err := ReadUint16(conn)
	if err != nil {
		return false, err
	}

	tuid := make([]byte, tl)
	_, err = conn.Read(tuid)
	if err != nil {
		return false, err
	}

	scpscu.uid = string(tuid)
	scpscu.SCURole, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	scpscu.SCPRole, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UserInformation - UserInformation
type UserInformation struct {
	ItemType        byte //0x50
	Reserved1       byte
	Length          uint16
	UserInfoBaggage uint32
	MaxSubLength    MaximumSubLength
	AsyncOpWindow   AsyncOperationWindow
	SCPSCURole      SCPSCURoleSelect
	ImpClass        UIDitem
	ImpVersion      UIDitem
}

// NewUserInformation - NewUserInformation
func NewUserInformation() *UserInformation {
	return &UserInformation{
		ItemType:      0x50,
		MaxSubLength:  *NewMaximumSubLength(),
		AsyncOpWindow: *NewAsyncOperationWindow(),
		SCPSCURole:    *NewSCPSCURoleSelect(),
	}
}

// Size - Size
func (ui *UserInformation) Size() uint16 {
	ui.Length = ui.MaxSubLength.Size()
	ui.Length += ui.ImpClass.Size()
	ui.Length += ui.ImpVersion.Size()
	return ui.Length + 4
}

// SetImpClassUID - SetImpClassUID
func (ui *UserInformation) SetImpClassUID(name string) {
	ui.ImpClass.ItemType = 0x52
	ui.ImpClass.Reserved1 = 0x00
	ui.ImpClass.UIDName = name
	ui.ImpClass.Length = uint16(len(name))
}

// SetImpVersionName - SetImpVersionName
func (ui *UserInformation) SetImpVersionName(name string) {
	ui.ImpVersion.ItemType = 0x55
	ui.ImpVersion.Reserved1 = 0x00
	ui.ImpVersion.UIDName = name
	ui.ImpVersion.Length = uint16(len(name))
}

func (ui *UserInformation) Write(conn net.Conn) (err error) {
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

func (ui *UserInformation) Read(conn net.Conn) (err error) {
	ui.ItemType, err = ReadByte(conn)
	if err != nil {
		return
	}
	return ui.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (ui *UserInformation) ReadDynamic(conn net.Conn) (err error) {
	ui.Reserved1, err = ReadByte(conn)
	if err != nil {
		return
	}
	ui.Length, err = ReadUint16(conn)
	if err != nil {
		return
	}

	var Count int
	Count = int(ui.Length)
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
			log.Println("ERROR, user::ReadDynamic, unknown TempByte: " + strconv.Itoa(int(TempByte)))
			break
		}
	}

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, user::ReadDynamic, Count is not zero")
}
