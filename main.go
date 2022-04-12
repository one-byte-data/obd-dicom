package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/services"
	"git.onebytedata.com/odb/go-dicom/tags"
)

var destination *network.Destination
var version string

func main() {
	log.Printf("Starting odb-dicom %s\n\n", version)

	media.InitDict()

	hostName := flag.String("host", "localhost", "Destination host name or IP")
	calledAE := flag.String("calledae", "DICOM_SCP", "AE of the destination")
	callingAE := flag.String("callingae", "DICOM_SCU", "AE of the client")
	port := flag.Int("port", 1040, "Port of the destination system")

	studyUID := flag.String("studyuid", "", "Study UID to be added to request")

	destinationAE := flag.String("destinationae", "", "AE of the destination for a C-Move request")

	fileName := flag.String("file", "", "DICOM file to be sent")

	cecho := flag.Bool("cecho", false, "Send C-Echo to the destination")
	cfind := flag.Bool("cfind", false, "Send C-Find request to the destination")
	cmove := flag.Bool("cmove", false, "Send C-Move request to the destination")
	cstore := flag.Bool("cstore", false, "Sends a C-Store request to the destination")

	query := flag.String("query", "", "Comma seperated query to be sent with request ex: 00080020=test")

	dump := flag.Bool("dump", false, "Dump contents of DICOM file to stdout")

	datastore := flag.String("datastore", "", "Directory to use as SCP storage")

	startSCP := flag.Bool("scp", false, "Start a SCP")

	flag.Parse()

	if *startSCP {
		if *datastore == "" {
			log.Fatalln("datastore is required for scp")
		}

		if *calledAE == "" {
			log.Fatalln("calledae is required for scp")
		}
		scp := services.NewSCP(*port)

		scp.OnAssociationRequest(func(request network.AAssociationRQ) bool {
			called := request.GetCalledAE()
			return *calledAE == called
		})

		scp.OnCFindRequest(func(request network.AAssociationRQ, queryLevel string, query media.DcmObj) []media.DcmObj {
			query.DumpTags()
			results := make([]media.DcmObj, 0)
			for i := 0; i < 10; i++ {
				results = append(results, media.GenerateCFindRequest())
			}
			return results
		})

		scp.OnCMoveRequest(func(request network.AAssociationRQ, moveLevel string, query media.DcmObj) {
			query.DumpTags()
		})

		scp.OnCStoreRequest(func(request network.AAssociationRQ, data media.DcmObj) {
			log.Printf("INFO, C-Store recieved %s", data.GetString(tags.SOPInstanceUID))
			directory := filepath.Join(*datastore, data.GetString(tags.PatientID), data.GetString(tags.StudyInstanceUID), data.GetString(tags.SeriesInstanceUID))
			os.MkdirAll(directory, 0755)

			path := filepath.Join(directory, data.GetString(tags.SOPInstanceUID)+".dcm")

			err := data.WriteToFile(path)
			if err != nil {
				log.Printf("ERROR: There was an error saving %s : %s", path, err.Error())
			}
		})

		err := scp.StartServer()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	destination = &network.Destination{
		Name:      *hostName,
		HostName:  *hostName,
		CalledAE:  *calledAE,
		CallingAE: *callingAE,
		Port:      *port,
		IsCFind:   true,
		IsCStore:  true,
		IsMWL:     true,
		IsTLS:     false,
	}

	if *cecho {
		scu := services.NewSCU(destination)
		err := scu.EchoSCU(30)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CEcho was successful")
	}
	if *cfind {
		request := media.DefaultCFindRequest()
		scu := services.NewSCU(destination)
		scu.SetOnCFindResult(func(result media.DcmObj) {
			log.Printf("Found study %s\n", result.GetString(tags.StudyInstanceUID))
			result.DumpTags()
		})

		if *query != "" {
			parts := strings.Split(*query, ",")
			for _, part := range parts {
				log.Println(part)
				// p := strings.Split(part, "=")
				// tag := media.DcmTag{

				// }
			}
		}

		count, status, err := scu.FindSCU(request, 0)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("CFind was successful")
		log.Printf("Found %d results with status %d\n\n", count, status)
		os.Exit(0)
	}
	if *cmove {
		if *destinationAE == "" {
			log.Fatalln("destinationae is required for a C-Move")
		}
		if *studyUID == "" {
			log.Fatalln("studyuid is required for a C-Move")
		}

		request := media.DefaultCMoveRequest(*studyUID)

		scu := services.NewSCU(destination)
		_, err := scu.MoveSCU(*destinationAE, request, 0)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CMove was successful")
		os.Exit(0)
	}
	if *cstore {
		if *fileName == "" {
			log.Fatalln("file is required for a C-Store")
		}
		scu := services.NewSCU(destination)
		err := scu.StoreSCU(*fileName, 0)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("CStore of %s was successful", *fileName)
		os.Exit(0)
	}
	if *dump {
		if *fileName == "" {
			log.Fatalln("file is required for a dump")
		}
		obj, err := media.NewDCMObjFromFile(*fileName)
		if err != nil {
			log.Panicln(err)
		}
		obj.DumpTags()
		os.Exit(0)
	}
}
