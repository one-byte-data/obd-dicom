package jpeglib

// #cgo CFLAGS: -I dcmjpeg/libijg16 -I dcmjpeg/win64
// #cgo LDFLAGS: -L dcmjpeg/win64 -lijg16
// #include "dcmjpeg/dijg16.c"
// #include "dcmjpeg/eijg16.c"
import  "C"
import (
	"errors"
	"unsafe"
)

// DIJG16decode - JPEG File to RAW
func DIJG16decode(jpegData []byte, jpegSize uint32, outputData []byte, outputSize uint32) error {
	if C.decode16((*C.uchar)(unsafe.Pointer(&jpegData[0])), C.int(jpegSize), (*C.uchar)(unsafe.Pointer(&outputData[0])), C.int(outputSize)) == 1 {
		return nil
	}
	return errors.New("ERROR, Decode16 JPEG failed")
}

// EIJG16encode - RAW File to JPEG
func EIJG16encode(rawData []uint8, width uint16, height uint16, samples uint16, outData *[]byte, outSize *int, mode int) error {
	var jpegData *C.uchar
	var jpegSize C.int
	if C.encode16((*C.ushort)(unsafe.Pointer(&rawData[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpegData, &jpegSize, C.int(mode)) == 1 {
		if jpegSize > 0 {
			*outData = C.GoBytes(unsafe.Pointer(jpegData), jpegSize)
			*outSize = int(jpegSize)
			C.free(unsafe.Pointer(jpegData))
			return nil
		}
	}
	return errors.New("ERROR, Encode16 JPEG failed")
}
