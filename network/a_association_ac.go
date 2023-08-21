package network

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/one-byte-data/obd-dicom/dictionary/sopclass"
	"github.com/one-byte-data/obd-dicom/dictionary/transfersyntax"
	"github.com/one-byte-data/obd-dicom/media"
)

// AAssociationAC AAssociationAC
type AAssociationAC interface {
	GetAppContext() UIDItem
	SetAppContext(context UIDItem)
	SetCallingAE(AET string)
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
	AppContext         UIDItem
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
		AppContext: &uidItem{
			itemType:  0x10,
			reserved1: 0x00,
			uid:       sopclass.DICOMApplicationContext.UID,
			length:    uint16(len(sopclass.DICOMApplicationContext.UID)),
		},
		PresContextAccepts: make([]PresentationContextAccept, 0),
		UserInfo:           NewUserInformation(),
	}
}

func (aaac *aassociationAC) GetAppContext() UIDItem {
	return aaac.AppContext
}

func (aaac *aassociationAC) SetAppContext(context UIDItem) {
	aaac.AppContext = context
}

func (aaac *aassociationAC) SetCallingAE(AET string) {
	copy(aaac.CallingAE[:], AET)
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
	aaac.Length += uint32(aaac.AppContext.GetSize())

	for _, PresContextAccept := range aaac.PresContextAccepts {
		aaac.Length += uint32(PresContextAccept.Size())
	}

	aaac.Length += uint32(aaac.UserInfo.Size())
	return aaac.Length + 6
}

func (aaac *aassociationAC) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	fmt.Println()

	log.Printf("INFO, ASSOC-AC: ImpClass: %s\n", aaac.UserInfo.GetImpClass().GetUID())
	log.Printf("INFO, ASSOC-AC: ImpVersion: %s\n\n", aaac.UserInfo.GetImpVersion().GetUID())

	log.Printf("INFO, ASSOC-AC: CalledAE: %s\n", aaac.CalledAE)
	log.Printf("INFO, ASSOC-AC: CallingAE: %s\n\n", aaac.CallingAE)

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

	if err := bd.Send(rw); err != nil {
		return err
	}
	log.Printf("INFO, ASSOC-AC: ApplicationContext: %s - %s\n", aaac.AppContext.GetUID(), sopclass.GetSOPClassFromUID(aaac.AppContext.GetUID()).Description)
	if err := aaac.AppContext.Write(rw); err != nil {
		return err
	}
	for presIndex, presContextAccept := range aaac.PresContextAccepts {
		log.Printf("INFO, ASSOC-AC: PresentationContext: %d\n", presIndex + 1)
		log.Printf("INFO, ASSOC-AC: \tAbstract Syntax: %s - %s\n", presContextAccept.GetAbstractSyntax().GetUID(), sopclass.GetSOPClassFromUID(presContextAccept.GetAbstractSyntax().GetUID()).Description)
		log.Printf("INFO, ASSOC-AC: \tTransfer Syntax: %s - %s\n", presContextAccept.GetTrnSyntax().GetUID(), transfersyntax.GetTransferSyntaxFromUID(presContextAccept.GetTrnSyntax().GetUID()).Description)
		if err := presContextAccept.Write(rw); err != nil {
			return err
		}
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
			Count = Count - int(aaac.AppContext.GetSize())
		case 0x21:
			PresContextAccept := NewPresentationContextAccept()
			PresContextAccept.ReadDynamic(ms)
			Count = Count - int(PresContextAccept.Size())
			aaac.PresContextAccepts = append(aaac.PresContextAccepts, PresContextAccept)
		case 0x50: // User Information
			aaac.UserInfo.ReadDynamic(ms)
			Count = Count - int(aaac.UserInfo.Size())
		default:
			Count = -1
			return errors.New("ERROR, aaac::ReadDynamic, unknown Item, " + strconv.Itoa(int(TempByte)))
		}
	}

	log.Printf("INFO, ASSOC-AC: CalledAE - %s\n", aaac.CalledAE)
	log.Printf("INFO, ASSOC-AC: CallingAE - %s\n", aaac.CallingAE)
	log.Printf("INFO, ASSOC-AC: \tImpClass %s\n", aaac.GetUserInformation().GetImpClass().GetUID())
	log.Printf("INFO, ASSOC-AC: \tImpVersion %s\n\n", aaac.GetUserInformation().GetImpVersion().GetUID())

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, aarq::ReadDynamic, Count is not zero")
}
