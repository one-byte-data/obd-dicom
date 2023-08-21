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

// AAssociationRQ - AAssociationRQ
type AAssociationRQ interface {
	GetAppContext() UIDItem
	SetAppContext(context UIDItem)
	GetCallingAE() string
	SetCallingAE(AET string)
	GetCalledAE() string
	SetCalledAE(AET string)
	GetPresContexts() []PresentationContext
	GetUserInformation() UserInformation
	SetUserInformation(userInfo UserInformation)
	GetMaxSubLength() uint32
	SetMaxSubLength(length uint32)
	GetImpClass() UIDItem
	SetImpClassUID(uid string)
	SetImpVersionName(name string)
	Size() uint32
	Write(rw *bufio.ReadWriter) error
	Read(ms media.MemoryStream) error
	AddPresContexts(presentationContext PresentationContext)
}

type aassociationRQ struct {
	ItemType        byte // 0x01
	Reserved1       byte
	Length          uint32
	ProtocolVersion uint16 // 0x01
	Reserved2       uint16
	CallingAE       [16]byte // 16 bytes transfered
	CalledAE        [16]byte // 16 bytes transfered
	Reserved3       [32]byte
	AppContext      UIDItem
	PresContexts    []PresentationContext
	UserInfo        UserInformation
}

// NewAAssociationRQ - NewAAssociationRQ
func NewAAssociationRQ() AAssociationRQ {
	return &aassociationRQ{
		ItemType:        0x01,
		Reserved1:       0x00,
		ProtocolVersion: 0x01,
		Reserved2:       0x00,
		AppContext: &uidItem{
			itemType:  0x10,
			reserved1: 0x00,
			uid:       sopclass.DICOMApplicationContext.UID,
			length:    uint16(len(sopclass.DICOMApplicationContext.UID)),
		},
		PresContexts: make([]PresentationContext, 0),
		UserInfo:     NewUserInformation(),
	}
}

func (aarq *aassociationRQ) GetAppContext() UIDItem {
	return aarq.AppContext
}

func (aarq *aassociationRQ) SetAppContext(context UIDItem) {
	aarq.AppContext = context
}

func (aarq *aassociationRQ) GetCallingAE() string {
	temp := []byte{}
	for _, b := range aarq.CallingAE {
		if b != 0x00 && b != 0x20 {
			temp = append(temp, b)
		}
	}
	return string(temp)
}

func (aarq *aassociationRQ) SetCallingAE(AET string) {
	copy(aarq.CallingAE[:], AET)
	for index, b := range aarq.CallingAE {
		if b == 0x00 {
			aarq.CallingAE[index] = 0x20
		}
	}
}

func (aarq *aassociationRQ) GetCalledAE() string {
	temp := []byte{}
	for _, b := range aarq.CalledAE {
		if b != 0x00 && b != 0x20 {
			temp = append(temp, b)
		}
	}
	return string(temp)
}

func (aarq *aassociationRQ) SetCalledAE(AET string) {
	copy(aarq.CalledAE[:], AET)
	for index, b := range aarq.CalledAE {
		if b == 0x00 {
			aarq.CalledAE[index] = 0x20
		}
	}
}

func (aarq *aassociationRQ) GetPresContexts() []PresentationContext {
	return aarq.PresContexts
}

func (aarq *aassociationRQ) GetUserInformation() UserInformation {
	return aarq.UserInfo
}

func (aarq *aassociationRQ) SetUserInformation(userInfo UserInformation) {
	aarq.UserInfo = userInfo
}

func (aarq *aassociationRQ) GetMaxSubLength() uint32 {
	return aarq.UserInfo.GetMaxSubLength().GetMaximumLength()
}

func (aarq *aassociationRQ) SetMaxSubLength(length uint32) {
	aarq.UserInfo.GetMaxSubLength().SetMaximumLength(length)
}

func (aarq *aassociationRQ) GetImpClass() UIDItem {
	return aarq.UserInfo.GetImpClass()
}

func (aarq *aassociationRQ) SetImpClassUID(uid string) {
	aarq.UserInfo.SetImpClassUID(uid)
}

func (aarq *aassociationRQ) SetImpVersionName(name string) {
	aarq.UserInfo.SetImpVersionName(name)
}

func (aarq *aassociationRQ) Size() uint32 {
	aarq.Length = 4 + 16 + 16 + 32
	aarq.Length += uint32(aarq.AppContext.GetSize())

	for _, PresContext := range aarq.PresContexts {
		aarq.Length += uint32(PresContext.Size())
	}

	aarq.Length += uint32(aarq.UserInfo.Size())
	return aarq.Length + 6
}

