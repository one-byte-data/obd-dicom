package network

import (
	"bufio"
	"log"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// AAssociationRJ association reject struct
type AAssociationRJ interface {
	Set(result byte, reason byte)
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type aassociationRJ struct {
	ItemType  byte // 0x03
	Reserved1 byte
	Length    uint32
	Reserved2 byte
	Result    byte
	Source    byte
	Reason    byte
}

// NewAAssociationRJ creates an association reject
func NewAAssociationRJ() AAssociationRJ {
	return &aassociationRJ{
		ItemType:  0x03,
		Reserved1: 0x00,
		Reserved2: 0x00,
		Result:    0x01,
		Source:    0x03,
		Reason:    1,
	}
}

func (aarj *aassociationRJ) Size() uint32 {
	aarj.Length = 4
	return aarj.Length + 6
}

func (aarj *aassociationRJ) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	log.Printf("INFO, ASSOC-RJ: Reason %x\n", aarj.Reason)

	bd.SetBigEndian(true)
	aarj.Size()
	bd.WriteByte(aarj.ItemType)
	bd.WriteByte(aarj.Reserved1)
	bd.WriteUint32(aarj.Length)
	bd.WriteByte(aarj.Reserved2)
	bd.WriteByte(aarj.Result)
	bd.WriteByte(aarj.Source)
	bd.WriteByte(aarj.Reason)

	return bd.Send(rw)
}

func (aarj *aassociationRJ) Set(result byte, reason byte) {
	aarj.Result = result
	aarj.Reason = reason
}

func (aarj *aassociationRJ) Read(ms media.MemoryStream) (err error) {
	aarj.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return aarj.ReadDynamic(ms)
}

func (aarj *aassociationRJ) ReadDynamic(ms media.MemoryStream) (err error) {
	aarj.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarj.Length, err = ms.GetUint32()
	if err != nil {
		return err
	}
	aarj.Reserved2, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarj.Result, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarj.Source, err = ms.GetByte()
	if err != nil {
		return err
	}
	aarj.Reason, err = ms.GetByte()
	if err != nil {
		return err
	}
	return
}
