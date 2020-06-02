package main

import (
	"fmt"
	"rafael/DICOM/dimsec"
	"rafael/DICOM/network"
	"rafael/DICOM/media"
)

func main() {
	var pdu network.PDUService

	media.InitDict()
	if pdu.Connect("localhost", "1040") {
		fmt.Println("Connection Success")
		if dimsec.CEchoWriteRQ(pdu, "1.2.840.10008.1.1") {
			if dimsec.CEchoReadRSP(pdu) {
				fmt.Println("DICOM Echo OK!!")
			} else {
				fmt.Println("DICOM Echo Failed!!")
			}
		}
	} else {
		fmt.Println("Connection failed")
	}
}
