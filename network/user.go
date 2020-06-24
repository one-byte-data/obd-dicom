package network

import (
	"net"

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
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	bd.WriteByte(maxim.ItemType)
	bd.WriteByte(maxim.Reserved1)
	bd.WriteUint16(maxim.Length)
	bd.WriteUint32(maxim.MaximumLength)
	flag = bd.Send(conn)
	return flag
}

func (maxim *MaximumSubLength) Read(conn net.Conn) bool {
	maxim.ItemType = ReadByte(conn)
	return maxim.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (maxim *MaximumSubLength) ReadDynamic(conn net.Conn) bool {
	maxim.Reserved1 = ReadByte(conn)
	maxim.Length = ReadUint16(conn)
	maxim.MaximumLength = ReadUint32(conn)
	return true
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

func (async *AsyncOperationWindow) Read(conn net.Conn) bool {
	async.ItemType = ReadByte(conn)
	return async.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (async *AsyncOperationWindow) ReadDynamic(conn net.Conn) bool {
	async.Reserved1 = ReadByte(conn)
	async.Length = ReadUint16(conn)
	async.MaxNumberOperationsInvoked = ReadUint16(conn)
	async.MaxNumberOperationsPerformed = ReadUint16(conn)
	return true
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
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	bd.WriteByte(scpscu.ItemType)
	bd.WriteByte(scpscu.Reserved1)
	bd.WriteUint16(scpscu.Length)
	bd.WriteUint16(uint16(len(scpscu.uid)))
	bd.Ms.Write([]byte(scpscu.uid), len(scpscu.uid))
	bd.WriteByte(scpscu.SCURole)
	bd.WriteByte(scpscu.SCPRole)
	flag = bd.Send(conn)
	return flag
}

func (scpscu *SCPSCURoleSelect) Read(conn net.Conn) bool {
	scpscu.ItemType = ReadByte(conn)
	return scpscu.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (scpscu *SCPSCURoleSelect) ReadDynamic(conn net.Conn) bool {
	scpscu.Reserved1 = ReadByte(conn)
	scpscu.Length = ReadUint16(conn)
	tl := ReadUint16(conn)
	tuid := make([]byte, tl)
	conn.Read(tuid)
	scpscu.uid = string(tuid)
	scpscu.SCURole = ReadByte(conn)
	scpscu.SCPRole = ReadByte(conn)
	return true
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

func (ui *UserInformation) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	ui.Size()
	bd.WriteByte(ui.ItemType)
	bd.WriteByte(ui.Reserved1)
	bd.WriteUint16(ui.Length)
	if bd.Send(conn) {
		ui.MaxSubLength.Write(conn)
		ui.ImpClass.Write(conn)
		ui.ImpVersion.Write(conn)
		flag = true
	}
	return flag
}

func (ui *UserInformation) Read(conn net.Conn) bool {
	ui.ItemType = ReadByte(conn)
	return ui.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (ui *UserInformation) ReadDynamic(conn net.Conn) bool {
	ui.Reserved1 = ReadByte(conn)
	ui.Length = ReadUint16(conn)
	var Count int
	Count = int(ui.Length)
	for Count > 0 {
		TempByte := ReadByte(conn)
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
			break
		}
	}
	if Count == 0 {
		return true
	}
	return (false)
}
