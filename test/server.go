package main

import (
	"fmt"
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/dimsec"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

func handleConnection(conn net.Conn) {
	pdu := *network.NewPDUService()
	if pdu.Multiplex(conn) {
		var DCO media.DcmObj
		flag := true
		for flag && pdu.Read(&DCO) {
			command := DCO.GetUShort(0x00, 0x0100)
			switch command {
			case 0x01: // C-Store
				var DDO media.DcmObj
				if dimsec.CStoreReadRQ(pdu, DCO, &DDO) {
					if dimsec.CStoreWriteRSP(pdu, DCO, 0) {
						DDO.Write("test.dcm")
						flag = true
					} else {
						flag = false
					}
				}
				break
			case 0x20: // C-Find
				var DDO media.DcmObj
				if dimsec.CFindReadRQ(pdu, DCO, &DDO) {
					QueryLevel := DDO.GetString(0x08, 0x52) // Get Query Level
					var Out media.DcmObj                    // This is for the result
					if QueryLevel == "STUDY" {
						// Process Study Query
					}
					if QueryLevel == "SERIES" {
						// Process Series Query
					}
					if QueryLevel == "IMAGE" {
						// Process Image Query
					}
					dimsec.CFindWriteRSP(pdu, DCO, Out, 0x00)
				}
				break
			case 0x21: // C-Move
				var DDO media.DcmObj
				if dimsec.CMoveReadRQ(pdu, DCO, &DDO) {
					MoveLevel := DDO.GetString(0x08, 0x52) // Get Move Level
					if MoveLevel == "STUDY" {
						// Process Study Move
					}
					if MoveLevel == "SERIES" {
						// Process Series Move
					}
					if MoveLevel == "IMAGE" {
						// Process Image Move
					}
					dimsec.CMoveWriteRSP(pdu, DCO, 0x00, 0x00)
				}
				break
			case 0x30: // C-Echo
				if dimsec.CEchoReadRQ(pdu, DCO) {
					if dimsec.CEchoWriteRSP(pdu, DCO) {
						fmt.Println("Echo Success!")
					} else {
						fmt.Println("Echo failed!")
						flag = false
					}
				}
				break
			default:
				fmt.Println("ERROR, service not implemented: " + string(command))
				flag = false
			}
		}
	}
}

func server(Port string) {
	l, err := net.Listen("tcp4", ":"+Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	media.InitDict()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
