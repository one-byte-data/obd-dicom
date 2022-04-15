package main

import (
	"flag"
	"log"
	"reflect"

	"git.onebytedata.com/odb/go-dicom/media"
)

var version string

func main() {
	log.Printf("Starting compare %s\n\n", version)

	media.InitDict()

	sourceFile := flag.String("s", "", "Source DICOM file")
	destinationFile := flag.String("d", "", "Destination DICOM file to compare against")

	flag.Parse()

	if *sourceFile == "" || *destinationFile == "" {
		log.Fatalln("Both a source and destination is required")
	}
	srcDicom, err := media.NewDCMObjFromFile(*sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dumping source tags")
	srcDicom.DumpTags()

	dstDicom, err := media.NewDCMObjFromFile(*destinationFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dumping destination tags")
	dstDicom.DumpTags()

	compare(srcDicom, dstDicom)
}

func compare(source media.DcmObj, destination media.DcmObj) {
	for _, st := range source.GetTags() {
		found := false
		for _, dt := range destination.GetTags() {
			if st.Group == dt.Group && st.Element == dt.Element {
				found = true
				if !reflect.DeepEqual(st.Data, dt.Data) {
					if len(st.Data) > 128 || len(dt.Data) > 128 {
						log.Printf("Tag: %s are not equal", st.Name)
					} else {
						log.Printf("Tag: %s are not equal, source: %s, destination: %s", st.Name, st.Data, dt.Data)
					}
				}
			}
		}
		if !found {
			log.Printf("Tag: %s not found in destination", st.Name)
		}
	}
}
