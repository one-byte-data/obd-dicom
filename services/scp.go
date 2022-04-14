package services

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"git.onebytedata.com/odb/go-dicom/dictionary/tags"
	"git.onebytedata.com/odb/go-dicom/dimsec"
	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/network/dicomcommand"
	"git.onebytedata.com/odb/go-dicom/network/dicomstatus"
)

// SCP - Interface to scp
type SCP interface {
	Start() error
	Stop() error
	OnAssociationRequest(f func(request network.AAssociationRQ) bool)
	OnCFindRequest(f func(request network.AAssociationRQ, findLevel string, data media.DcmObj) ([]media.DcmObj, uint16))
	OnCMoveRequest(f func(request network.AAssociationRQ, moveLevel string, data media.DcmObj) uint16)
	OnCStoreRequest(f func(request network.AAssociationRQ, data media.DcmObj) uint16)
	handleConnection(conn net.Conn)
}

type scp struct {
	Port                 int
	listener             net.Listener
	onAssociationRequest func(request network.AAssociationRQ) bool
	onCFindRequest       func(request network.AAssociationRQ, findLevel string, data media.DcmObj) ([]media.DcmObj, uint16)
	onCMoveRequest       func(request network.AAssociationRQ, moveLevel string, data media.DcmObj) uint16
	onCStoreRequest      func(request network.AAssociationRQ, data media.DcmObj) uint16
}

// NewSCP - Creates an interface to scu
func NewSCP(port int) SCP {
	media.InitDict()

	return &scp{
		Port: port,
	}
}

func (s *scp) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Println("INFO, handleConnection, new connection from: ", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *scp) Stop() error {
	return s.listener.Close()
}

func (s *scp) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	pdu := network.NewPDUService()
	pdu.SetConn(rw)

	if s.onAssociationRequest != nil {
		pdu.SetOnAssociationRequest(s.onAssociationRequest)
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

			if s.onCStoreRequest == nil {
				panic("OnCStoreRequest() not implemented")
			}

			status := s.onCStoreRequest(pdu.GetAAssociationRQ(), ddo)

			if err := dimsec.CStoreWriteRSP(pdu, dco, status); err != nil {
				log.Printf("ERROR, handleConnection, C-Store failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Store Success")
		case dicomcommand.CFindRequest:
			ddo, err := dimsec.CFindReadRQ(pdu)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Find failed to read request!")
				conn.Close()
				return
			}
			queryLevel := ddo.GetString(tags.QueryRetrieveLevel)

			results := make([]media.DcmObj, 0)
			status := dicomstatus.Success

			if s.onCFindRequest == nil {
				panic("OnCFindRequest() not implemented")
			}

			results, status = s.onCFindRequest(pdu.GetAAssociationRQ(), queryLevel, ddo)

			for _, result := range results {
				err = dimsec.CFindWriteRSP(pdu, dco, result, dicomstatus.Pending)
				if err != nil {
					log.Printf("ERROR, handleConnection, C-Find failed to write response: %s", err.Error())
					conn.Close()
					return
				}
			}

			if err := dimsec.CFindWriteRSP(pdu, dco, dco, status); err != nil {
				log.Printf("ERROR, handleConnection, C-Find failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Find Success")
		case dicomcommand.CMoveRequest:
			ddo, err := dimsec.CMoveReadRQ(pdu)
			if err != nil {
				log.Println("ERROR, handleConnection, C-Move failed to read request!")
				conn.Close()
				return
			}
			moveLevel := ddo.GetString(tags.QueryRetrieveLevel)

			if s.onCMoveRequest == nil {
				panic("OnCMoveRequest() not implemented")
			}

			status := s.onCMoveRequest(pdu.GetAAssociationRQ(), moveLevel, ddo)

			if err := dimsec.CMoveWriteRSP(pdu, dco, status, 0x00); err != nil {
				log.Printf("ERROR, handleConnection, C-Move failed to write response: %s", err.Error())
				conn.Close()
				return
			}
			log.Println("INFO, handleConnection, C-Move Success")
		case dicomcommand.CEchoRequest:
			if dimsec.CEchoReadRQ(dco) {
				if err := dimsec.CEchoWriteRSP(pdu, dco); err != nil {
					log.Println("ERROR, handleConnection, C-Echo failed to write response!")
					conn.Close()
					return
				}
				log.Println("INFO, handleConnection, C-Echo Success!")
			}
		default:
			log.Printf("ERROR, handleConnection, service not implemented: %d\n", command)
			conn.Close()
			return
		}
	}

	if err != nil {
		conn.Close()
	}
}

func (s *scp) OnAssociationRequest(f func(request network.AAssociationRQ) bool) {
	s.onAssociationRequest = f
}

func (s *scp) OnCFindRequest(f func(request network.AAssociationRQ, findLevel string, data media.DcmObj) ([]media.DcmObj, uint16)) {
	s.onCFindRequest = f
}

func (s *scp) OnCMoveRequest(f func(request network.AAssociationRQ, moveLevel string, data media.DcmObj) uint16) {
	s.onCMoveRequest = f
}

func (s *scp) OnCStoreRequest(f func(request network.AAssociationRQ, data media.DcmObj) uint16) {
	s.onCStoreRequest = f
}
