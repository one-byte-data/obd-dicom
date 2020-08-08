package jpeglib

// #cgo CFLAGS: -Idcmjpeg/libijg16 -Idcmjpeg/win64
// #cgo LDFLAGS: -Ldcmjpeg/win64 -lijg16
// #include "dcmjpeg/dijg16.c"
// #include "dcmjpeg/eijg16.c"
import "C"
import (
	"errors"
	"unsafe"
)

// DIJG16decode - JPEG File to RAW
func DIJG16decode(jpegData []byte, jpegSize int, outputData []byte, outputSize int) error {
	if C.decode16((*C.uchar)(unsafe.Pointer(&jpegData[0])), C.int(jpegSize), (*C.uchar)(unsafe.Pointer(&outputData[0])), C.int(outputSize)) == 1 {
		return nil
	}
	return errors.New("ERROR, Decode16 JPEG failed")
}

// EIJG16encode - RAW File to JPEG
func EIJG16encode(rawData []uint8, width int, height int, samples int, outData *[]byte) error {
	var jpegData *C.uchar
	var jpegSize C.int
	if C.encode16((*C.ushort)(unsafe.Pointer(&rawData[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpegData, &jpegSize, C.int(0)) == 1 {
		if jpegSize > 0 {
			*outData = C.GoBytes(unsafe.Pointer(jpegData), jpegSize)
			C.free(unsafe.Pointer(jpegData))
			return nil
		}
	}
	return errors.New("ERROR, Encode16 JPEG failed")
}
