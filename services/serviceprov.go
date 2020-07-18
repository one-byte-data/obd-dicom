package services

import (
	"fmt"
	"log"
	"net"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/dimsec"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
)

// SCP - Interface to scp
type SCP interface {
	StartServer() error
	handleConnection(conn net.Conn)
}

type scp struct {
	CalledAEs []string
	Port      int
}

// NewSCP - Creates an interface to scu
func NewSCP(calledAEs []string, port int) SCP {
	media.InitDict()

	return &scp{
		CalledAEs: calledAEs,
		Port:      port,
	}
}

func (s *scp) StartServer() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *scp) handleConnection(conn net.Conn) {
	pdu := network.NewPDUService()
	log.Println("INFO, handleConnection, new connection from: ", conn.RemoteAddr())
	err := pdu.Multiplex(conn)
	if err != nil {
		log.Print(err)
		return
	}

	DCO := media.NewEmptyDCMObj()
	for true {
		err := pdu.Read(DCO)
		if err != nil {
			break
		}
		command := DCO.GetUShort(0x00, 0x0100)
		switch command {
		case 0x01: // C-Store
			DDO := media.NewEmptyDCMObj()
			err := dimsec.CStoreReadRQ(pdu, DCO, DDO)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Store failed to read request!")
				return
			}
			err = dimsec.CStoreWriteRSP(pdu, DCO, 0)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Store failed to write response!")
				return
			}
			DDO.WriteToFile("test.dcm")
			log.Println("INFO, handleConnection, CStore Success")
			break
		case 0x20: // C-Find
			DDO := media.NewEmptyDCMObj()
			err := dimsec.CFindReadRQ(pdu, DCO, DDO)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Find failed to read request!")
				break
			}
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
			err = dimsec.CFindWriteRSP(pdu, DCO, Out, 0x00)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Find failed to write response!")
				return
			}
			break
		case 0x21: // C-Move
			DDO := media.NewEmptyDCMObj()
			err := dimsec.CMoveReadRQ(pdu, DCO, DDO)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Move failed to read request!")
				return
			}
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
			err = dimsec.CMoveWriteRSP(pdu, DCO, 0x00, 0x00)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Move failed to write response!")
				return
			}
			break
		case 0x30: // C-Echo
			if dimsec.CEchoReadRQ(pdu, DCO) {
				err := dimsec.CEchoWriteRSP(pdu, DCO)
				if err != nil {
					log.Println("ERROR, handleConnection, C-Echo failed to write response!")
					return
				}
				log.Println("INFO, handleConnection, C-Echo Success!")
			}
			break
		default:
			log.Println("ERROR, handleConnection, service not implemented: " + string(command))
			return
		}
	}
}
