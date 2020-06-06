package network

import (
	"net"
	"rafael/DICOM/media"
)

type PDV struct {
	Length                uint32
	PresentationContextID byte
	MsgHeader             byte
}

type PDataTF struct {
	ItemType              byte
	Reserved1             byte
	Length                uint32
	Buffer                media.BufData
	BlockSize             uint32
	MsgStatus             uint32
	Endian                uint32
	pdv                   PDV
	PresentationContextID byte
	MsgHeader             byte
}

func (pd *PDataTF) ReadDynamic(conn net.Conn) bool {
	var Count uint32

	if pd.Length == 0 {
		pd.Reserved1 = ReadByte(conn)
		pd.Length = ReadUint32(conn)
	}
	Count = pd.Length
	pd.MsgStatus = 0
	for Count > 0 {
		pd.pdv.Length = ReadUint32(conn)
		pd.pdv.PresentationContextID = ReadByte(conn)
		pd.pdv.MsgHeader = ReadByte(conn)
		buff := make([]byte, pd.pdv.Length-2)
		_, err := conn.Read(buff)
		if err != nil {
			return false
		}
		pd.Buffer.Ms.Write(buff, int(pd.pdv.Length-2))
		Count = Count - pd.pdv.Length - 4
		pd.Length = pd.Length - pd.pdv.Length - 4
		if pd.pdv.MsgHeader&0x02 > 0 {
			pd.MsgStatus = 1
			pd.PresentationContextID = pd.pdv.PresentationContextID
			return true
		}
	}
	if pd.pdv.MsgHeader&0x02 > 0 {
		pd.MsgStatus = 1
	}
	pd.PresentationContextID = pd.pdv.PresentationContextID
	return true
}

func (pd *PDataTF) Write(conn net.Conn) bool {
TotalSize := uint32(pd.Buffer.Ms.Size)
pd.Buffer.Ms.Position=0
if pd.BlockSize==0 {
	pd.BlockSize=4096
	}
SentSize:= uint32(0)
TLength := pd.Length
for (SentSize < TotalSize){
	if (TotalSize-SentSize) < pd.BlockSize {
		pd.BlockSize = TotalSize-SentSize
		}
	if (pd.BlockSize+SentSize) == TotalSize {
			pd.MsgHeader=pd.MsgHeader | 0x02
		} else{
			pd.MsgHeader=pd.MsgHeader & 0x01
		}
	pd.pdv.PresentationContextID = pd.PresentationContextID
	pd.pdv.MsgHeader = pd.MsgHeader
	pd.pdv.Length = pd.BlockSize+2
	pd.Length = pd.pdv.Length+4
	pd.ItemType = 0x04
	pd.Reserved1=0
	var bd media.BufData

	bd.BigEndian = true
	bd.WriteByte(pd.ItemType)
	bd.WriteByte(pd.Reserved1)
	bd.WriteUint32(pd.Length)
	bd.WriteUint32(pd.pdv.Length)
	bd.WriteByte(pd.pdv.PresentationContextID)
	bd.WriteByte(pd.MsgHeader)
	if bd.Send(conn) {
		buff := make([]byte, pd.BlockSize)
		pd.Buffer.Ms.Read(buff, int(pd.BlockSize))
		n, err:=conn.Write(buff)
		if err != nil {
			return false
		}
		if n!=int(pd.BlockSize) {
			return false
		}
	} else {
		return false
		}
	SentSize += pd.BlockSize
	}
	pd.Length=TLength
	return true
}