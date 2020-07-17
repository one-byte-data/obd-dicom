package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// PresentationContext - PresentationContext
type PresentationContext interface {
	GetPresentationContextID() byte
	SetPresentationContextID(id byte)
	GetAbstractSyntax() UIDitem
	SetAbstractSyntax(Abst string)
	AddTransferSyntax(Tran string)
	GetTransferSyntaxes() []UIDitem
	Size() uint16
	Write(conn net.Conn) error
	Read(conn net.Conn) (bool, error)
	ReadDynamic(conn net.Conn) (bool, error)
}

type presentationContext struct {
	ItemType              byte //0x20
	Reserved1             byte
	Length                uint16
	PresentationContextID byte
	Reserved2             byte
	Reserved3             byte
	Reserved4             byte
	AbsSyntax             UIDitem
	TrnSyntaxs            []UIDitem
}

// NewPresentationContext - NewPresentationContext
func NewPresentationContext() PresentationContext {
	return &presentationContext{
		ItemType:              0x20,
		PresentationContextID: Uniq8odd(),
	}
}

func (pc *presentationContext) GetPresentationContextID() byte {
	return pc.PresentationContextID
}

func (pc *presentationContext) SetPresentationContextID(id byte) {
	pc.PresentationContextID = id
}

func (pc *presentationContext) GetAbstractSyntax() UIDitem {
	return pc.AbsSyntax
}

func (pc *presentationContext) SetAbstractSyntax(Abst string) {
	pc.AbsSyntax.ItemType = 0x30
	pc.AbsSyntax.Reserved1 = 0x00
	pc.AbsSyntax.UIDName = Abst
	pc.AbsSyntax.Length = uint16(len(Abst))
}

func (pc *presentationContext) AddTransferSyntax(Tran string) {
	TrnSyntax := NewUIDitem(Tran, 0x40)
	pc.TrnSyntaxs = append(pc.TrnSyntaxs, *TrnSyntax)
}

func (pc *presentationContext) GetTransferSyntaxes() []UIDitem {
	return pc.TrnSyntaxs
}

func (pc *presentationContext) Size() uint16 {
	pc.Length = 4
	pc.Length += pc.AbsSyntax.Size()
	for i := 0; i < len(pc.TrnSyntaxs); i++ {
		TrnSyntax := pc.TrnSyntaxs[i]
		pc.Length += TrnSyntax.Size()
	}
	return pc.Length + 4
}

