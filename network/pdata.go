package network

import (
	"errors"
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// PDV - PDV
type PDV struct {
	Length                uint32
	PresentationContextID byte
	MsgHeader             byte
}

// PDataTF - PDataTF
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

// ReadDynamic - ReadDynamic
func (pd *PDataTF) ReadDynamic(conn net.Conn) (err error) {
	var Count uint32

	if pd.Length == 0 {
		pd.Reserved1, err = ReadByte(conn)
		if err != nil {
			return
		}
		pd.Length, err = ReadUint32(conn)
		if err != nil {
			return
		}
	}

	Count = pd.Length
	pd.MsgStatus = 0

	for Count > 0 {
		pd.pdv.Length, err = ReadUint32(conn)
		if err != nil {
			return
		}
		pd.pdv.PresentationContextID, err = ReadByte(conn)
		if err != nil {
			return
		}
		pd.pdv.MsgHeader, err = ReadByte(conn)
		if err != nil {
			return
		}

		buff := make([]byte, pd.pdv.Length-2)
		_, err := conn.Read(buff)
		if err != nil {
			return errors.New("ERROR, pdata::ReadDynamic, " + err.Error())
		}

		pd.Buffer.Write(buff, int(pd.pdv.Length-2))
		Count = Count - pd.pdv.Length - 4
		pd.Length = pd.Length - pd.pdv.Length - 4

		if pd.pdv.MsgHeader&0x02 > 0 {
			pd.MsgStatus = 1
			pd.PresentationContextID = pd.pdv.PresentationContextID
			return nil
		}
	}

	if pd.pdv.MsgHeader&0x02 > 0 {
		pd.MsgStatus = 1
	}

	pd.PresentationContextID = pd.pdv.PresentationContextID
	return nil
}

func (pd *PDataTF) Write(conn net.Conn) error {
	TotalSize := uint32(pd.Buffer.GetSize())
	pd.Buffer.SetPosition(0)
	if pd.BlockSize == 0 {
		pd.BlockSize = 4096
	}

	SentSize := uint32(0)
	TLength := pd.Length

	for SentSize < TotalSize {
		if (TotalSize - SentSize) < pd.BlockSize {
			pd.BlockSize = TotalSize - SentSize
		}
		if (pd.BlockSize + SentSize) == TotalSize {
			pd.MsgHeader = pd.MsgHeader | 0x02
		} else {
			pd.MsgHeader = pd.MsgHeader & 0x01
		}

		pd.pdv.PresentationContextID = pd.PresentationContextID
		pd.pdv.MsgHeader = pd.MsgHeader
		pd.pdv.Length = pd.BlockSize + 2
		pd.Length = pd.pdv.Length + 4
		pd.ItemType = 0x04
		pd.Reserved1 = 0
		bd := media.NewEmptyBufData()

		bd.SetBigEndian(true)
		bd.WriteByte(pd.ItemType)
		bd.WriteByte(pd.Reserved1)
		bd.WriteUint32(pd.Length)
		bd.WriteUint32(pd.pdv.Length)
		bd.WriteByte(pd.pdv.PresentationContextID)
		bd.WriteByte(pd.MsgHeader)
		if err := bd.Send(conn); err == nil {
			buff, err := pd.Buffer.Read(int(pd.BlockSize))
			if err != nil {
				return errors.New("ERROR, pdata::Write, " + err.Error())
			}

			n, err := conn.Write(buff)
			if err != nil {
				return errors.New("ERROR, pdata::Write, " + err.Error())
			}

			if n != int(pd.BlockSize) {
				return errors.New("ERROR, pdata::Write, n!=int(pd.BlockSize)")
			}
		} else {
			return errors.New("ERROR, pdata::Write, bd.Send(conn) failed")
		}
		SentSize += pd.BlockSize
	}
	pd.Length = TLength
	return nil
}
