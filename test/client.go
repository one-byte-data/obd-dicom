package main

import (
	"fmt"

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
	} else {
		fmt.Println("DICOM Echo Failed!!")
	}

}
