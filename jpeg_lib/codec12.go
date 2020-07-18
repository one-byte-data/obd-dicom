package main

// #cgo CFLAGS: -Idcmjpeg/libijg12 -Idcmjpeg/win64
// #cgo LDFLAGS: -Ldcmjpeg/win64 -lijg12
// #include "dcmjpeg/dijg12.c"
// #include "dcmjpeg/eijg12.c"
import "C"
import (
	"errors"
	"unsafe"
)

// Decode12 - JPEG File to RAW
func Decode12(jpegData []byte, jpegSize int, outputData []byte, outputSize int) error {
	if C.decode12((*C.uchar)(unsafe.Pointer(&jpegData[0])), C.int(jpegSize), (*C.uchar)(unsafe.Pointer(&outputData[0])), C.int(outputSize)) == 1 {
		return nil
	}
	return errors.New("ERROR, Decode12 JPEG failed")
}

// Encode12 - RAW File to JPEG
func Encode12(rawData []uint8, width int, height int, samples int, outData *[]byte) error {
	var jpegData *C.uchar
	var jpegSize C.int
	if C.encode12((*C.ushort)(unsafe.Pointer(&rawData[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpegData, &jpegSize, C.int(0)) == 1 {
		if jpegSize > 0 {
			*outData = C.GoBytes(unsafe.Pointer(jpegData), jpegSize)
			C.free(unsafe.Pointer(jpegData))
			return nil
		}
	}
	return errors.New("ERROR, Encode12 JPEG failed")
}
