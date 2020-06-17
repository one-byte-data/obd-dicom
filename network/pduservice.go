package network

import (
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
)

// PDUService - struct for PDUService
type PDUService struct {
	AcceptedPresentationContexts []PresentationContextAccept
	conn                         net.Conn
	AssocRQ                      AAssociationRQ
	AssocAC                      AAssociationAC
	AssocRJ                      AAssociationRJ
	ReleaseRQ                    AReleaseRQ
	ReleaseRP                    AReleaseRP
	AbortRQ                      AAbortRQ
	Pdata                        PDataTF
}

// NewPDUService - creates a pointer to PDUService
func NewPDUService() *PDUService {
	return &PDUService{
		AssocRQ:   *NewAAssociationRQ(),
		AssocAC:   *NewAAssociationAC(),
		AssocRJ:   *NewAAssociationRJ(),
		ReleaseRQ: *NewAReleaseRQ(),
		ReleaseRP: *NewAReleaseRP(),
		AbortRQ:   *NewAAbortRQ(),
	}
}

// InterogateAAssociateAC - InterogateAAssociateAC
func (pdu *PDUService) InterogateAAssociateAC() bool {
	var PresentationContextID byte
	flag := false
	TS := ""

	for i := 0; i < len(pdu.AssocAC.PresContextAccepts); i++ {
		PresContextAccept := pdu.AssocAC.PresContextAccepts[i]
		if PresContextAccept.Result == 0 {
			pdu.AcceptedPresentationContexts = append(pdu.AcceptedPresentationContexts, PresContextAccept)
			if len(TS) == 0 {
				if PresContextAccept.TrnSyntax.UIDName == "1.2.840.10008.1.2.1" {
					TS = PresContextAccept.TrnSyntax.UIDName
					PresentationContextID = PresContextAccept.PresentationContextID
				}
			}
			if len(TS) == 0 {
				if PresContextAccept.TrnSyntax.UIDName == "1.2.840.10008.1.2" {
					TS = PresContextAccept.TrnSyntax.UIDName
					PresentationContextID = PresContextAccept.PresentationContextID
				}
			}
		}
	}
	if (len(TS) > 0) && (len(pdu.AcceptedPresentationContexts) > 0) {
		pdu.Pdata.PresentationContextID = PresentationContextID
		flag = true
	}
	return flag
}

// InterogateAAssociateRQ - InterogateAAssociateRQ
func (pdu *PDUService) InterogateAAssociateRQ(conn net.Conn) bool {
	pdu.AssocAC.CalledApTitle = pdu.AssocRQ.CalledApTitle
	pdu.AssocAC.CallingApTitle = pdu.AssocRQ.CallingApTitle
	pdu.AssocAC.AppContext = pdu.AssocRQ.AppContext
	pdu.AssocAC.UserInfo = pdu.AssocRQ.UserInfo
	for i := 0; i < len(pdu.AssocRQ.PresContexts); i++ {
		PresContext := pdu.AssocRQ.PresContexts[i]
		PresContextAccept := *NewPresentationContextAccept()
		PresContextAccept.Result = 4
		PresContextAccept.SetTransferSyntax("")
		PresContextAccept.SetAbstractSyntax(PresContext.AbsSyntax.UIDName)
		TS := ""
		for j := 0; j < len(PresContext.TrnSyntaxs); j++ {
			TrnSyntax := PresContext.TrnSyntaxs[j]
			if TrnSyntax.UIDName == "1.2.840.10008.1.2.1" {
				TS = TrnSyntax.UIDName
			}
		}
		if TS == "" {
			for j := 0; j < len(PresContext.TrnSyntaxs); j++ {
				TrnSyntax := PresContext.TrnSyntaxs[j]
				if TrnSyntax.UIDName == "1.2.840.10008.1.2" {
					TS = TrnSyntax.UIDName
				}
			}
		}
		if len(TS) > 0 {
			PresContextAccept.Result = 0
			PresContextAccept.SetTransferSyntax(TS)
			PresContextAccept.PresentationContextID = PresContext.PresentationContextID
			pdu.AcceptedPresentationContexts = append(pdu.AcceptedPresentationContexts, PresContextAccept)
		}
		pdu.AssocAC.PresContextAccepts = append(pdu.AssocAC.PresContextAccepts, PresContextAccept)
	}
	if len(pdu.AcceptedPresentationContexts) > 0 {
		MaxSubLength := *NewMaximumSubLength()
		UserInfo := *NewUserInformation()

		MaxSubLength.MaximumLength = 16384
		UserInfo.SetImpClassUID("1.2.826.0.1.3680043.10.90.999")
		UserInfo.SetImpVersionName("One-Byte-Data")
		UserInfo.MaxSubLength = MaxSubLength
		pdu.AssocAC.SetUserInformation(UserInfo)
		return pdu.AssocAC.Write(conn)

	}
	pdu.AssocAC.Write(conn)
	return false
}

// ParseDCMIntoRaw - ParseDCMIntoRaw
func (pdu *PDUService) ParseDCMIntoRaw(DCO media.DcmObj) bool {
	pdu.Pdata.Buffer.WriteObj(&DCO)
	return true
}

