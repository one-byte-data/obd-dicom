package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/test/HL7Broker/hl7"
)

const (
	startBlock = '\x0b'
	endBlock   = '\x1c'
	cr         = '\x0d'
)

func checkByte(msg []byte, pos int, expected byte) error {
	if msg[pos] != expected {
		return fmt.Errorf("invalid message %v, expected %v at position %v but got %v",
			msg, expected, pos, msg[pos])
	}
	return nil
}

func decapsulate(msg []byte) ([]byte, error) {
	if len(msg) < 3 {
		return nil, fmt.Errorf("short message, length %v", len(msg))
	}
	if err := checkByte(msg, 0, startBlock); err != nil {
		return nil, err
	}
	if err := checkByte(msg, len(msg)-2, endBlock); err != nil {
		return nil, err
	}
	if err := checkByte(msg, len(msg)-1, cr); err != nil {
		return nil, err
	}
	return msg[1 : len(msg)-2], nil
}

func handleConnection(c net.Conn) {
	flag := false
	fmt.Print(".")
	reader := bufio.NewReader(c)
	rawMsg, err := reader.ReadBytes(endBlock)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read one more byte for the carriage return.
	lastByte, err := reader.ReadByte()
	if err != nil {
		return
	}
	msg, err := decapsulate(append(rawMsg, lastByte))
	for _, line := range strings.Split(string(msg), "\r") {
		//	fmt.Println(line)
		if strings.Contains(line, "MSH|") {
			if !strings.Contains(line, "ORU^R01") {
				break
			}
		}
		if strings.Contains(line, "PID") {
			hl7.ParsePID(line)
		}
		if strings.Contains(line, "PV1") {
			hl7.ParsePV1(line)
		}
		if strings.Contains(line, "ORC") {
			hl7.ParseORC(line)
		}
		if strings.Contains(line, "OBR") {
			hl7.ParseOBR(line)
		}
		if strings.Contains(line, "OBX") {
			hl7.ParseOBX(line)
			flag = true
		}
	}
	if err != nil {
		return
	}
	if flag == true {
		hl7.SaveDICOMSR("test.dcm")
		c.Write([]byte("Ack"))
	} else {
		c.Write([]byte("Nack"))
	}
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
