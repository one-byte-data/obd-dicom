package network

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"git.onebytedata.com/odb/go-dicom/imp"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network/pdutype"
	"git.onebytedata.com/odb/go-dicom/uid"
)

// PDUService - struct for PDUService
type PDUService interface {
	GetTransferSyntaxUID(pcid byte) string
	SetTimeout(timeout int)
	Connect(IP string, Port string) error
	Close()
	GetAAssociationRQ() AAssociationRQ
	GetCalledAE() string
	GetCallingAE() string
	SetCalledAE(calledAE string)
	SetCallingAE(callingAE string)
	SetConn(rw *bufio.ReadWriter)
	NextPDU() (media.DcmObj, error)
	AddPresContexts(presentationContext PresentationContext)
	GetPresentationContextID() byte
	SetOnAssociationRequest(f func(request AAssociationRQ) bool)
	Write(DCO media.DcmObj, SOPClass string, ItemType byte) error
	interogateAAssociateAC() bool
	interogateAAssociateRQ(rw *bufio.ReadWriter) error
	parseDCMIntoRaw(DCO media.DcmObj) bool
	parseRawVRIntoDCM(DCO media.DcmObj) bool
	readPDU()
}

type pduService struct {
	AcceptedPresentationContexts []PresentationContextAccept
	readWriter                   *bufio.ReadWriter
	ms                           media.MemoryStream
	pdutype                      int
	pdulength                    uint32
	AssocRQ                      AAssociationRQ
	AssocAC                      AAssociationAC
	AssocRJ                      AAssociationRJ
	ReleaseRQ                    AReleaseRQ
	ReleaseRP                    AReleaseRP
	AbortRQ                      AAbortRQ
	Pdata                        PDataTF
	Timeout                      int
	OnAssociationRequest         func(request AAssociationRQ) bool
}

// NewPDUService - creates a pointer to PDUService
func NewPDUService() PDUService {
	return &pduService{
		ms:        media.NewEmptyMemoryStream(),
		AssocRQ:   NewAAssociationRQ(),
		AssocAC:   NewAAssociationAC(),
		AssocRJ:   NewAAssociationRJ(),
		ReleaseRQ: NewAReleaseRQ(),
		ReleaseRP: NewAReleaseRP(),
		AbortRQ:   NewAAbortRQ(),
	}
}

var maxPduLength uint32 = 16384

func (pdu *pduService) SetConn(rw *bufio.ReadWriter) {
	pdu.readWriter = rw
}

func (pdu *pduService) GetTransferSyntaxUID(pcid byte) string {
	for _, pca := range pdu.AcceptedPresentationContexts {
		if pca.GetPresentationContextID() == pcid {
			return pca.GetTrnSyntax().UIDName
		}
	}
	return ""
}

func (pdu *pduService) SetTimeout(timeout int) {
	pdu.Timeout = timeout
}

func (pdu *pduService) Connect(IP string, Port string) error {
	conn, err := net.Dial("tcp", IP+":"+Port)
	if err != nil {
		return errors.New("ERROR, pduservice::Connect - " + err.Error())
	}

	if pdu.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Duration(int32(pdu.Timeout)) * time.Second))
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	pdu.readWriter = rw
	pdu.AssocRQ.SetMaxSubLength(maxPduLength)
	pdu.AssocRQ.SetImpClassUID(imp.GetImpClassUID())
	pdu.AssocRQ.SetImpVersionName(imp.GetImpVersion())

	err = pdu.AssocRQ.Write(pdu.readWriter)
	if err != nil {
		return err
	}

	pdu.ms = media.NewEmptyMemoryStream()

	pdu.ms.ReadFully(rw, 10)

	ItemType, err := pdu.ms.GetByte()
	if err != nil {
		return err
	}

	_, err = pdu.ms.GetByte()
	if err != nil {
		return err
	}

	pdu.pdulength, err = pdu.ms.GetUint32()
	if err != nil {
		return err
	}

	switch ItemType {
	case pdutype.AssociationAccept:
		pdu.readPDU()
		pdu.ms.SetPosition(1)
		pdu.AssocAC.ReadDynamic(pdu.ms)
		if !pdu.interogateAAssociateAC() {
			return errors.New("ERROR, pduservice::Connect - InterogateAAssociateAC failed")
		}
		return nil
	case pdutype.AssociationReject:
		pdu.readPDU()
		pdu.ms.SetPosition(1)
		pdu.AssocRJ.ReadDynamic(pdu.ms)
		return fmt.Errorf("ERROR, pduservice::Connect - %s", pdu.AssocRJ.GetReason())
	default:
		return errors.New("ERROR, pduservice::Connect - Corrupt Transmision")
	}
}

