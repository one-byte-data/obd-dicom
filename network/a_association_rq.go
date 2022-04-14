package network

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"git.onebytedata.com/odb/go-dicom/dictionary/sopclass"
	"git.onebytedata.com/odb/go-dicom/media"
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
	temp := strings.ReplaceAll(fmt.Sprintf("%s", aarq.CallingAE), "\x20", "\x00")
	return strings.ReplaceAll(temp, "\x00", "")
}

func (aarq *aassociationRQ) SetCallingAE(AET string) {
	copy(aarq.CallingAE[:], AET)
}

func (aarq *aassociationRQ) GetCalledAE() string {
	temp := strings.ReplaceAll(fmt.Sprintf("%s", aarq.CalledAE), "\x20", "\x00")
	return strings.ReplaceAll(temp, "\x00", "")
}

func (aarq *aassociationRQ) SetCalledAE(AET string) {
	copy(aarq.CalledAE[:], AET)
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

	log.Printf("INFO, ASSOC-RQ: CalledAE - %s\n", aarq.CalledAE)
	log.Printf("INFO, ASSOC-RQ: CallingAE - %s\n", aarq.CallingAE)
	log.Printf("INFO, ASSOC-RQ: \tImpClass %s\n", aarq.GetUserInformation().GetImpClass().GetUID())
	log.Printf("INFO, ASSOC-RQ: \tImpVersion %s\n\n", aarq.GetUserInformation().GetImpVersion().GetUID())

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

	if err := bd.Send(rw); err == nil {
		err = aarq.AppContext.Write(rw)
		if err != nil {
			return err
		}
		for _, PresContext := range aarq.PresContexts {
			PresContext.Write(rw)
		}
		aarq.UserInfo.Write(rw)
	}
	return nil
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
			break
		case 0x20:
			PresContext := NewPresentationContext()
			PresContext.ReadDynamic(ms)
			Count = Count - int(PresContext.Size())
			aarq.PresContexts = append(aarq.PresContexts, PresContext)
			break
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
