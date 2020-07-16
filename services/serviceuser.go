package services

import (
	"log"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/dimsec"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

func OpenAssociation(pdu *network.PDUService, LAET string, RAET string, RIP string, RPort string, AbstractSyntax string, timeout int) bool {
	// DICOM Application Context
	pdu.AssocRQ.SetCallingApTitle(LAET)
	pdu.AssocRQ.SetCalledApTitle(RAET)
	pdu.SetTimeout(timeout)

	network.Resetuniq()
	PresContext := network.NewPresentationContext()
	PresContext.SetAbstractSyntax(AbstractSyntax)
	PresContext.AddTransferSyntax("1.2.840.10008.1.2")
	PresContext.AddTransferSyntax("1.2.840.10008.1.2.1")
	pdu.AssocRQ.PresContexts = append(pdu.AssocRQ.PresContexts, *PresContext)

	return (pdu.Connect(RIP, RPort))
}

func EchoSCU(LAET string, RAET string, RIP string, RPort string, timeout int) bool {
	flag := false

	pdu := network.NewPDUService()
	if OpenAssociation(pdu, LAET, RAET, RIP, RPort, "1.2.840.10008.1.1", timeout) {
		if dimsec.CEchoWriteRQ(*pdu, "1.2.840.10008.1.1") {
			if dimsec.CEchoReadRSP(*pdu) {
				flag = true
			} else {
				log.Println("ERROR, serviceuser::EchoSCU, dimsec.CEchoReadRSP failed")
			}
		} else {
			log.Println("ERROR, serviceuser::EchoSCU, dimsec.CEchoWriteRQ failed")
		}
	} else {
		log.Println("ERROR, serviceuser::EchoSCU, OpenAssociation failed, RAET: "+RAET)
	}
	pdu.Close()
	return flag
}

func WriteStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) int {
	status := -1

	PCID := pdu.Pdata.PresentationContextID
	if PCID == 0 {
		log.Println("ERROR, serviceuser::WriteStoreRQ, PCID==0")
		return -1
	}
	TrnSyntOUT := pdu.GetTransferSyntaxUID(PCID)

	if len(TrnSyntOUT) == 0 {
		log.Println("ERROR, serviceuser::WriteStoreRQ, TrnSyntOut is empty")
		return -2
	}

	if TrnSyntOUT == DDO.TransferSyntax {
		if dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID) {
			status = 0
		} else {
			log.Println("ERROR, serviceuser::WriteStoreRQ, dimsec.CStoreWriteRQ failed")
			status = -4
		}
	} else {
		DDO.TransferSyntax = TrnSyntOUT
		DDO.ExplicitVR = true
		DDO.BigEndian = false
		if TrnSyntOUT == "1.2.840.10008.1.2" {
			DDO.ExplicitVR = false
		}
		if TrnSyntOUT == "1.2.840.10008.1.2.2" {
			DDO.BigEndian = true
		}
		if dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID) {
			status = 0
		} else {
			log.Println("ERROR, serviceuser::WriteStoreRQ, dimsec.CStoreWriteRQ failed")
			status = -4
		}
	}

	return status
}

func StoreSCU(LAET string, RAET string, RIP string, RPort string, FileName string, timeout int) bool {
	flag := false
	var DDO media.DcmObj

	if DDO.Read(FileName) {
		SOPClassUID := DDO.GetString(0x08, 0x16)
		if len(SOPClassUID) > 0 {
			pdu := network.NewPDUService()
			if OpenAssociation(pdu, LAET, RAET, RIP, RPort, SOPClassUID, timeout) {
				if WriteStoreRQ(*pdu, DDO, SOPClassUID) == 0x00 {
					if dimsec.CStoreReadRSP(*pdu) == 0x00 {
						flag = true
					} else {
						log.Println("ERROR, serviceuser::StoreSCU, dimsec.CStoreReadRSP failed")
					}
				} else {
					log.Println("ERROR, serviceuser::StoreSCU, dimsec.CStoreWriteRQ failed")
				}
			} else {
				log.Println("ERROR, serviceuser::StoreSCU, OpenAssociation failed, RAET: "+RAET)
			}
			pdu.Close()
		} else {
			log.Println("ERROR, serviceuser::StoreSCU, SOPClassUID is empty")			
		}
	} else {
		log.Println("ERROR, serviceuser::StoreSCU, DDO.Read failed for: "+FileName)
	}
	return flag
}

func FindSCU(LAET string, RAET string, RIP string, RPort string, Query media.DcmObj, Results *[]media.DcmObj, timeout int) int {
	status := 1
	var DDO media.DcmObj
	SOPClassUID := "1.2.840.10008.5.1.4.1.2.2.1"

	pdu := network.NewPDUService()
	if OpenAssociation(pdu, LAET, RAET, RIP, RPort, SOPClassUID, timeout) {
		if dimsec.CFindWriteRQ(*pdu, Query, SOPClassUID) {
			for (status!=-1) && (status!=0) {
				status = dimsec.CFindReadRSP(*pdu, &DDO)
				if (status==0xFF00)||(status==0xFF01) {
					*Results = append(*Results, DDO)
				} 
				DDO.Clear();
			}
		} else {
			log.Println("ERROR, serviceuser::FindSCU, dimsec.CFindWriteRQ failed")
		}
	} else {
		log.Println("ERROR, serviceuser::FindSCU, OpenAssociation failed, RAET: "+RAET)
	}
	pdu.Close()
	return status
}

func MoveSCU(LAET string, RAET string, RIP string, RPort string, destAET string, Query media.DcmObj, timeout int) int {
	var pending int
	status := 0xFF00
	SOPClassUID := "1.2.840.10008.5.1.4.1.2.2.2"

	pdu := network.NewPDUService()
	if OpenAssociation(pdu, LAET, RAET, RIP, RPort, SOPClassUID, timeout) {
		if dimsec.CMoveWriteRQ(*pdu, Query, SOPClassUID, destAET) {
			var DDO media.DcmObj
			for status == 0xFF00 {
				status = dimsec.CMoveReadRSP(*pdu, &DDO, &pending)
				DDO.Clear()
			}
		} else {
			log.Println("ERROR, serviceuser::MoveSCU, dimsec.CMoveWriteRQ failed")
		}
	} else {
		log.Println("ERROR, serviceuser::MoveSCU, OpenAssociation failed, RAET: "+RAET)
	}
	pdu.Close()
	return status
}
