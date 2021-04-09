package main

import (
	"log"
	"os"

	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/services"
	"git.onebytedata.com/odb/go-dicom/tags"
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
		if queryLevel == "STUDY" {
			var study DCMStudy
			err, results := study.Select(obj)
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			return results
		}
		if queryLevel == "SERIES" {
			var series DCMSeries
			err, results := series.Select(obj)
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			return results
		}
		return nil
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
