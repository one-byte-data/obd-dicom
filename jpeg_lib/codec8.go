package main

// #cgo CFLAGS: -Idcmjpeg/libijg8 -Idcmjpeg/win64
// #cgo LDFLAGS: -Ldcmjpeg/win64 -lijg8
// #include "dcmjpeg/dijg8.c"
// #include "dcmjpeg/eijg8.c"
import "C"
import (
	"log"
	"unsafe"
)

// Decode JPEG File to RAW
func Decode8(jpeg_data []byte, jpeg_size int, output_data []byte, output_size int) bool {
	flag := false
	if C.decode8((*C.uchar)(unsafe.Pointer(&jpeg_data[0])), C.int(jpeg_size), (*C.uchar)(unsafe.Pointer(&output_data[0])), C.int(output_size)) == 1 {
		flag = true
	} else {
		log.Println("ERROR, Decode8, JPEG failed!!")
	}
	return flag
}

// Encode RAW File to JPEG
func Encode8(raw_data []byte, width int, height int, samples int, out_data *[]byte) bool {
	flag := false
	var jpeg_data *C.uchar
	var jpegSize C.int
	if C.encode8((*C.uchar)(unsafe.Pointer(&raw_data[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpeg_data, &jpegSize, C.int(0)) == 1 {
		if jpegSize > 0 {
			*out_data = C.GoBytes(unsafe.Pointer(jpeg_data), jpegSize)
			C.free(unsafe.Pointer(jpeg_data))
			flag=true
		}
	}
	return flag
}
