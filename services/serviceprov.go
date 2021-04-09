package services

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"git.onebytedata.com/odb/go-dicom/dimsec"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/network/dicomcommand"
	"git.onebytedata.com/odb/go-dicom/network/dicomstatus"
	"git.onebytedata.com/odb/go-dicom/tags"
)

// SCP - Interface to scp
type SCP interface {
	StartServer() error
	SetOnAssociationRequest(f func(request network.AAssociationRQ) bool)
	SetOnCFindRequest(f func(request network.AAssociationRQ, findLevel string, data media.DcmObj) []media.DcmObj)
	SetOnCMoveRequest(f func(request network.AAssociationRQ, moveLevel string, data media.DcmObj))
	SetOnCStoreRequest(f func(request network.AAssociationRQ, data media.DcmObj))
	handleConnection(conn net.Conn)
}

type scp struct {
	CalledAEs            []string
	Port                 int
	OnAssociationRequest func(request network.AAssociationRQ) bool
	OnCFindRequest       func(request network.AAssociationRQ, findLevel string, data media.DcmObj) []media.DcmObj
	OnCMoveRequest       func(request network.AAssociationRQ, moveLevel string, data media.DcmObj)
	OnCStoreRequest      func(request network.AAssociationRQ, data media.DcmObj)
}

// NewSCP - Creates an interface to scu
func NewSCP(port int) SCP {
	media.InitDict()

	return &scp{
		Port: port,
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
		log.Println("INFO, handleConnection, new connection from: ", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *scp) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	pdu := network.NewPDUService()
	pdu.SetConn(rw)

	if s.OnAssociationRequest != nil {
		pdu.SetOnAssociationRequest(s.OnAssociationRequest)
	}

	var err error
	var dco media.DcmObj
	for err == nil {
		dco, err = pdu.NextPDU()
		if dco == nil {
			continue
		}
		command := dco.GetUShort(tags.CommandField)
		switch command {
		case dicomcommand.CStoreRequest:
			ddo, err := dimsec.CStoreReadRQ(pdu, dco)
			if err != nil {
				log.Printf("ERROR, handleConnection, C-Store failed to read request : %s", err.Error())
				conn.Close()
				return
			}

			if s.OnCStoreRequest != nil {
				s.OnCStoreRequest(pdu.GetAAssociationRQ(), ddo)
			}

			err = dimsec.CStoreWriteRSP(pdu, dco, 0)
			if err != nil {
				log.Printf("ERROR, handleConnection, C-Store failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Store Success")
			break
		case dicomcommand.CFindRequest:
			ddo, err := dimsec.CFindReadRQ(pdu)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Find failed to read request!")
				conn.Close()
				return
			}
			QueryLevel := ddo.GetString(tags.QueryRetrieveLevel)

			Results := make([]media.DcmObj, 0)

			if s.OnCFindRequest != nil {
				Results = s.OnCFindRequest(pdu.GetAAssociationRQ(), QueryLevel, ddo)
			}

			for _, result := range Results {
				err = dimsec.CFindWriteRSP(pdu, dco, result, dicomstatus.Pending)
				if err != nil {
					log.Printf("ERROR, handleConnection, C-Find failed to write response: %s", err.Error())
					conn.Close()
					return
				}
			}

			err = dimsec.CFindWriteRSP(pdu, dco, dco, dicomstatus.Success)
			if err != nil {
				log.Printf("ERROR, handleConnection, C-Find failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Find Success")
			break
		case dicomcommand.CMoveRequest:
			ddo, err := dimsec.CMoveReadRQ(pdu)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Move failed to read request!")
				conn.Close()
				return
			}
			MoveLevel := ddo.GetString(tags.QueryRetrieveLevel)

			if s.OnCMoveRequest != nil {
				s.OnCMoveRequest(pdu.GetAAssociationRQ(), MoveLevel, ddo)
			}

			err = dimsec.CMoveWriteRSP(pdu, dco, dicomstatus.Success, 0x00)
			if err != nil {
				log.Printf("ERROR, handleConnection, C-Move failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Move Success")
			break
		case dicomcommand.CEchoRequest:
			if dimsec.CEchoReadRQ(pdu, dco) {
				err := dimsec.CEchoWriteRSP(pdu, dco)
				if err != nil {
					log.Println("ERROR, handleConnection, C-Echo failed to write response!")
					conn.Close()
					return
				}
				log.Println("INFO, handleConnection, C-Echo Success!")
			}
			break
		default:
			log.Println("ERROR, handleConnection, service not implemented: " + string(command))
			conn.Close()
			return
		}
	}

	if err != nil {
		conn.Close()
	}
}

func (s *scp) SetOnAssociationRequest(f func(request network.AAssociationRQ) bool) {
	s.OnAssociationRequest = f
}

func (s *scp) SetOnCFindRequest(f func(request network.AAssociationRQ, findLevel string, data media.DcmObj) []media.DcmObj) {
	s.OnCFindRequest = f
}

func (s *scp) SetOnCMoveRequest(f func(request network.AAssociationRQ, moveLevel string, data media.DcmObj)) {
	s.OnCMoveRequest = f
}

func (s *scp) SetOnCStoreRequest(f func(request network.AAssociationRQ, data media.DcmObj)) {
	s.OnCStoreRequest = f
}