func (pdu *pduService) Close() {
	pdu.ReleaseRQ.Write(pdu.readWriter)
	pdu.ReleaseRP.Read(pdu.ms)
}

func (pdu *pduService) NextPDU() (command media.DcmObj, err error) {
	if pdu.Pdata.Buffer != nil {
		pdu.Pdata.Buffer.ClearMemoryStream()
	} else {
		pdu.Pdata.Buffer = media.NewEmptyBufData()
	}

	pdu.Pdata.MsgStatus = 0
	if pdu.Pdata.Length != 0 {
		DCO := media.NewEmptyDCMObj()
		pdu.Pdata.ReadDynamic(pdu.ms)
		if pdu.Pdata.MsgStatus > 0 {
			if !pdu.parseRawVRIntoDCM(DCO) {
				pdu.AbortRQ.Write(pdu.readWriter)
				return nil, errors.New("ERROR, pduservice::Read, ParseRawVRIntoDCM failed")
			}
			return DCO, nil
		}
	}

	for {
		pdu.ms = media.NewEmptyMemoryStream()

		pdu.ms.ReadFully(pdu.readWriter, 10)
		pdu.ms.SetPosition(0)

		pdu.pdutype, err = pdu.ms.Get()
		if err != nil {
			return nil, err
		}

		_, err = pdu.ms.Get()
		if err != nil {
			return nil, err
		}

		pdu.pdulength, err = pdu.ms.GetUint32()
		if err != nil {
			return nil, err
		}

		switch pdu.pdutype {
		case pdutype.AssocicationRequest:
			pdu.readPDU()
			err := pdu.AssocRQ.Read(pdu.ms)
			if err != nil {
				return nil, err
			}
			err = pdu.interogateAAssociateRQ(pdu.readWriter)
			if err != nil {
				return nil, err
			}
			return nil, nil
		case pdutype.AssociationAccept:
			pdu.readPDU()
			return nil, nil
		case pdutype.PDUDataTransfer:
			pdu.readPDU()
			pdu.ms.SetPosition(1)
			err := pdu.Pdata.ReadDynamic(pdu.ms)
			if err != nil {
				return nil, err
			}
			if pdu.Pdata.MsgStatus > 0 {
				DCO := media.NewEmptyDCMObj()
				if !pdu.parseRawVRIntoDCM(DCO) {
					pdu.AbortRQ.Write(pdu.readWriter)
					return nil, errors.New("ERROR, pduservice::Read - ParseRawVRIntoDCM failed")
				}
				return DCO, nil
			}
			break
		case pdutype.AssociationReleaseRequest:
			log.Printf("INFO, ASSOC-R-RQ: %s --> %s\n", pdu.AssocRQ.GetCallingAE(), pdu.AssocRQ.GetCalledAE())
			pdu.ReleaseRQ.ReadDynamic(pdu.ms)
			pdu.ReleaseRP.Write(pdu.readWriter)
			return nil, errors.New("ERROR, pduservice::Read - A-Release-RQ")
		case pdutype.AssociationReleaseResponse:
			log.Printf("INFO, ASSOC-R-RP: %s <-- %s\n", pdu.AssocRQ.GetCallingAE(), pdu.AssocRQ.GetCalledAE())
			return nil, errors.New("ERROR, pduservice::Read - A-Release-RP")
		case pdutype.AssociationAbortRequest:
			log.Printf("INFO, ASSOC-ABORT-RQ: %s --> %s\n", pdu.AssocRQ.GetCallingAE(), pdu.AssocRQ.GetCalledAE())
			return nil, errors.New("ERROR, pduservice::Read - A-Abort-RQ")
		default:
			pdu.AbortRQ.Write(pdu.readWriter)
			return nil, errors.New("ERROR, pduservice::Read - unknown ItemType")
		}
	}
}

