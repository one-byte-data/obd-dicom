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
	SetOnAssociationRequest(f func(calledAE string) bool)
	SetOnCFindRequest(f func(findLevel string, data media.DcmObj, Result media.DcmObj))
	SetOnCMoveRequest(f func(moveLevel string, data media.DcmObj))
	SetOnCStoreRequest(f func(data media.DcmObj))
	handleConnection(conn net.Conn)
}

type scp struct {
	CalledAEs            []string
	Port                 int
	OnAssociationRequest func(calledAE string) bool
	OnCFindRequest       func(findLevel string, data media.DcmObj, Result media.DcmObj)
	OnCMoveRequest       func(moveLevel string, data media.DcmObj)
	OnCStoreRequest      func(data media.DcmObj)
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

	if s.OnAssociationRequest != nil {
		pdu.SetOnAssociationRequest(s.OnAssociationRequest)
	}

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
				log.Printf("ERROR, handleConnection, C-Store failed to read request : %s", err.Error())
				return
			}

			if s.OnCStoreRequest != nil {
				s.OnCStoreRequest(DDO)
			}

			err = dimsec.CStoreWriteRSP(pdu, DCO, 0)
			if err != nil {
				log.Printf("ERROR, handleConnection, C-Store failed to write response: %s", err.Error())
				return
			}
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

			Result := media.NewEmptyDCMObj()

			if s.OnCFindRequest != nil {
				s.OnCFindRequest(QueryLevel, DDO, Result)
			}

			err = dimsec.CFindWriteRSP(pdu, DCO, Result, 0x00)
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

			if s.OnCMoveRequest != nil {
				s.OnCMoveRequest(MoveLevel, DDO)
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

func (s *scp) SetOnAssociationRequest(f func(calledAE string) bool) {
	s.OnAssociationRequest = f
}

func (s *scp) SetOnCFindRequest(f func(findLevel string, data media.DcmObj, Result media.DcmObj)) {
	s.OnCFindRequest = f
}

func (s *scp) SetOnCMoveRequest(f func(moveLevel string, data media.DcmObj)) {
	s.OnCMoveRequest = f
}

func (s *scp) SetOnCStoreRequest(f func(data media.DcmObj)) {
	s.OnCStoreRequest = f
}
