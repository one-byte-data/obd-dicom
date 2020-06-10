package network

import(
	 "net"
	 "git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

type AAssociationRJ struct {
	ItemType        byte // 0x03
	Reserved1       byte
	Length          uint32
	Reserved2 byte
	Result byte
	Source byte
	Reason byte
}

func NewAAssociationRJ() *AAssociationRJ {
	aarj := &AAssociationRJ{}
	aarj.ItemType = 0x03
	aarj.Reserved1 = 0x00
	aarj.Reserved2 = 0x00
	aarj.Result = 0x01
	aarj.Source = 0x03
	aarj.Reason =1
	return aarj
}

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

	flag=bd.Send(conn)
	return flag
}

func (aarj *AAssociationRJ) Read(conn net.Conn) bool {
	aarj.ItemType=ReadByte(conn)
	return aarj.ReadDynamic(conn)
}

func (aarj *AAssociationRJ) ReadDynamic(conn net.Conn) bool {
	aarj.Reserved1=ReadByte(conn)
	aarj.Length=ReadUint32(conn)
	aarj.Reserved2=ReadByte(conn)
	aarj.Result=ReadByte(conn)
	aarj.Source=ReadByte(conn)
	aarj.Reason=ReadByte(conn)
	return true
}

type AReleaseRQ struct {
	ItemType        byte // 0x05
	Reserved1       byte
	Length          uint32
	Reserved2 uint32
}

func NewAReleaseRQ() *AReleaseRQ {
	arrq := &AReleaseRQ{}
	arrq.ItemType = 0x05
	arrq.Reserved1 = 0x00
	arrq.Reserved2 = 0x00
	return arrq
}

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

	flag=bd.Send(conn)
	return flag
}

func (arrq *AReleaseRQ) Read(conn net.Conn) bool {
	arrq.ItemType=ReadByte(conn)
	return arrq.ReadDynamic(conn)
}

func (arrq *AReleaseRQ) ReadDynamic(conn net.Conn) bool {
	arrq.Reserved1=ReadByte(conn)
	arrq.Length=ReadUint32(conn)
	arrq.Reserved2=ReadUint32(conn)
	return true
}

type AReleaseRP struct {
	ItemType        byte // 0x06
	Reserved1       byte
	Length          uint32
	Reserved2 uint32
}

func NewAReleaseRP() *AReleaseRP {
	arrp := &AReleaseRP{}
	arrp.ItemType = 0x06
	arrp.Reserved1 = 0x00
	arrp.Reserved2 = 0x00
	return arrp
}

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

	flag=bd.Send(conn)
	return flag
}

func (arrp *AReleaseRP) Read(conn net.Conn) bool {
	arrp.ItemType=ReadByte(conn)
	return arrp.ReadDynamic(conn)
}

func (arrp *AReleaseRP) ReadDynamic(conn net.Conn) bool {
	arrp.Reserved1=ReadByte(conn)
	arrp.Length=ReadUint32(conn)
	arrp.Reserved2=ReadUint32(conn)
	return true
}

type AAbortRQ struct {
	ItemType        byte // 0x07
	Reserved1       byte
	Length          uint32
	Reserved2 byte
	Reserved3 byte
	Source byte
	Reason byte
}

func NewAAbortRQ() *AAbortRQ {
	aarq := &AAbortRQ{}
	aarq.ItemType = 0x07
	aarq.Reserved1 = 0x00
	aarq.Reserved2 = 0x00
	aarq.Reserved3 = 0x01
	aarq.Source = 0x03
	aarq.Reason =0x01
	return aarq
}

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

	flag=bd.Send(conn)
	return flag
}

func (aarq *AAbortRQ) Read(conn net.Conn) bool {
	aarq.ItemType=ReadByte(conn)
	return aarq.ReadDynamic(conn)
}

func (aarq *AAbortRQ) ReadDynamic(conn net.Conn) bool {
	aarq.Reserved1=ReadByte(conn)
	aarq.Length=ReadUint32(conn)
	aarq.Reserved2=ReadByte(conn)
	aarq.Reserved3=ReadByte(conn)
	aarq.Source=ReadByte(conn)
	aarq.Reason=ReadByte(conn)
	return true
}