func (pdu *PDUService) Write(DCO media.DcmObj, SOPClass string, ItemType byte) bool {
	// Limpiar el buffer aqui.
	pdu.Pdata.Buffer.Ms.Clear()
	if pdu.Pdata.PresentationContextID == 0 {
		return false
	}
	if ItemType == 0x01 {
		if pdu.ParseDCMIntoRaw(DCO) == false {
			return false
		}
	} else {
		if pdu.ParseDCMIntoRaw(DCO) == false {
			return false
		}
	}
	pdu.Pdata.MsgHeader = ItemType
	if pdu.AssocAC.UserInfo.MaxSubLength.MaximumLength > 16384 {
		pdu.AssocAC.UserInfo.MaxSubLength.MaximumLength = 16384
	}
	pdu.Pdata.BlockSize = pdu.AssocAC.UserInfo.MaxSubLength.MaximumLength
	return pdu.Pdata.Write(pdu.conn)
}

// GetTransferSyntaxUID - Gets transfer syntax UID
func (pdu *PDUService) GetTransferSyntaxUID(pcid byte) string {
	for i := 0; i < len(pdu.AcceptedPresentationContexts); i++ {
		pca := pdu.AcceptedPresentationContexts[i]
		if pca.PresentationContextID == pcid {
			return pca.TrnSyntax.UIDName
		}
	}
	return ""
}

// ParseRawVRIntoDCM - ParseRawVRIntoDCM
func (pdu *PDUService) ParseRawVRIntoDCM(DCO *media.DcmObj) bool {

	TrnSyntax := pdu.GetTransferSyntaxUID(pdu.Pdata.PresentationContextID)
	if len(TrnSyntax) == 0 {
		return false
	}
	DCO.TransferSyntax = TrnSyntax
	if TrnSyntax == "1.2.840.10008.1.2.1" {
		DCO.ExplicitVR = true
	}
	if TrnSyntax == "1.2.840.10008.1.2.2" {
		DCO.BigEndian = true
	}
	pdu.Pdata.Buffer.Ms.Position = 0
	return pdu.Pdata.Buffer.ReadObj(DCO)
}

func (pdu *PDUService) Read(DCO *media.DcmObj) bool {
	// Ver donde tendria que limpiar el buffer.
	pdu.Pdata.Buffer.Ms.Clear()
	pdu.Pdata.MsgStatus = 0
	if pdu.Pdata.Length != 0 {
		pdu.Pdata.ReadDynamic(pdu.conn)
		if pdu.Pdata.MsgStatus > 0 {
			if pdu.ParseRawVRIntoDCM(DCO) == false {
				pdu.AbortRQ.Write(pdu.conn)
				pdu.conn.Close()
				return false
			}
			return true
		}
	}
	for true {
		ItemType := ReadByte(pdu.conn)
		switch ItemType {
		case 0x01: // A-Associate-RQ, should not get here
			pdu.AssocRQ.Read(pdu.conn)
			pdu.AbortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
		case 0x02: // A-Associate-AC, should not get here
			pdu.AssocAC.Read(pdu.conn)
			pdu.AbortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
		case 0x03: // A-Associate-RJ, should not get here
			pdu.AbortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
		case 0x04: // P-Data-TF
			pdu.Pdata.ReadDynamic(pdu.conn)
			if pdu.Pdata.MsgStatus > 0 {
				if pdu.ParseRawVRIntoDCM(DCO) == false {
					pdu.AbortRQ.Write(pdu.conn)
					pdu.conn.Close()
					return false
				}
				return true
			}
			break
		case 0x05: // A-Release-RQ
			pdu.ReleaseRQ.ReadDynamic(pdu.conn)
			pdu.ReleaseRP.Write(pdu.conn)
			return false
		case 0x06: // A-Release-RP
			pdu.conn.Close()
			return false
		case 0x07: //A-Abort-RQ
			pdu.conn.Close()
			return false
		default:
			pdu.AbortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
		}
	}
	return false
}

// SetTimeout - SetTimeout
func (pdu *PDUService) SetTimeout(timeout int) {
}

// Connect - Connect
func (pdu *PDUService) Connect(IP string, Port string) bool {
	conn, err := net.Dial("tcp", IP+":"+Port)
	if err != nil {
		return false
	}

	pdu.conn = conn
	pdu.AssocRQ.UserInfo.MaxSubLength.MaximumLength = 16384
	pdu.AssocRQ.UserInfo.SetImpClassUID("1.2.826.0.1.3680043.10.90.999")
	pdu.AssocRQ.UserInfo.SetImpVersionName("One-Byte-Data")
	if pdu.AssocRQ.Write(pdu.conn) {
		var ItemType byte
		ItemType = ReadByte(pdu.conn)
		switch ItemType {
		case 0x02:
			pdu.AssocAC.ReadDynamic(pdu.conn)
			if !pdu.InterogateAAssociateAC() {
				pdu.conn.Close()
				return false
			}
			return true
		case 0x03:
			// Error, Assoc. Rejected.
			pdu.AssocRJ.ReadDynamic(pdu.conn)
			pdu.conn.Close()
			return false
		default:
			// Error, Corrupt Transmission
			pdu.conn.Close()
			return false
		}
	}
	//Indeterminate state.
	return false
}

// Close - close connection
func (pdu *PDUService) Close() {
	pdu.ReleaseRQ.Write(pdu.conn)
	pdu.ReleaseRP.Read(pdu.conn)
	pdu.conn.Close()
}

// Multiplex - Multiplex
func (pdu *PDUService) Multiplex(conn net.Conn) bool {
	pdu.conn = conn
	if pdu.AssocRQ.Read(pdu.conn) {
		if pdu.InterogateAAssociateRQ(pdu.conn) {
			return true
		}
		conn.Close()
	}
	return false
}
