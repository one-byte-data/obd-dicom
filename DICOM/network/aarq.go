package network

import (
	"net"
	"rafael/DICOM/media"
)

type PresentationContext struct {
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

func NewPresentationContext() *PresentationContext {
	pc := &PresentationContext{}
	pc.ItemType = 0x20
	pc.PresentationContextID = uniq8odd()
	return pc
}

func (pc *PresentationContext) Size() uint16 {
	pc.Length = 4
	pc.Length += pc.AbsSyntax.Size()
	for i := 0; i < len(pc.TrnSyntaxs); i++ {
		TrnSyntax := pc.TrnSyntaxs[i]
		pc.Length += TrnSyntax.Size()
	}
	return pc.Length + 4
}

func (pc *PresentationContext) SetAbstractSyntax(Abst string) {
	pc.AbsSyntax.ItemType = 0x30
	pc.AbsSyntax.Reserved1 = 0x00
	pc.AbsSyntax.UIDName = Abst
	pc.AbsSyntax.Length = uint16(len(Abst))
}

func (pc *PresentationContext) AddTransferSyntax(Tran string) {
	TrnSyntax := NewUIDitem(Tran, 0x40)
	pc.TrnSyntaxs = append(pc.TrnSyntaxs, *TrnSyntax)
}

func (pc *PresentationContext) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	pc.Size()
	bd.WriteByte(pc.ItemType)
	bd.WriteByte(pc.Reserved1)
	bd.WriteUint16(pc.Length)
	bd.WriteByte(pc.PresentationContextID)
	bd.WriteByte(pc.Reserved2)
	bd.WriteByte(pc.Reserved3)
	bd.WriteByte(pc.Reserved4)
	if bd.Send(conn) {
		flag = pc.AbsSyntax.Write(conn)
		for i := 0; i < len(pc.TrnSyntaxs); i++ {
			TrnSyntax := pc.TrnSyntaxs[i]
			TrnSyntax.Write(conn)
		}
	}
	return flag
}

type AAssociationRQ struct {
	ItemType        byte // 0x01
	Reserved1       byte
	Length          uint32
	ProtocolVersion uint16 // 0x01
	Reserved2       uint16
	CallingApTitle  [16]byte // 16 bytes transfered
	CalledApTitle   [16]byte // 16 bytes transfered
	Reserved3       [32]byte
	AppContext      UIDitem
	PresContexts    []PresentationContext
	UserInfo UserInformation
}

func NewAAAssociationRQ() *AAssociationRQ {
	aarq := &AAssociationRQ{}
	aarq.ItemType = 0x01
	aarq.Reserved1 = 0x00
	aarq.ProtocolVersion = 0x01
	aarq.Reserved2 = 0x00
	aarq.AppContext.ItemType = 0x10
	aarq.AppContext.Reserved1 = 0x00
	aarq.AppContext.UIDName = "1.2.840.10008.3.1.1.1"
	aarq.AppContext.Length = uint16(len(aarq.AppContext.UIDName))
	return aarq
}

func (aarq *AAssociationRQ) Size() uint32 {
	aarq.Length = 4 + 16 + 16 + 32
	aarq.Length += uint32(aarq.AppContext.Size())

	for i := 0; i < len(aarq.PresContexts); i++ {
		PresContext := aarq.PresContexts[i]
		aarq.Length += uint32(PresContext.Size())
	}
	aarq.Length += uint32(aarq.UserInfo.Size())
	return aarq.Length + 6
}

func (aarq *AAssociationRQ) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	aarq.Size()
	bd.WriteByte(aarq.ItemType)
	bd.WriteByte(aarq.Reserved1)
	bd.WriteUint32(aarq.Length)
	bd.WriteUint16(aarq.ProtocolVersion)
	bd.WriteUint16(aarq.Reserved2)
	bd.Ms.Write(aarq.CalledApTitle[:], 16)
	bd.Ms.Write(aarq.CallingApTitle[:], 16)
	bd.Ms.Write(aarq.Reserved3[:], 32)

	if bd.Send(conn) {
		flag = aarq.AppContext.Write(conn)
		for i := 0; i < len(aarq.PresContexts); i++ {
			PresContext := aarq.PresContexts[i]
			PresContext.Write(conn)
		}
		aarq.UserInfo.Write(conn)
	}
	return flag
}
