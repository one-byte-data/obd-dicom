package network

import (
	"encoding/binary"
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// ReadByte reads a byte
func ReadByte(conn net.Conn) byte {
	c := make([]byte, 1)
	_, err := conn.Read(c)
	if err != nil {
		return 0
	}
	return c[0]
}

// ReadUint16 read unsigned int
func ReadUint16(conn net.Conn) uint16 {
	var val uint16
	c := make([]byte, 2)
	_, err := conn.Read(c)
	if err != nil {
		return 0
	}
	val = binary.BigEndian.Uint16(c)
	return val
}

// ReadUint32 read unsigned int
func ReadUint32(conn net.Conn) uint32 {
	var val uint32
	c := make([]byte, 4)
	_, err := conn.Read(c)
	if err != nil {
		return 0
	}
	val = binary.BigEndian.Uint32(c)
	return val
}

// PresentationContextAccept accepted presentation context
type PresentationContextAccept struct {
	ItemType              byte //0x21
	Reserved1             byte
	Length                uint16
	PresentationContextID byte
	Reserved2             byte
	Result                byte
	Reserved4             byte
	AbsSyntax             UIDitem
	TrnSyntax             UIDitem
}

// NewPresentationContextAccept creates a PresentationContextAccept
func NewPresentationContextAccept() *PresentationContextAccept {
	pc := &PresentationContextAccept{}
	pc.ItemType = 0x21
	pc.PresentationContextID = Uniq8()
	pc.Result = 2
	return pc
}

// Size gets the size of presentation
func (pc *PresentationContextAccept) Size() uint16 {
	pc.Length = 4
	pc.Length += pc.TrnSyntax.Size()
	return pc.Length + 4
}

// SetAbstractSyntax sets abstrct syntax
func (pc *PresentationContextAccept) SetAbstractSyntax(Abst string) {
	pc.AbsSyntax.ItemType = 0x30
	pc.AbsSyntax.Reserved1 = 0x00
	pc.AbsSyntax.UIDName = Abst
	pc.AbsSyntax.Length = uint16(len(Abst))
}

// SetTransferSyntax sets the transfer syntax
func (pc *PresentationContextAccept) SetTransferSyntax(Tran string) {
	pc.TrnSyntax.ItemType = 0x40
	pc.TrnSyntax.Reserved1 = 0
	pc.TrnSyntax.UIDName = Tran
	pc.TrnSyntax.Length = uint16(len(Tran))
}

func (pc *PresentationContextAccept) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	pc.Size()
	bd.WriteByte(pc.ItemType)
	bd.WriteByte(pc.Reserved1)
	bd.WriteUint16(pc.Length)
	bd.WriteByte(pc.PresentationContextID)
	bd.WriteByte(pc.Reserved2)
	bd.WriteByte(pc.Result)
	bd.WriteByte(pc.Reserved4)
	if bd.Send(conn) {
		flag = pc.TrnSyntax.Write(conn)
	}
	return flag
}
func (pc *PresentationContextAccept) Read(conn net.Conn) bool {
	pc.ItemType = ReadByte(conn)
	return pc.ReadDynamic(conn)
}

// ReadDynamic ReadDynamic
func (pc *PresentationContextAccept) ReadDynamic(conn net.Conn) bool {
	pc.Reserved1 = ReadByte(conn)
	pc.Length = ReadUint16(conn)
	pc.PresentationContextID = ReadByte(conn)
	pc.Reserved2 = ReadByte(conn)
	pc.Result = ReadByte(conn)
	pc.Reserved4 = ReadByte(conn)
	pc.TrnSyntax.Read(conn)
	return true
}

// AAssociationAC AAssociationAC
type AAssociationAC struct {
	ItemType           byte // 0x02
	Reserved1          byte
	Length             uint32
	ProtocolVersion    uint16 // 0x01
	Reserved2          uint16
	CallingApTitle     [16]byte // 16 bytes transfered
	CalledApTitle      [16]byte // 16 bytes transfered
	Reserved3          [32]byte
	AppContext         UIDitem
	PresContextAccepts []PresentationContextAccept
	UserInfo           UserInformation
}

