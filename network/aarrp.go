package network

import (
	"bufio"
	"log"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// AReleaseRP - AReleaseRP
type AReleaseRP interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type areleaseRP struct {
	ItemType  byte // 0x06
	Reserved1 byte
	Length    uint32
	Reserved2 uint32
}

// NewAReleaseRP - NewAReleaseRP
func NewAReleaseRP() AReleaseRP {
	return &areleaseRP{
		ItemType:  0x06,
		Reserved1: 0x00,
		Reserved2: 0x00,
	}
}

func (arrp *areleaseRP) Size() uint32 {
	arrp.Length = 4
	return arrp.Length + 6
}

func (arrp *areleaseRP) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	log.Printf("INFO, ASSOC-R-RP: %x -->\n", arrp.Reserved1)

	bd.SetBigEndian(true)
	arrp.Size()
	bd.WriteByte(arrp.ItemType)
	bd.WriteByte(arrp.Reserved1)
	bd.WriteUint32(arrp.Length)
	bd.WriteUint32(arrp.Reserved2)

	return bd.Send(rw)
}

func (arrp *areleaseRP) Read(ms media.MemoryStream) (err error) {
	arrp.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return arrp.ReadDynamic(ms)
}

func (arrp *areleaseRP) ReadDynamic(ms media.MemoryStream) (err error) {
	arrp.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	arrp.Length, err = ms.GetUint32()
	if err != nil {
		return err
	}
	arrp.Reserved2, err = ms.GetUint32()
	return
}