package network

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// AAssociationAC AAssociationAC
type AAssociationAC interface {
	GetAppContext() UIDitem
	SetAppContext(context UIDitem)
	GetCallingAE() string
	SetCallingAE(AET string)
	GetCalledAE() string
	SetCalledAE(AET string)
	AddPresContextAccept(context PresentationContextAccept)
	GetPresContextAccepts() []PresentationContextAccept
	GetUserInformation() UserInformation
	SetUserInformation(UserInfo UserInformation)
	GetMaxSubLength() uint32
	SetMaxSubLength(length uint32)
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type aassociationAC struct {
	ItemType           byte
	Reserved1          byte
	Length             uint32
	ProtocolVersion    uint16
	Reserved2          uint16
	CallingAE          [16]byte
	CalledAE           [16]byte
	Reserved3          [32]byte
	AppContext         UIDitem
	PresContextAccepts []PresentationContextAccept
	UserInfo           UserInformation
}

// NewAAssociationAC NewAAssociationAC
func NewAAssociationAC() AAssociationAC {
	return &aassociationAC{
		ItemType:        0x02,
		Reserved1:       0x00,
		ProtocolVersion: 0x01,
		Reserved2:       0x00,
		AppContext: UIDitem{
			ItemType:  0x10,
			Reserved1: 0x00,
			UIDName:   "1.2.840.10008.3.1.1.1",
			Length:    uint16(len("1.2.840.10008.3.1.1.1")),
		},
		PresContextAccepts: make([]PresentationContextAccept, 0),
		UserInfo:           NewUserInformation(),
	}
}

func (aaac *aassociationAC) GetAppContext() UIDitem {
	return aaac.AppContext
}

func (aaac *aassociationAC) SetAppContext(context UIDitem) {
	aaac.AppContext = context
}

func (aaac *aassociationAC) GetCallingAE() string {
	temp := strings.ReplaceAll(fmt.Sprintf("%s", aaac.CallingAE), "\x20", "\x00")
	return strings.ReplaceAll(temp, "\x00", "")
}

func (aaac *aassociationAC) SetCallingAE(AET string) {
	copy(aaac.CallingAE[:], AET)
}

func (aaac *aassociationAC) GetCalledAE() string {
	temp := strings.ReplaceAll(fmt.Sprintf("%s", aaac.CalledAE), "\x20", "\x00")
	return strings.ReplaceAll(temp, "\x00", "")
}

func (aaac *aassociationAC) SetCalledAE(AET string) {
	copy(aaac.CalledAE[:], AET)
}

func (aaac *aassociationAC) AddPresContextAccept(context PresentationContextAccept) {
	aaac.PresContextAccepts = append(aaac.PresContextAccepts, context)
}

func (aaac *aassociationAC) GetPresContextAccepts() []PresentationContextAccept {
	return aaac.PresContextAccepts
}

func (aaac *aassociationAC) GetUserInformation() UserInformation {
	return aaac.UserInfo
}

func (aaac *aassociationAC) SetUserInformation(UserInfo UserInformation) {
	aaac.UserInfo = UserInfo
}

func (aaac *aassociationAC) GetMaxSubLength() uint32 {
	return aaac.UserInfo.GetMaxSubLength().GetMaximumLength()
}

func (aaac *aassociationAC) SetMaxSubLength(length uint32) {
	aaac.UserInfo.GetMaxSubLength().SetMaximumLength(length)
}

func (aaac *aassociationAC) Size() uint32 {
	aaac.Length = 4 + 16 + 16 + 32
	aaac.Length += uint32(aaac.AppContext.Size())

	for _, PresContextAccept := range aaac.PresContextAccepts {
		aaac.Length += uint32(PresContextAccept.Size())
	}

	aaac.Length += uint32(aaac.UserInfo.Size())
	return aaac.Length + 6
}

func (aaac *aassociationAC) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	fmt.Println()

	log.Printf("INFO, ASSOC-AC: %s <-- %s\n", aaac.CallingAE, aaac.CalledAE)
	log.Printf("INFO, ASSOC-AC: \tImpClass %s\n", aaac.UserInfo.GetImpClass().UIDName)
	log.Printf("INFO, ASSOC-AC: \tImpVersion %s\n\n", aaac.UserInfo.GetImpVersion().UIDName)

	bd.SetBigEndian(true)
	aaac.Size()
	bd.WriteByte(aaac.ItemType)
	bd.WriteByte(aaac.Reserved1)
	bd.WriteUint32(aaac.Length)
	bd.WriteUint16(aaac.ProtocolVersion)
	bd.WriteUint16(aaac.Reserved2)
	bd.Write(aaac.CalledAE[:], 16)
	bd.Write(aaac.CallingAE[:], 16)
	bd.Write(aaac.Reserved3[:], 32)

	err := bd.Send(rw)
	if err != nil {
		return err
	}
	err = aaac.AppContext.Write(rw)
	if err != nil {
		return err
	}
	for _, PresContextAccept := range aaac.PresContextAccepts {
		PresContextAccept.Write(rw)
	}
	return aaac.UserInfo.Write(rw)
}

func (aaac *aassociationAC) Read(ms media.MemoryStream) (err error) {
	aaac.ItemType, err = ms.GetByte()
	if err != nil {
		return err
	}
	return aaac.ReadDynamic(ms)
}

func (aaac *aassociationAC) ReadDynamic(ms media.MemoryStream) (err error) {
	aaac.Reserved1, err = ms.GetByte()
	if err != nil {
		return err
	}
	aaac.Length, err = ms.GetUint32()
	if err != nil {
		return err
	}
	aaac.ProtocolVersion, err = ms.GetUint16()
	if err != nil {
		return err
	}
	aaac.Reserved2, err = ms.GetUint16()
	if err != nil {
		return err
	}

	ms.ReadData(aaac.CalledAE[:])
	ms.ReadData(aaac.CallingAE[:])
	ms.ReadData(aaac.Reserved3[:])

	Count := int(aaac.Length - 4 - 16 - 16 - 32)
	for Count > 0 {
		TempByte, err := ms.GetByte()
		if err != nil {
			return err
		}

		switch TempByte {
		case 0x10:
			aaac.AppContext.ReadDynamic(ms)
			Count = Count - int(aaac.AppContext.Size())
			break
		case 0x21:
			PresContextAccept := NewPresentationContextAccept()
			PresContextAccept.ReadDynamic(ms)
			Count = Count - int(PresContextAccept.Size())
			aaac.PresContextAccepts = append(aaac.PresContextAccepts, PresContextAccept)
			break
		case 0x50: // User Information
			aaac.UserInfo.ReadDynamic(ms)
			Count = Count - int(aaac.UserInfo.Size())
			break
		default:
			Count = -1
			return errors.New("ERROR, aaac::ReadDynamic, unknown Item, " + strconv.Itoa(int(TempByte)))
		}
	}

	log.Printf("INFO, ASSOC-AC: %s --> %s\n", aaac.GetCallingAE(), aaac.GetCalledAE())
	log.Printf("INFO, ASSOC-AC: \tImpClass %s\n", aaac.GetUserInformation().GetImpClass().UIDName)
	log.Printf("INFO, ASSOC-AC: \tImpVersion %s\n\n", aaac.GetUserInformation().GetImpVersion().UIDName)

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, aarq::ReadDynamic, Count is not zero")
}
