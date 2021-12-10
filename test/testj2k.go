package main

import (
	"fmt"
	"log"

	"git.onebytedata.com/odb/go-dicom/openjpeg"
)

func testj2k() {
	var jpegData []byte
	var outData []byte

	if LoadFromFile("test.j2k", &jpegData) {
		outSize := 1576 * 1134 * 3 // Image Size, have to know in advance.
		outData = make([]byte, outSize)

		err := openjpeg.J2Kdecode(jpegData, uint32(len(jpegData)), outData)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Decode Success!")
		if SaveToFile("out.raw", outData) {
			fmt.Println("Saved out.raw")
		} else {
			fmt.Println("ERROR, Saving out.raw")
		}
	} else {
		fmt.Println("ERROR, Decode Failed!")
	}

	jpegData = nil
	var jpegSize int
	outData = nil
	if LoadFromFile("test.raw", &outData) {
		err := openjpeg.J2Kencode(outData, 1576, 1134, 3, 8, &jpegData, &jpegSize, 10)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Encode Success!")
		if SaveToFile("out.j2k", jpegData) {
			fmt.Println("Saved out.j2k")
		} else {
			fmt.Println("ERROR, Saving out.j2k")
		}
	} else {
		fmt.Println("ERROR, Encode Failed!")
	}
}

/*
func main(){
	testj2k()
}
*/
