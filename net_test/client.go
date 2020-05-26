package main

import (
	"fmt"
	"net"
	"rafael/DICOM/network"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1040")
	if err != nil {
		fmt.Println(err)
		return
	}

	network.Resetuniq()
	aarq := network.NewAAAssociationRQ()
	copy(aarq.CallingApTitle[:], "TESTSCU")
	copy(aarq.CalledApTitle[:], "CHARRUAPACS")

	PresContext := network.NewPresentationContext()
	PresContext.SetAbstractSyntax("1.2.840.10008.1.1")
	PresContext.AddTransferSyntax("1.2.840.10008.1.2")
	aarq.PresContexts = append(aarq.PresContexts, *PresContext)

	aarq.UserInfo.MaxSubLength.MaximumLength=16384
	aarq.UserInfo.SetImpClassUID("1.2.826.0.1.3680043.2.1396.999")
	aarq.UserInfo.SetImpVersionName("CharruaSoft")
	if aarq.Write(conn) {
		aaac := network.NewAAssociationAC()
		aaac.Read(conn)
	}
}
