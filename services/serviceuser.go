package services

import (
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

	pdu :=network.NewPDUService()
	if OpenAssociation(pdu, LAET, RAET, RIP, RPort, "1.2.840.10008.1.1", timeout) {
		if dimsec.CEchoWriteRQ(*pdu, "1.2.840.10008.1.1") {
			if dimsec.CEchoReadRSP(*pdu) {
				flag = true
			}
		}
	}
	pdu.Close()
	return flag
}

func WriteStoreRQ(pdu network.PDUService, DDO media.DcmObj, SOPClassUID string) int {
	status:=-1

	PCID:=pdu.Pdata.PresentationContextID;
	if(PCID==0) {
		return -1
	}
	TrnSyntOUT:=pdu.GetTransferSyntaxUID(PCID)

	if(len(TrnSyntOUT)==0) {
	return -2
	}

	if TrnSyntOUT==DDO.TransferSyntax {
		if(dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID)){
			status=0
		} else {
			status=-4
		}
	} else{
		DDO.TransferSyntax=TrnSyntOUT
		DDO.ExplicitVR=true
		DDO.BigEndian=false
		if TrnSyntOUT=="1.2.840.10008.1.2" {
			DDO.ExplicitVR=false
		}
		if TrnSyntOUT=="1.2.840.10008.1.2.2" {
			DDO.BigEndian=true
		}
		if dimsec.CStoreWriteRQ(pdu, DDO, SOPClassUID) {
			  status=0
		} else {
			  status=-4
		}
	}

	return status
}

func StoreSCU(LAET string, RAET string, RIP string, RPort string, FileName string, timeout int) bool{
	flag:=false
	var DDO media.DcmObj

	if DDO.Read(FileName) {
		SOPClassUID:=DDO.GetString(0x08, 0x16)
		if len(SOPClassUID) > 0 {
			pdu:=network.NewPDUService()
			if OpenAssociation(pdu, LAET, RAET, RIP, RPort, SOPClassUID, timeout) {
				if(WriteStoreRQ(*pdu, DDO, SOPClassUID)==0x00){
					if(dimsec.CStoreReadRSP(*pdu)==0x00) {
						flag=true
					}
				}
			}
			pdu.Close()
		}
	}
	return flag
}

func FindSCU(LAET string, RAET string, RIP string, RPort string, Query media.DcmObj, Results []media.DcmObj, timeout int) bool{
	flag:=false
	status:=0
	var DDO media.DcmObj
	SOPClassUID:="1.2.840.10008.5.1.4.1.2.2.1"

	pdu:=network.NewPDUService()
	if OpenAssociation(pdu, LAET, RAET, RIP, RPort, SOPClassUID, timeout) {
		if(CFindWriteRQ(pdu, Query, SOPClassUID)){
			for(status!=-1){
				status=CFindReadRSP(pdu, &DDO)
				if status!=-1 {
					Results = append(Results, DDO)
				}
			}
		}
	}
	return flag
}