func (aarq *aassociationRQ) Write(rw *bufio.ReadWriter) error {
	bd := media.NewEmptyBufData()

	fmt.Println()

	log.Printf("INFO, ASSOC-RQ: ImpClass: %s\n", aarq.GetUserInformation().GetImpClass().GetUID())
	log.Printf("INFO, ASSOC-RQ: ImpVersion: %s\n\n", aarq.GetUserInformation().GetImpVersion().GetUID())

	log.Printf("INFO, ASSOC-RQ: CalledAE: %s\n", aarq.CalledAE)
	log.Printf("INFO, ASSOC-RQ: CallingAE: %s\n\n", aarq.CallingAE)

	bd.SetBigEndian(true)
	aarq.Size()
	bd.WriteByte(aarq.ItemType)
	bd.WriteByte(aarq.Reserved1)
	bd.WriteUint32(aarq.Length)
	bd.WriteUint16(aarq.ProtocolVersion)
	bd.WriteUint16(aarq.Reserved2)
	bd.Write(aarq.CalledAE[:], 16)
	bd.Write(aarq.CallingAE[:], 16)
	bd.Write(aarq.Reserved3[:], 32)

	if err := bd.Send(rw); err != nil {
		return err
	}
	log.Printf("INFO, ASSOC-RQ: ApplicationContext: %s - %s\n", aarq.AppContext.GetUID(), sopclass.GetSOPClassFromUID(aarq.AppContext.GetUID()).Description)
	if err := aarq.AppContext.Write(rw); err != nil {
		return err
	}
	for presIndex, presContext := range aarq.PresContexts {
		log.Printf("INFO, ASSOC-RQ: PresentationContext: %d\n", presIndex+1)
		log.Printf("INFO, ASSOC-RQ: \tAbstract Syntax: %s - %s\n", presContext.GetAbstractSyntax().GetUID(), sopclass.GetSOPClassFromUID(presContext.GetAbstractSyntax().GetUID()).Description)
		for _, transSyntax := range presContext.GetTransferSyntaxes() {
			log.Printf("INFO, ASSOC-RQ: \tTransfer Syntax: %s - %s\n", transSyntax.GetUID(), transfersyntax.GetTransferSyntaxFromUID(transSyntax.GetUID()).Description)
		}
		if err := presContext.Write(rw); err != nil {
			return err
		}
	}
	return aarq.UserInfo.Write(rw)
}

func (aarq *aassociationRQ) Read(ms media.MemoryStream) (err error) {
	aarq.ProtocolVersion, err = ms.GetUint16()
	if err != nil {
		return err
	}

	aarq.Reserved2, err = ms.GetUint16()
	if err != nil {
		return err
	}

	ms.ReadData(aarq.CalledAE[:])
	ms.ReadData(aarq.CallingAE[:])
	ms.ReadData(aarq.Reserved3[:])

	Count := int(ms.GetSize() - 4 - 16 - 16 - 32)
	for Count > 0 {
		TempByte, err := ms.GetByte()
		if err != nil {
			return err
		}

		switch TempByte {
		case 0x10:
			aarq.AppContext.SetType(TempByte)
			aarq.AppContext.ReadDynamic(ms)
			Count = Count - int(aarq.AppContext.GetSize())
		case 0x20:
			PresContext := NewPresentationContext()
			PresContext.ReadDynamic(ms)
			Count = Count - int(PresContext.Size())
			aarq.PresContexts = append(aarq.PresContexts, PresContext)
		case 0x50: // User Information
			aarq.UserInfo.SetItemType(TempByte)
			aarq.UserInfo.ReadDynamic(ms)
			return nil
		default:
			log.Println("ERROR, aarq::ReadDynamic, unknown Item, " + strconv.Itoa(int(TempByte)))
			Count = -1
		}
	}

	log.Printf("INFO, ASSOC-RQ: CalledAE - %s\n", aarq.CalledAE)
	log.Printf("INFO, ASSOC-RQ: CallingAE - %s\n", aarq.CallingAE)
	log.Printf("INFO, ASSOC-RQ: \tImpClass %s\n", aarq.GetUserInformation().GetImpClass().GetUID())
	log.Printf("INFO, ASSOC-RQ: \tImpVersion %s\n\n", aarq.GetUserInformation().GetImpVersion().GetUID())

	if Count == 0 {
		return nil
	}

	return errors.New("ERROR, aarq::ReadDynamic, Count is not zero")
}

func (aarq *aassociationRQ) AddPresContexts(presentationContext PresentationContext) {
	aarq.PresContexts = append(aarq.PresContexts, presentationContext)
}
