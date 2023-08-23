package network

import (
	"bufio"

	"github.com/one-byte-data/obd-dicom/media"
)

// MaximumSubLength - MaximumSubLength
type MaximumSubLength interface {
	GetMaximumLength() uint32
	SetMaximumLength(length uint32)
	Size() uint16
	Write(rw *bufio.ReadWriter) bool
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type maximumSubLength struct {
	ItemType      byte //0x51
	Reserved1     byte
	Length        uint16
	MaximumLength uint32
}

// NewMaximumSubLength - NewMaximumSubLength
func NewMaximumSubLength() MaximumSubLength {
	return &maximumSubLength{
		ItemType: 0x51,
		Length:   4,
	}
}

func (maxim *maximumSubLength) GetMaximumLength() uint32 {
	return maxim.MaximumLength
}

func (maxim *maximumSubLength) SetMaximumLength(length uint32) {
	maxim.MaximumLength = length
}

func (maxim *maximumSubLength) Size() uint16 {
	return maxim.Length + 4
}

func (maxim *maximumSubLength) Write(rw *bufio.ReadWriter) bool {
	bd := media.NewEmptyBufData()

	bd.SetBigEndian(true)
	bd.WriteByte(maxim.ItemType)
	bd.WriteByte(maxim.Reserved1)
	bd.WriteUint16(maxim.Length)
	bd.WriteUint32(maxim.MaximumLength)

	if err := bd.Send(rw); err != nil {
		return false
	}
	return true
}

func (maxim *maximumSubLength) Read(ms media.MemoryStream) (err error) {
	if maxim.ItemType, err = ms.GetByte(); err != nil {
		return err
	}
	return maxim.ReadDynamic(ms)
}

func (maxim *maximumSubLength) ReadDynamic(ms media.MemoryStream) (err error) {
	if maxim.Reserved1, err = ms.GetByte(); err != nil {
		return err
	}
	if maxim.Length, err = ms.GetUint16(); err != nil {
		return err
	}
	if maxim.MaximumLength, err = ms.GetUint32(); err != nil {
		return err
	}
	return
}
