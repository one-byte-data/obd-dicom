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
	Read(rw *bufio.ReadWriter) (err error)
	ReadDynamic(rw *bufio.ReadWriter) (err error)
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

func (aarj *aassociationRJ) Read(rw *bufio.ReadWriter) (err error) {
	aarj.ItemType, err = ReadByte(rw)
	return aarj.ReadDynamic(rw)
}

func (aarj *aassociationRJ) ReadDynamic(rw *bufio.ReadWriter) (err error) {
	aarj.Reserved1, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarj.Length, err = ReadUint32(rw)
	if err != nil {
		return
	}
	aarj.Reserved2, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarj.Result, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarj.Source, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarj.Reason, err = ReadByte(rw)
	if err != nil {
		return
	}
	return
}

// AReleaseRQ AReleaseRQ
type AReleaseRQ interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(rw *bufio.ReadWriter) (err error)
	ReadDynamic(rw *bufio.ReadWriter) (err error)
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

func (arrq *areleaseRQ) Read(rw *bufio.ReadWriter) (err error) {
	arrq.ItemType, err = ReadByte(rw)
	if err != nil {
		return
	}
	return arrq.ReadDynamic(rw)
}

func (arrq *areleaseRQ) ReadDynamic(rw *bufio.ReadWriter) (err error) {
	arrq.Reserved1, err = ReadByte(rw)
	if err != nil {
		return
	}
	arrq.Length, err = ReadUint32(rw)
	if err != nil {
		return
	}
	arrq.Reserved2, err = ReadUint32(rw)
	return
}

// AReleaseRP - AReleaseRP
type AReleaseRP interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(rw *bufio.ReadWriter) (err error)
	ReadDynamic(rw *bufio.ReadWriter) (err error)
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

func (arrp *areleaseRP) Read(rw *bufio.ReadWriter) (err error) {
	arrp.ItemType, err = ReadByte(rw)
	if err != nil {
		return
	}
	return arrp.ReadDynamic(rw)
}

func (arrp *areleaseRP) ReadDynamic(rw *bufio.ReadWriter) (err error) {
	arrp.Reserved1, err = ReadByte(rw)
	if err != nil {
		return
	}
	arrp.Length, err = ReadUint32(rw)
	if err != nil {
		return
	}
	arrp.Reserved2, err = ReadUint32(rw)
	return
}

// AAbortRQ - AAbortRQ
type AAbortRQ interface {
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(rw *bufio.ReadWriter) (err error)
	ReadDynamic(rw *bufio.ReadWriter) (err error)
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

func (aarq *aabortRQ) Read(rw *bufio.ReadWriter) (err error) {
	aarq.ItemType, err = ReadByte(rw)
	if err != nil {
		return
	}
	return aarq.ReadDynamic(rw)
}

func (aarq *aabortRQ) ReadDynamic(rw *bufio.ReadWriter) (err error) {
	aarq.Reserved1, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarq.Length, err = ReadUint32(rw)
	if err != nil {
		return
	}
	aarq.Reserved2, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarq.Reserved3, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarq.Source, err = ReadByte(rw)
	if err != nil {
		return
	}
	aarq.Reason, err = ReadByte(rw)
	return
}
