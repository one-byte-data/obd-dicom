package main

import (
	"log"
	"os"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/services"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/tags"
)

var destination *network.Destination

func main() {
	media.InitDict()
	calledAE := "DICOM_SCP"
	port := 1040
	scp := services.NewSCP(port)

	scp.SetOnAssociationRequest(func(request network.AAssociationRQ) bool {
		called := request.GetCalledAE()
		return calledAE == called
		})

	scp.SetOnCFindRequest(func(request network.AAssociationRQ, queryLevel string, obj media.DcmObj) []media.DcmObj {
		// Make Query from query.
		var study DCMStudy 
		QueryString := study.DICOM2Query(obj)
		err, results:= study.QueryDB(QueryString)
		if err!= nil {
			log.Println(err.Error())
			return nil
		}
		return results
	})

	scp.SetOnCMoveRequest(func(request network.AAssociationRQ, moveLevel string, query media.DcmObj) {
		query.DumpTags()
	})

	scp.SetOnCStoreRequest(func(request network.AAssociationRQ, data media.DcmObj) {
		log.Printf("INFO, C-Store received %s", data.GetString(tags.SOPInstanceUID))
	})

	err := scp.StartServer()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