func (pc *presentationContext) Write(conn net.Conn) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	pc.Size()
	bd.WriteByte(pc.ItemType)
	bd.WriteByte(pc.Reserved1)
	bd.WriteUint16(pc.Length)
	bd.WriteByte(pc.PresentationContextID)
	bd.WriteByte(pc.Reserved2)
	bd.WriteByte(pc.Reserved3)
	bd.WriteByte(pc.Reserved4)
	err := bd.Send(conn)
	if err != nil {
		return err
	}
	err = pc.AbsSyntax.Write(conn)
	if err != nil {
		return err
	}
	for i := 0; i < len(pc.TrnSyntaxs); i++ {
		TrnSyntax := pc.TrnSyntaxs[i]
		err := TrnSyntax.Write(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pc *presentationContext) Read(conn net.Conn) (bool, error) {
	var err error
	pc.ItemType, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	return pc.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (pc *presentationContext) ReadDynamic(conn net.Conn) (bool, error) {
	var err error
	pc.Reserved1, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	pc.Length, err = ReadUint16(conn)
	if err != nil {
		return false, err
	}
	pc.PresentationContextID, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	pc.Reserved2, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	pc.Reserved3, err = ReadByte(conn)
	if err != nil {
		return false, err
	}
	pc.Reserved4, err = ReadByte(conn)
	if err != nil {
		return false, err
	}

	pc.AbsSyntax.Read(conn)
	Count := pc.Length - 4 - pc.AbsSyntax.Size()
	for Count > 0 {
		var TrnSyntax UIDitem
		TrnSyntax.Read(conn)
		Count = Count - TrnSyntax.Size()
		if TrnSyntax.Size() > 0 {
			pc.TrnSyntaxs = append(pc.TrnSyntaxs, TrnSyntax)
		}
	}
	if Count == 0 {
		return true, nil
	}
	return false, errors.New("ERROR, pc::ReadDynamic, Count is not zero")
}

// AAssociationRQ - AAssociationRQ
type AAssociationRQ interface {
	GetAppContext() UIDitem
	SetAppContext(context UIDitem)
	GetCallingAE() string
	SetCallingAE(AET string)
	GetCalledAE() string
	SetCalledAE(AET string)
	GetPresContexts() []PresentationContext
	GetUserInformation() UserInformation
	SetUserInformation(userInfo UserInformation)
	GetMaxSubLength() uint32
	SetMaxSubLength(length uint32)
	GetImpClass() UIDitem
	SetImpClassUID(uid string)
	SetImpVersionName(name string)
	Size() uint32
	Write(conn net.Conn) error
	Read(conn net.Conn) error
	ReadDynamic(conn net.Conn) error
	AddPresContexts(presentationContext PresentationContext)
}

type aassociationRQ struct {
	ItemType        byte // 0x01
	Reserved1       byte
	Length          uint32
	ProtocolVersion uint16 // 0x01
	Reserved2       uint16
	CallingAE       [16]byte // 16 bytes transfered
	CalledAE        [16]byte // 16 bytes transfered
	Reserved3       [32]byte
	AppContext      UIDitem
	PresContexts    []PresentationContext
	UserInfo        UserInformation
}

// NewAAssociationRQ - NewAAssociationRQ
func NewAAssociationRQ() AAssociationRQ {
	return &aassociationRQ{
		ItemType:        0x01,
		Reserved1:       0x00,
		ProtocolVersion: 0x01,
		Reserved2:       0x00,
		AppContext: UIDitem{
			ItemType:  0x10,
			Reserved1: 0x00,
			UIDName:   "1.2.840.10008.3.1.1.1",
			Length:    uint16(len("1.2.840.10008.3.1.1.1")),
		},
		PresContexts: make([]PresentationContext, 0),
		UserInfo:     *NewUserInformation(),
	}
}

func (aarq *aassociationRQ) GetAppContext() UIDitem {
	return aarq.AppContext
}

func (aarq *aassociationRQ) SetAppContext(context UIDitem) {
	aarq.AppContext = context
}

func (aarq *aassociationRQ) GetCallingAE() string {
	return fmt.Sprintf("%s", aarq.CallingAE)
}

func (aarq *aassociationRQ) SetCallingAE(AET string) {
	copy(aarq.CallingAE[:], AET)
}

func (aarq *aassociationRQ) GetCalledAE() string {
	return fmt.Sprintf("%s", aarq.CalledAE)
}

func (aarq *aassociationRQ) SetCalledAE(AET string) {
	copy(aarq.CalledAE[:], AET)
}

func (aarq *aassociationRQ) GetPresContexts() []PresentationContext {
	return aarq.PresContexts
}

func (aarq *aassociationRQ) GetUserInformation() UserInformation {
	return aarq.UserInfo
}

func (aarq *aassociationRQ) SetUserInformation(userInfo UserInformation) {
	aarq.UserInfo = userInfo
}

func (aarq *aassociationRQ) GetMaxSubLength() uint32 {
	return aarq.UserInfo.MaxSubLength.MaximumLength
}

func (aarq *aassociationRQ) SetMaxSubLength(length uint32) {
	aarq.UserInfo.MaxSubLength.MaximumLength = length
}

func (aarq *aassociationRQ) GetImpClass() UIDitem {
	return aarq.UserInfo.ImpClass
}

func (aarq *aassociationRQ)	SetImpClassUID(uid string) {
	aarq.UserInfo.SetImpClassUID(uid)
}

func (aarq *aassociationRQ) SetImpVersionName(name string) {
	aarq.UserInfo.SetImpVersionName(name)
}

func (aarq *aassociationRQ) Size() uint32 {
	aarq.Length = 4 + 16 + 16 + 32
	aarq.Length += uint32(aarq.AppContext.Size())

	for i := 0; i < len(aarq.PresContexts); i++ {
		PresContext := aarq.PresContexts[i]
		aarq.Length += uint32(PresContext.Size())
	}
	aarq.Length += uint32(aarq.UserInfo.Size())
	return aarq.Length + 6
}

func (aarq *aassociationRQ) Write(conn net.Conn) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	aarq.Size()
	bd.WriteByte(aarq.ItemType)
	bd.WriteByte(aarq.Reserved1)
	bd.WriteUint32(aarq.Length)
	bd.WriteUint16(aarq.ProtocolVersion)
	bd.WriteUint16(aarq.Reserved2)
	bd.Write(aarq.CalledAE[:], 16)
	bd.Write(aarq.CallingAE[:], 16)
	bd.Write(aarq.Reserved3[:], 32)

	if err := bd.Send(conn); err == nil {
		err = aarq.AppContext.Write(conn)
		if err != nil {
			return err
		}
		for i := 0; i < len(aarq.PresContexts); i++ {
			PresContext := aarq.PresContexts[i]
			PresContext.Write(conn)
		}
		aarq.UserInfo.Write(conn)
	}
	return nil
}

func (aarq *aassociationRQ) Read(conn net.Conn) error {
	var err error
	aarq.ItemType, err = ReadByte(conn)
	if err != nil {
		return err
	}
	return aarq.ReadDynamic(conn)
}

func (aarq *aassociationRQ) ReadDynamic(conn net.Conn) error {
	var err error
	aarq.Reserved1, err = ReadByte(conn)
	if err != nil {
		return err
	}
	aarq.Length, err = ReadUint32(conn)
	if err != nil {
		return err
	}
	aarq.ProtocolVersion, err = ReadUint16(conn)
	if err != nil {
		return err
	}
	aarq.Reserved2, err = ReadUint16(conn)
	if err != nil {
		return err
	}

	conn.Read(aarq.CalledAE[:])
	conn.Read(aarq.CallingAE[:])
	conn.Read(aarq.Reserved3[:])

	var Count int
	Count = int(aarq.Length - 4 - 16 - 16 - 32)
	for Count > 0 {
		TempByte, err := ReadByte(conn)
		if err != nil {
			return err
		}

		switch TempByte {
		case 0x10:
			aarq.AppContext.ItemType = TempByte
			aarq.AppContext.ReadDynamic(conn)
			Count = Count - int(aarq.AppContext.Size())
			break
		case 0x20:
			PresContext := NewPresentationContext()
			PresContext.ReadDynamic(conn)
			Count = Count - int(PresContext.Size())
			aarq.PresContexts = append(aarq.PresContexts, PresContext)
			break
		case 0x50: // User Information
			aarq.UserInfo.ItemType = TempByte
			aarq.UserInfo.ReadDynamic(conn)
			Count = Count - int(aarq.UserInfo.Size())
			break
		default:
			log.Println("ERROR, aarq::ReadDynamic, unknown Item, " + strconv.Itoa(int(TempByte)))
			Count = -1
		}
	}
	if Count == 0 {
		return nil
	}
	return errors.New("ERROR, aarq::ReadDynamic, Count is not zero")
}

func (aarq *aassociationRQ) AddPresContexts(presentationContext PresentationContext) {
	aarq.PresContexts = append(aarq.PresContexts, presentationContext)
}
