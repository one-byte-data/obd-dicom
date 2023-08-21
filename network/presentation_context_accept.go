package network

import (
	"bufio"
	"log"

	"github.com/one-byte-data/obd-dicom/dictionary/sopclass"
	"github.com/one-byte-data/obd-dicom/dictionary/transfersyntax"
	"github.com/one-byte-data/obd-dicom/media"
)

// PresentationContextAccept accepted presentation context
type PresentationContextAccept interface {
	GetPresentationContextID() byte
	SetPresentationContextID(id byte)
	GetResult() byte
	SetResult(result byte)
	GetTrnSyntax() UIDItem
	Size() uint16
	GetAbstractSyntax() UIDItem
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
	AbsSyntax             uidItem
	TrnSyntax             uidItem
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

func (pc *presentationContextAccept) GetTrnSyntax() UIDItem {
	return &pc.TrnSyntax
}

// Size gets the size of presentation
func (pc *presentationContextAccept) Size() uint16 {
	pc.Length = 4
	pc.Length += pc.TrnSyntax.GetSize()
	return pc.Length + 4
}

func (pc *presentationContextAccept) GetAbstractSyntax() UIDItem {
	return &pc.AbsSyntax
}

func (pc *presentationContextAccept) SetAbstractSyntax(Abst string) {
	pc.AbsSyntax.SetType(0x30)
	pc.AbsSyntax.SetReserved(0x00)
	pc.AbsSyntax.SetUID(Abst)
	pc.AbsSyntax.SetLength(uint16(len(Abst)))
}

func (pc *presentationContextAccept) SetTransferSyntax(Tran string) {
	pc.TrnSyntax.SetType(0x40)
	pc.TrnSyntax.SetReserved(0)
	pc.TrnSyntax.SetUID(Tran)
	pc.TrnSyntax.SetLength(uint16(len(Tran)))
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

	sopName := ""
	tsName := ""
	sopClass := sopclass.GetSOPClassFromUID(pc.GetAbstractSyntax().GetUID())
	if sopClass != nil {
		sopName = sopClass.Name
	}
	transferSyntax := transfersyntax.GetTransferSyntaxFromUID(pc.GetTrnSyntax().GetUID())
	if transferSyntax != nil {
		tsName = transferSyntax.Name
	}
	log.Printf("INFO, ASSOC-AC: \tAccepted PresentationContext: %s - %s\n", pc.GetAbstractSyntax().GetUID(), sopName)
	log.Printf("INFO, ASSOC-AC: \tAccepted TransferSynxtax: %s - %s\n", pc.GetTrnSyntax().GetUID(), tsName)

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
