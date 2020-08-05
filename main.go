package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/network"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/services"
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

		scp.SetOnAssociationRequest(func(request network.AAssociationRQ) bool {
			called := request.GetCalledAE()

			log.Printf("%x", []byte(called))

			log.Printf("%x", []byte(*calledAE))

			return *calledAE == called
		})

		scp.SetOnCFindRequest(func(request network.AAssociationRQ, queryLevel string, query media.DcmObj, result media.DcmObj) {
			query.DumpTags()
		})

		scp.SetOnCMoveRequest(func(request network.AAssociationRQ, moveLevel string, query media.DcmObj) {
			query.DumpTags()
		})

		scp.SetOnCStoreRequest(func(request network.AAssociationRQ, data media.DcmObj) {
			log.Printf("INFO, C-Store recieved %s", data.GetString(0x0008, 0x0018))
			directory := filepath.Join(*datastore, data.GetString(0x0010, 0x0020), data.GetString(0x0020, 0x000d), data.GetString(0x0020, 0x000e))
			os.MkdirAll(directory, 0755)

			path := filepath.Join(directory, data.GetString(0x0008, 0x0018)+".dcm")

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

		results := make([]media.DcmObj, 0)
		_, err := scu.FindSCU(request, &results, 30)
		if err != nil {
			log.Fatalln(err)
		}

		if *query != "" {
			parts := strings.Split(*query, ",")
			for _, part := range parts {
				log.Println(part)
				// p := strings.Split(part, "=")
				// tag := media.DcmTag{

				// }
			}
		}

		log.Println("CFind was successful")
		log.Printf("Found %d results\n\n", len(results))
		for _, result := range results {
			log.Printf("Found study %s\n", result.GetString(0x0020, 0x000D))
			result.DumpTags()
		}
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
		_, err := scu.MoveSCU(*destinationAE, request, 30)
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
		err := scu.StoreSCU(*fileName, 30)
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