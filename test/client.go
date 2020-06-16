package main

import (
	"fmt"
	"strconv"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/services"
)

func client(AET string, IP string, Port string) {
	media.InitDict()
	if services.EchoSCU("TESTSCU", AET, IP, Port, 30) {
		fmt.Println("DICOM Echo OK!!")
		if services.StoreSCU("TESTSCU", AET, IP, Port, "test.dcm", 30) {
			fmt.Println("DICOM Store OK!!")
		} else {
			fmt.Println("DICOM Store Failed!!")
		}
		var Results []media.DcmObj
		var Query media.DcmObj
		Query.WriteString(0x08, 0x20, "DA", "")
		Query.WriteString(0x08, 0x30, "TM", "")
		Query.WriteString(0x08, 0x50, "SH", "")
		Query.WriteString(0x08, 0x52, "CS", "STUDY") // Use Study Level
		Query.WriteString(0x08, 0x61, "CS", "MR")    // Look for Modality=MR
		Query.WriteString(0x08, 0x1030, "LO", "")
		Query.WriteString(0x10, 0x10, "PN", "")
		Query.WriteString(0x10, 0x20, "LO", "")
		Query.WriteString(0x10, 0x30, "DA", "")
		Query.WriteString(0x10, 0x40, "CS", "")
		Query.WriteString(0x20, 0x0D, "UI", "")
		if services.FindSCU("TESTSCU", AET, IP, Port, Query, &Results, 30) == 0 {
			fmt.Println("DICOM Query OK!! Results: " + strconv.Itoa(len(Results)))
		} else {
			fmt.Println("DICOM Query Failed!!")
		}
		if services.MoveSCU("TESTSCU", AET, IP, Port, "DESTAET", Query, 30) == 0 {
			fmt.Println("DICOM Move OK!!")
		} else {
			fmt.Println("DICOM Move Failed!!")
		}
	} else {
		fmt.Println("DICOM Echo Failed!!")
	}

}
