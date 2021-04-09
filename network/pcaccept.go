package network

import (
	"bufio"
	"log"

	"git.onebytedata.com/odb/go-dicom/media"
)

// PresentationContextAccept accepted presentation context
type PresentationContextAccept interface {
	GetPresentationContextID() byte
	SetPresentationContextID(id byte)
	GetResult() byte
	SetResult(result byte)
	GetTrnSyntax() UIDitem
	Size() uint16
	GetAbstractSyntax() UIDitem
	SetAbstractSyntax(Abst string)
	SetTransferSyntax(Tran string)
	Write(rw *bufio.ReadWriter) (err error)
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type presentationContextAccept struct {
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
func NewPresentationContextAccept() PresentationContextAccept {
	return &presentationContextAccept{
		ItemType:              0x21,
		PresentationContextID: Uniq8(),
		Result:                2,
	}
}

func (pc *presentationContextAccept) GetPresentationContextID() byte {
	return pc.PresentationContextID
}

func (pc *presentationContextAccept) SetPresentationContextID(id byte) {
	pc.PresentationContextID = id
}

func (pc *presentationContextAccept) GetResult() byte {
	return pc.Result
}

func (pc *presentationContextAccept) SetResult(result byte) {
	pc.Result = result
}

func (pc *presentationContextAccept) GetTrnSyntax() UIDitem {
	return pc.TrnSyntax
}

// Size gets the size of presentation
func (pc *presentationContextAccept) Size() uint16 {
	pc.Length = 4
	pc.Length += pc.TrnSyntax.Size()
	return pc.Length + 4
}

func (pc *presentationContextAccept) GetAbstractSyntax() UIDitem {
	return pc.AbsSyntax
}

func (pc *presentationContextAccept) SetAbstractSyntax(Abst string) {
	pc.AbsSyntax.ItemType = 0x30
	pc.AbsSyntax.Reserved1 = 0x00
	pc.AbsSyntax.UIDName = Abst
	pc.AbsSyntax.Length = uint16(len(Abst))
}

func (pc *presentationContextAccept) SetTransferSyntax(Tran string) {
	pc.TrnSyntax.ItemType = 0x40
	pc.TrnSyntax.Reserved1 = 0
	pc.TrnSyntax.UIDName = Tran
	pc.TrnSyntax.Length = uint16(len(Tran))
}

func (pc *presentationContextAccept) Write(rw *bufio.ReadWriter) (err error) {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	pc.Size()
	bd.WriteByte(pc.ItemType)
	bd.WriteByte(pc.Reserved1)
	bd.WriteUint16(pc.Length)
	bd.WriteByte(pc.PresentationContextID)
	bd.WriteByte(pc.Reserved2)
	bd.WriteByte(pc.Result)
	bd.WriteByte(pc.Reserved4)

	log.Printf("INFO, ASSOC-AC: \tAccepted Presentation Context %s\n", pc.GetAbstractSyntax().UIDName)
	log.Printf("INFO, ASSOC-AC: \tAccepted Transfer Synxtax %s\n", pc.GetTrnSyntax().UIDName)

	if err = bd.Send(rw); err == nil {
		return pc.TrnSyntax.Write(rw)
	}
	return
}

func (pc *presentationContextAccept) Read(ms media.MemoryStream) (err error) {
	pc.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return pc.ReadDynamic(ms)
}

func (pc *presentationContextAccept) ReadDynamic(ms media.MemoryStream) (err error) {
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
	pc.Result, err = ms.GetByte()
	if err != nil {
		return err
	}
	pc.Reserved4, err = ms.GetByte()
	if err != nil {
		return err
	}

	return pc.TrnSyntax.Read(ms)
}