// NewAAssociationAC NewAAssociationAC
func NewAAssociationAC() *AAssociationAC {
	aaac := &AAssociationAC{}
	aaac.ItemType = 0x02
	aaac.Reserved1 = 0x00
	aaac.ProtocolVersion = 0x01
	aaac.Reserved2 = 0x00
	aaac.AppContext.ItemType = 0x10
	aaac.AppContext.Reserved1 = 0x00
	aaac.AppContext.UIDName = "1.2.840.10008.3.1.1.1"
	aaac.AppContext.Length = uint16(len(aaac.AppContext.UIDName))
	return aaac
}

// Size size of association
func (aaac *AAssociationAC) Size() uint32 {
	aaac.Length = 4 + 16 + 16 + 32
	aaac.Length += uint32(aaac.AppContext.Size())

	for i := 0; i < len(aaac.PresContextAccepts); i++ {
		PresContextAccept := aaac.PresContextAccepts[i]
		aaac.Length += uint32(PresContextAccept.Size())
	}
	aaac.Length += uint32(aaac.UserInfo.Size())
	return aaac.Length + 6
}

// SetUserInformation SetUserInformation
func (aaac *AAssociationAC) SetUserInformation(UserInfo UserInformation) {
	aaac.UserInfo = UserInfo
}

func (aaac *AAssociationAC) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	aaac.Size()
	bd.WriteByte(aaac.ItemType)
	bd.WriteByte(aaac.Reserved1)
	bd.WriteUint32(aaac.Length)
	bd.WriteUint16(aaac.ProtocolVersion)
	bd.WriteUint16(aaac.Reserved2)
	bd.Ms.Write(aaac.CalledApTitle[:], 16)
	bd.Ms.Write(aaac.CallingApTitle[:], 16)
	bd.Ms.Write(aaac.Reserved3[:], 32)

	if bd.Send(conn) {
		flag = aaac.AppContext.Write(conn)
		for i := 0; i < len(aaac.PresContextAccepts); i++ {
			PresContextAccept := aaac.PresContextAccepts[i]
			PresContextAccept.Write(conn)
		}
		flag = aaac.UserInfo.Write(conn)
	}
	return flag
}

func (aaac *AAssociationAC) Read(conn net.Conn) bool {
	aaac.ItemType = ReadByte(conn)
	return aaac.ReadDynamic(conn)
}

// ReadDynamic ReadDynamic
func (aaac *AAssociationAC) ReadDynamic(conn net.Conn) bool {
	aaac.Reserved1 = ReadByte(conn)
	aaac.Length = ReadUint32(conn)
	aaac.ProtocolVersion = ReadUint16(conn)
	aaac.Reserved2 = ReadUint16(conn)
	conn.Read(aaac.CalledApTitle[:])
	conn.Read(aaac.CallingApTitle[:])
	conn.Read(aaac.Reserved3[:])

	var Count int
	Count = int(aaac.Length - 4 - 16 - 16 - 32)
	for Count > 0 {
		TempByte := ReadByte(conn)
		switch TempByte {
		case 0x10:
			aaac.AppContext.ReadDynamic(conn)
			Count = Count - int(aaac.AppContext.Size())
			break
		case 0x21:
			PresContextAccept := NewPresentationContextAccept()
			PresContextAccept.ReadDynamic(conn)
			Count = Count - int(PresContextAccept.Size())
			aaac.PresContextAccepts = append(aaac.PresContextAccepts, *PresContextAccept)
			break
		case 0x50: // User Information
			aaac.UserInfo.ReadDynamic(conn)
			Count = Count - int(aaac.UserInfo.Size())
			break
		default:
			conn.Close()
			Count = -1
		}
	}
	if Count == 0 {
		return true
	}
	return (false)
}
