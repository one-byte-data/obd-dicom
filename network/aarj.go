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

	log.Printf("INFO, ASSOC-R-RQ: <-- %x\n", arrq.Reserved1)

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
