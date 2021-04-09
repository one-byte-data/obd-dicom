package network

import (
	"bufio"

	"git.onebytedata.com/odb/go-dicom/media"
)

// AAbortRQ - AAbortRQ
type AAbortRQ interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type aabortRQ struct {
	ItemType  byte // 0x07
	Reserved1 byte
	Length    uint32
	Reserved2 byte
	Reserved3 byte
	Source    byte
	Reason    byte
}

// NewAAbortRQ - NewAAbortRQ
func NewAAbortRQ() AAbortRQ {
	return &aabortRQ{
		ItemType:  0x07,
		Reserved1: 0x00,
		Reserved2: 0x00,
		Reserved3: 0x01,
		Source:    0x03,
		Reason:    0x01,
	}
}

func (aarq *aabortRQ) Size() uint32 {
	aarq.Length = 4
	return aarq.Length + 6
}

func (aarq *aabortRQ) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	aarq.Size()
	bd.WriteByte(aarq.ItemType)
	bd.WriteByte(aarq.Reserved1)
	bd.WriteUint32(aarq.Length)
	bd.WriteByte(aarq.Reserved2)
	bd.WriteByte(aarq.Reserved3)
	bd.WriteByte(aarq.Source)
	bd.WriteByte(aarq.Reason)

	return bd.Send(rw)
}

func (aarq *aabortRQ) Read(ms media.MemoryStream) (err error) {
	aarq.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return aarq.ReadDynamic(ms)
}

func (aarq *aabortRQ) ReadDynamic(ms media.MemoryStream) (err error) {
	aarq.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarq.Length, err = ms.GetUint32()
	if err != nil {
		return err
	}
	aarq.Reserved2, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarq.Reserved3, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarq.Source, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarq.Reason, err = ms.GetByte()
	if err != nil {
		return err
	}
	return
}
