package network

import (
	"bufio"
	"errors"

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
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (error)
	ReadDynamic(ms media.MemoryStream) (error)
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
	for _, TrnSyntax := range pc.TrnSyntaxs {
		pc.Length += TrnSyntax.Size()
	}
	return pc.Length + 4
}

func (pc *presentationContext) Write(rw *bufio.ReadWriter) error {
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
	err := bd.Send(rw)
	if err != nil {
		return err
	}
	err = pc.AbsSyntax.Write(rw)
	if err != nil {
		return err
	}
	for _, TrnSyntax := range pc.TrnSyntaxs {
		err := TrnSyntax.Write(rw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pc *presentationContext) Read(ms media.MemoryStream) (err error) {
	pc.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return pc.ReadDynamic(ms)
}

func (pc *presentationContext) ReadDynamic(ms media.MemoryStream) (err error) {
	pc.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	pc.Length, err = ms.GetUint16()
	if err != nil {
		return err
	}
	pc.PresentationContextID, err = ms.GetByte()
	if err != nil {
		return err
	}
	pc.Reserved2, err = ms.GetByte()
	if err != nil {
		return err
	}
	pc.Reserved3, err = ms.GetByte()
	if err != nil {
		return err
	}
	pc.Reserved4, err = ms.GetByte()
	if err != nil {
		return err
	}

	pc.AbsSyntax.Read(ms)

	Count := pc.Length - 4 - pc.AbsSyntax.Size()
	for Count > 0 {
		var TrnSyntax UIDitem
		TrnSyntax.Read(ms)
		Count = Count - TrnSyntax.Size()
		if TrnSyntax.Size() > 0 {
			pc.TrnSyntaxs = append(pc.TrnSyntaxs, TrnSyntax)
		}
	}

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, pc::ReadDynamic, Count is not zero")
}