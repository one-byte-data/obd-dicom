package network

import (
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// AAssociationRJ association reject struct
type AAssociationRJ struct {
	ItemType  byte // 0x03
	Reserved1 byte
	Length    uint32
	Reserved2 byte
	Result    byte
	Source    byte
	Reason    byte
}

// NewAAssociationRJ creates an association reject
func NewAAssociationRJ() *AAssociationRJ {
	return &AAssociationRJ{
		ItemType:  0x03,
		Reserved1: 0x00,
		Reserved2: 0x00,
		Result:    0x01,
		Source:    0x03,
		Reason:    1,
	}
}

// Size gets the size
func (aarj *AAssociationRJ) Size() uint32 {
	aarj.Length = 4
	return aarj.Length + 6
}

func (aarj *AAssociationRJ) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	aarj.Size()
	bd.WriteByte(aarj.ItemType)
	bd.WriteByte(aarj.Reserved1)
	bd.WriteUint32(aarj.Length)
	bd.WriteByte(aarj.Reserved2)
	bd.WriteByte(aarj.Result)
	bd.WriteByte(aarj.Source)
	bd.WriteByte(aarj.Reason)

	flag = bd.Send(conn)
	return flag
}

func (aarj *AAssociationRJ) Read(conn net.Conn) bool {
	aarj.ItemType = ReadByte(conn)
	return aarj.ReadDynamic(conn)
}

// ReadDynamic ReadDynamic
func (aarj *AAssociationRJ) ReadDynamic(conn net.Conn) bool {
	aarj.Reserved1 = ReadByte(conn)
	aarj.Length = ReadUint32(conn)
	aarj.Reserved2 = ReadByte(conn)
	aarj.Result = ReadByte(conn)
	aarj.Source = ReadByte(conn)
	aarj.Reason = ReadByte(conn)
	return true
}

// AReleaseRQ AReleaseRQ
type AReleaseRQ struct {
	ItemType  byte // 0x05
	Reserved1 byte
	Length    uint32
	Reserved2 uint32
}

// NewAReleaseRQ NewAReleaseRQ
func NewAReleaseRQ() *AReleaseRQ {
	return &AReleaseRQ{
		ItemType:  0x05,
		Reserved1: 0x00,
		Reserved2: 0x00,
	}
}

// Size gets the size
func (arrq *AReleaseRQ) Size() uint32 {
	arrq.Length = 4
	return arrq.Length + 6
}

func (arrq *AReleaseRQ) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	arrq.Size()
	bd.WriteByte(arrq.ItemType)
	bd.WriteByte(arrq.Reserved1)
	bd.WriteUint32(arrq.Length)
	bd.WriteUint32(arrq.Reserved2)

	flag = bd.Send(conn)
	return flag
}

func (arrq *AReleaseRQ) Read(conn net.Conn) bool {
	arrq.ItemType = ReadByte(conn)
	return arrq.ReadDynamic(conn)
}

// ReadDynamic ReadDynamic
func (arrq *AReleaseRQ) ReadDynamic(conn net.Conn) bool {
	arrq.Reserved1 = ReadByte(conn)
	arrq.Length = ReadUint32(conn)
	arrq.Reserved2 = ReadUint32(conn)
	return true
}

// AReleaseRP - AReleaseRP
type AReleaseRP struct {
	ItemType  byte // 0x06
	Reserved1 byte
	Length    uint32
	Reserved2 uint32
}

// NewAReleaseRP - NewAReleaseRP
func NewAReleaseRP() *AReleaseRP {
	return &AReleaseRP{
		ItemType:  0x06,
		Reserved1: 0x00,
		Reserved2: 0x00,
	}
}

// Size gets the size
func (arrp *AReleaseRP) Size() uint32 {
	arrp.Length = 4
	return arrp.Length + 6
}

func (arrp *AReleaseRP) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	arrp.Size()
	bd.WriteByte(arrp.ItemType)
	bd.WriteByte(arrp.Reserved1)
	bd.WriteUint32(arrp.Length)
	bd.WriteUint32(arrp.Reserved2)

	flag = bd.Send(conn)
	return flag
}

func (arrp *AReleaseRP) Read(conn net.Conn) bool {
	arrp.ItemType = ReadByte(conn)
	return arrp.ReadDynamic(conn)
}

// ReadDynamic ReadDynamic
func (arrp *AReleaseRP) ReadDynamic(conn net.Conn) bool {
	arrp.Reserved1 = ReadByte(conn)
	arrp.Length = ReadUint32(conn)
	arrp.Reserved2 = ReadUint32(conn)
	return true
}

// AAbortRQ - AAbortRQ
type AAbortRQ struct {
	ItemType  byte // 0x07
	Reserved1 byte
	Length    uint32
	Reserved2 byte
	Reserved3 byte
	Source    byte
	Reason    byte
}

// NewAAbortRQ - NewAAbortRQ
func NewAAbortRQ() *AAbortRQ {
	return &AAbortRQ{
		ItemType:  0x07,
		Reserved1: 0x00,
		Reserved2: 0x00,
		Reserved3: 0x01,
		Source:    0x03,
		Reason:    0x01,
	}
}

// Size gets the size
func (aarq *AAbortRQ) Size() uint32 {
	aarq.Length = 4
	return aarq.Length + 6
}

func (aarq *AAbortRQ) Write(conn net.Conn) bool {
	flag := false
	var bd media.BufData

	bd.BigEndian = true
	aarq.Size()
	bd.WriteByte(aarq.ItemType)
	bd.WriteByte(aarq.Reserved1)
	bd.WriteUint32(aarq.Length)
	bd.WriteByte(aarq.Reserved2)
	bd.WriteByte(aarq.Reserved3)
	bd.WriteByte(aarq.Source)
	bd.WriteByte(aarq.Reason)

	flag = bd.Send(conn)
	return flag
}

func (aarq *AAbortRQ) Read(conn net.Conn) bool {
	aarq.ItemType = ReadByte(conn)
	return aarq.ReadDynamic(conn)
}

// ReadDynamic - ReadDynamic
func (aarq *AAbortRQ) ReadDynamic(conn net.Conn) bool {
	aarq.Reserved1 = ReadByte(conn)
	aarq.Length = ReadUint32(conn)
	aarq.Reserved2 = ReadByte(conn)
	aarq.Reserved3 = ReadByte(conn)
	aarq.Source = ReadByte(conn)
	aarq.Reason = ReadByte(conn)
	return true
}