func (pdu *pduService) GetAAssociationRQ() AAssociationRQ {
	return pdu.AssocRQ
}

func (pdu *pduService) GetCalledAE() string {
	return pdu.AssocRQ.GetCalledAE()
}

func (pdu *pduService) GetCallingAE() string {
	return pdu.AssocRQ.GetCallingAE()
}

func (pdu *pduService) SetCalledAE(calledAE string) {
	pdu.AssocRQ.SetCalledAE(calledAE)
}

func (pdu *pduService) SetCallingAE(callingAE string) {
	pdu.AssocRQ.SetCallingAE(callingAE)
}

func (pdu *pduService) AddPresContexts(presentationContext PresentationContext) {
	pdu.AssocRQ.AddPresContexts(presentationContext)
}

func (pdu *pduService) GetPresentationContextID() byte {
	return pdu.Pdata.PresentationContextID
}

func (pdu *pduService) SetOnAssociationRequest(f func(request AAssociationRQ) bool) {
	pdu.OnAssociationRequest = f
}

func (pdu *pduService) Write(DCO media.DcmObj, SOPClass string, ItemType byte) error {
	if pdu.Pdata.Buffer != nil {
		pdu.Pdata.Buffer.ClearMemoryStream()
	} else {
		pdu.Pdata.Buffer = media.NewEmptyBufData()
	}

	if pdu.Pdata.PresentationContextID == 0 {
		return errors.New("ERROR, pduservice::Write - PresentationContextID==0")
	}
	if ItemType == 0x01 {
		if pdu.parseDCMIntoRaw(DCO) == false {
			return errors.New("ERROR, pduservice::Write - ParseDCMIntoRaw failed")
		}
	} else {
		if pdu.parseDCMIntoRaw(DCO) == false {
			return errors.New("ERROR, pduservice::Write - ParseDCMIntoRaw failed")
		}
	}

	pdu.Pdata.MsgHeader = ItemType
	if pdu.AssocAC.GetUserInformation().GetMaxSubLength().GetMaximumLength() > maxPduLength {
		pdu.AssocAC.SetMaxSubLength(maxPduLength)
	}

	// Fixed MaxLength - 6 20200811
	pdu.Pdata.BlockSize = pdu.AssocAC.GetMaxSubLength() - 6

	log.Printf("INFO, PDU-Service: %s --> %s", SOPClass, pdu.GetCallingAE())

	return pdu.Pdata.Write(pdu.readWriter)
}

func (pdu *pduService) interogateAAssociateAC() bool {
	var PresentationContextID byte
	TS := ""

	for _, presContextAccept := range pdu.AssocAC.GetPresContextAccepts() {
		if presContextAccept.GetResult() == 0 {
			pdu.AcceptedPresentationContexts = append(pdu.AcceptedPresentationContexts, presContextAccept)
			if len(TS) == 0 {
				if presContextAccept.GetTrnSyntax().UIDName == uid.ExplicitVRLittleEndian {
					TS = presContextAccept.GetTrnSyntax().UIDName
					PresentationContextID = presContextAccept.GetPresentationContextID()
				}
			}
			if len(TS) == 0 {
				if presContextAccept.GetTrnSyntax().UIDName == uid.ImplicitVRLittleEndian {
					TS = presContextAccept.GetTrnSyntax().UIDName
					PresentationContextID = presContextAccept.GetPresentationContextID()
				}
			}
		}
	}
	if (len(TS) > 0) && (len(pdu.AcceptedPresentationContexts) > 0) {
		pdu.Pdata.PresentationContextID = PresentationContextID
		return true
	}
	return false
}

