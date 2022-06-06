package network

import (
	"bufio"

	"github.com/one-byte-data/obd-dicom/media"
)

// AReleaseRQ AReleaseRQ
type AReleaseRQ interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type areleaseRQ struct {
	ItemType  byte // 0x05
	Reserved1 byte
	Length    uint32
	Reserved2 uint32
}

// NewAReleaseRQ NewAReleaseRQ
func NewAReleaseRQ() AReleaseRQ {
	return &areleaseRQ{
		ItemType:  0x05,
		Reserved1: 0x00,
		Reserved2: 0x00,
	}
}

func (arrq *areleaseRQ) Size() uint32 {
	arrq.Length = 4
	return arrq.Length + 6
}

func (arrq *areleaseRQ) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	arrq.Size()
	bd.WriteByte(arrq.ItemType)
	bd.WriteByte(arrq.Reserved1)
	bd.WriteUint32(arrq.Length)
	bd.WriteUint32(arrq.Reserved2)

	return bd.Send(rw)
}

func (arrq *areleaseRQ) Read(ms media.MemoryStream) (err error) {
	arrq.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return arrq.ReadDynamic(ms)
}

func (arrq *areleaseRQ) ReadDynamic(ms media.MemoryStream) (err error) {
	arrq.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	arrq.Length, err = ms.GetUint32()
	if err != nil {
		return err
	}
	arrq.Reserved2, err = ms.GetUint32()
	if err != nil {
		return err
	}
	return
}
