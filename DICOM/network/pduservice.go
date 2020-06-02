package network

import (
	"net"
	"rafael/DICOM/media"
)

type PDUService struct {
	AcceptedPresentationContexts []PresentationContextAccept
	conn                         net.Conn
	assocRQ                      AAssociationRQ
	assocAC                      AAssociationAC
	assocRJ                      AAssociationRJ
	releaseRQ                    AReleaseRQ
	releaseRP                    AReleaseRP
	abortRQ                      AAbortRQ
	pdata                        PDataTF
}

func (pdu *PDUService) InterogateAAssociateAC() bool {
	var PresentationContextID byte
	flag := false
	TS := ""

	for i := 0; i < len(pdu.assocAC.PresContextAccepts); i++ {
		PresContextAccept := pdu.assocAC.PresContextAccepts[i]
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
		pdu.pdata.PresentationContextID = PresentationContextID
		flag = true
	}
	return flag
}

func (pdu *PDUService) ParseDCMIntoRaw(DCO media.DcmObj) bool {
	pdu.pdata.Buffer.WriteObj(&DCO)
	return true
}

func (pdu *PDUService) Write(DCO media.DcmObj, SOPClass string, ItemType byte) bool {
	if pdu.pdata.PresentationContextID == 0 {
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
	pdu.pdata.MsgHeader = ItemType
	if pdu.assocAC.UserInfo.MaxSubLength.MaximumLength > 16384 {
		pdu.assocAC.UserInfo.MaxSubLength.MaximumLength = 16384
	}
	pdu.pdata.BlockSize = pdu.assocAC.UserInfo.MaxSubLength.MaximumLength
	return pdu.pdata.Write(pdu.conn)
}

func (pdu *PDUService) GetTransferSyntaxUID(pcid byte) string {
	for i := 0; i < len(pdu.AcceptedPresentationContexts); i++ {
		pca := pdu.AcceptedPresentationContexts[i]
		if pca.PresentationContextID == pcid {
			return pca.TrnSyntax.UIDName
		}
	}
	return ""
}

func (pdu *PDUService) ParseRawVRIntoDCM(DCO *media.DcmObj) bool {

	TrnSyntax := pdu.GetTransferSyntaxUID(pdu.pdata.PresentationContextID)
	if len(TrnSyntax) == 0 {
		return false
	}
	DCO.TransferSyntax = TrnSyntax
	pdu.pdata.Buffer.Ms.Position = 0
	return pdu.pdata.Buffer.ReadObj(DCO)
}

func (pdu *PDUService) Read(DCO *media.DcmObj) bool {
	pdu.pdata.MsgStatus = 0
	if pdu.pdata.Length != 0 {
		pdu.pdata.ReadDynamic(pdu.conn)
		if pdu.pdata.MsgStatus > 0 {
			if pdu.ParseRawVRIntoDCM(DCO) == false {
				pdu.abortRQ.Write(pdu.conn)
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
			pdu.assocRQ.Read(pdu.conn)
			pdu.abortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
			break
		case 0x02: // A-Associate-AC, should not get here
			pdu.assocAC.Read(pdu.conn)
			pdu.abortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
			break
		case 0x03: // A-Associate-RJ, should not get here
			pdu.abortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
			break
		case 0x04: // P-Data-TF
			pdu.pdata.ReadDynamic(pdu.conn)
			if pdu.pdata.MsgStatus > 0 {
				if pdu.ParseRawVRIntoDCM(DCO) == false {
					pdu.abortRQ.Write(pdu.conn)
					pdu.conn.Close()
					return false
				}
				return true
			}
			break
		case 0x05: // A-Release-RQ
			pdu.releaseRQ.ReadDynamic(pdu.conn)
			pdu.releaseRP.Write(pdu.conn)
			return false
			break
		case 0x06: // A-Release-RP
			pdu.conn.Close()
			return false
			break
		case 0x07: //A-Abort-RQ
			pdu.conn.Close()
			return false
			break
		default:
			pdu.abortRQ.Write(pdu.conn)
			pdu.conn.Close()
			return false
		}
	}
	return false
}

func (pdu *PDUService) Connect(IP string, Port string) bool {
	conn, err := net.Dial("tcp", IP+":"+Port)
	if err != nil {
		return false
	}

	pdu.conn = conn
	Resetuniq()
	pdu.assocRQ = *NewAAAssociationRQ()
	copy(pdu.assocRQ.CallingApTitle[:], "TESTSCU")
	copy(pdu.assocRQ.CalledApTitle[:], "CHARRUAPACS")

	PresContext := NewPresentationContext()
	PresContext.SetAbstractSyntax("1.2.840.10008.1.1") // DICOM-Echo
	PresContext.AddTransferSyntax("1.2.840.10008.1.2")
	pdu.assocRQ.PresContexts = append(pdu.assocRQ.PresContexts, *PresContext)

	pdu.assocRQ.UserInfo = *NewUserInformation()
	pdu.assocRQ.UserInfo.MaxSubLength.MaximumLength = 16384
	pdu.assocRQ.UserInfo.SetImpClassUID("1.2.826.0.1.3680043.2.1396.999")
	pdu.assocRQ.UserInfo.SetImpVersionName("CharruaSoft")
	if pdu.assocRQ.Write(pdu.conn) {
		var ItemType byte
		pdu.assocAC = *NewAAssociationAC()
		ItemType = ReadByte(pdu.conn)
		switch ItemType {
		case 0x02:
			pdu.assocAC.ReadDynamic(pdu.conn)
			if !pdu.InterogateAAssociateAC() {
				pdu.conn.Close()
				return false
			}
			return true
			break
		case 0x03:
			// Error, Assoc. Rejected.
			pdu.assocRJ.ReadDynamic(pdu.conn)
			pdu.conn.Close()
			return false
			break
		default:
			// Error, Corrupt Transmission
			pdu.conn.Close()
			return false
		}
	}
	//Indeterminate state.
	return false
}