func (pdu *pduService) interogateAAssociateRQ(rw *bufio.ReadWriter) error {
	if pdu.OnAssociationRequest == nil || !pdu.OnAssociationRequest(pdu.AssocRQ) {
		pdu.AssocRJ.Set(1, 7)
		return pdu.AssocRJ.Write(rw)
	}

	pdu.AssocAC.SetCalledAE(pdu.AssocRQ.GetCalledAE())
	pdu.AssocAC.SetCallingAE(pdu.AssocRQ.GetCallingAE())

	pdu.AssocAC.SetAppContext(pdu.AssocRQ.GetAppContext())

	pdu.AssocAC.SetUserInformation(pdu.AssocRQ.GetUserInformation())

	for _, PresContext := range pdu.AssocRQ.GetPresContexts() {
		log.Printf("INFO, ASSOC-RQ: \tPresentation Context %s\n", PresContext.GetAbstractSyntax().UIDName)
		for _, TransferSyn := range PresContext.GetTransferSyntaxes() {
			log.Printf("INFO, ASSOC-RQ: \t\tTransfer Synxtax %s\n", TransferSyn.UIDName)
		}

		PresContextAccept := NewPresentationContextAccept()
		PresContextAccept.SetResult(4)
		PresContextAccept.SetTransferSyntax("")
		PresContextAccept.SetAbstractSyntax(PresContext.GetAbstractSyntax().UIDName)
		TS := ""
		for _, TrnSyntax := range PresContext.GetTransferSyntaxes() {
			if TrnSyntax.UIDName == uid.ExplicitVRLittleEndian {
				TS = TrnSyntax.UIDName
			}
		}
		if TS == "" {
			for _, TrnSyntax := range PresContext.GetTransferSyntaxes() {
				if TrnSyntax.UIDName == uid.ImplicitVRLittleEndian {
					TS = TrnSyntax.UIDName
				}
			}
		}
		if len(TS) > 0 {
			PresContextAccept.SetResult(0)
			PresContextAccept.SetTransferSyntax(TS)
			PresContextAccept.SetPresentationContextID(PresContext.GetPresentationContextID())
			pdu.AcceptedPresentationContexts = append(pdu.AcceptedPresentationContexts, PresContextAccept)
		}
		pdu.AssocAC.AddPresContextAccept(PresContextAccept)
	}

	if len(pdu.AcceptedPresentationContexts) > 0 {
		MaxSubLength := NewMaximumSubLength()
		UserInfo := NewUserInformation()

		MaxSubLength.SetMaximumLength(maxPduLength)
		UserInfo.SetImpClassUID(imp.GetImpClassUID())
		UserInfo.SetImpVersionName(imp.GetImpVersion())
		UserInfo.SetMaxSubLength(MaxSubLength)
		pdu.AssocAC.SetUserInformation(UserInfo)
		return pdu.AssocAC.Write(rw)
	}

	log.Println("ERROR, pduservice::InterogateAAssociateRQ - No valid AcceptedPresentationContexts")
	return pdu.AssocRJ.Write(rw)
}

func (pdu *pduService) parseDCMIntoRaw(DCO media.DcmObj) bool {
	pdu.Pdata.Buffer.WriteObj(DCO)
	return true
}

func (pdu *pduService) parseRawVRIntoDCM(DCO media.DcmObj) bool {
	TrnSyntax := pdu.GetTransferSyntaxUID(pdu.Pdata.PresentationContextID)
	if len(TrnSyntax) == 0 {
		log.Println("ERROR, pduservice::ParseRawVRIntoDCM - Transfer syntax length is 0")
		return false
	}
	DCO.SetTransferSyntax(TrnSyntax)
	if TrnSyntax == uid.ExplicitVRLittleEndian {
		DCO.SetExplicitVR(true)
	}
	if TrnSyntax == uid.ExplicitVRBigEndian {
		DCO.SetBigEndian(true)
	}
	pdu.Pdata.Buffer.SetPosition(0)
	return pdu.Pdata.Buffer.ReadObj(DCO)
}

func (pdu *pduService) readPDU() {
	pdu.ms.ReadFully(pdu.readWriter, int(pdu.pdulength)-4)
}
