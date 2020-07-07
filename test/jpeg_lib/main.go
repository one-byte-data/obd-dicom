package main

// #cgo CFLAGS: -Idcmjpeg/libijg8
// #cgo LDFLAGS: -Ldcmjpeg/libijg8 -lijg8
// #include "dcmjpeg/dijg8.c"
// #include "dcmjpeg/eijg8.c"
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// LoadFromFile - Load from File into MemoryStream
func LoadFromFile(FileName string, buffer *[]byte) bool {
	flag := false

	file, err := os.Open(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file: " + FileName)
		return flag
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("ERROR, getting file Stats: " + FileName)
		return flag
	}
	size := int(stat.Size())
	bs := make([]byte, size)
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("ERROR, reading file: " + FileName)
		return flag
	}
	*buffer = append(*buffer, bs...)
	return true
}

// SaveToFile - Save MemoryStream to File
func SaveToFile(FileName string, buffer []byte) bool {
	flag := false

	file, err := os.Create(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file: " + FileName)
		return flag
	}
	defer file.Close()
	_, err = file.Write(buffer)
	if err != nil {
		fmt.Println("ERROR, writing to file: " + FileName)
		return flag
	}
	return true
}

// Decode JPEG File to RAW
func Decode(filename string, width int, height int, samples int) bool {
	flag := false
	var jpegData []byte
	if LoadFromFile(filename, &jpegData) {
		outputSize := width * height * samples // Image Size, have to know in advance.
		outputData := make([]byte, outputSize)
		if C.decode8((*C.uchar)(unsafe.Pointer(&jpegData[0])), C.int(len(jpegData)), (*C.uchar)(unsafe.Pointer(&outputData[0])), C.int(outputSize)) == 1 {
			if SaveToFile("out.raw", outputData) {
				fmt.Println("INFO, saved raw data")
				flag = true
			}
		} else {
			fmt.Println("ERROR, Decode JPEG failed!!")
		}
	}
	return flag
}

// Encode RAW File to JPEG
func Encode(filename string, width int, height int, samples int) bool {
	flag := false
	var rawData []byte
	if LoadFromFile(filename, &rawData) {
		var jpegData *C.uchar
		var jpegSize C.int
		if C.encode8((*C.uchar)(unsafe.Pointer(&rawData[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpegData, &jpegSize, C.int(0)) == 1 {
			if jpegSize > 0 {
				imgdata := C.GoBytes(unsafe.Pointer(jpegData), jpegSize)
				if SaveToFile("out.jpg", imgdata) {
					fmt.Println("INFO, saved jpeg data")
					flag = true
				}
				C.free(unsafe.Pointer(jpegData))
			}
		}
	}
	return flag
}

func main() {
	if Decode("test.jpg", 1576, 1134, 3) {
		fmt.Println("Decode Success!")
	} else {
		fmt.Println("Decode Failed!")
	}
	if Encode("test.raw", 1576, 1134, 3) {
		fmt.Println("Encode Success!")
	} else {
		fmt.Println("Encode Failed!")
	}
}
