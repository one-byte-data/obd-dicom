package main

// #cgo CFLAGS: -Idcmjpeg/libijg16 -Idcmjpeg/win64
// #cgo LDFLAGS: -Ldcmjpeg/win64 -lijg16
// #include "dcmjpeg/dijg16.c"
// #include "dcmjpeg/eijg16.c"
import "C"
import (
	"log"
	"unsafe"
)

// Decode JPEG File to RAW
func Decode16(jpeg_data []byte, jpeg_size int, output_data []byte, output_size int) bool {
	flag := false
	if C.decode16((*C.uchar)(unsafe.Pointer(&jpeg_data[0])), C.int(jpeg_size), (*C.uchar)(unsafe.Pointer(&output_data[0])), C.int(output_size)) == 1 {
		flag = true
	} else {
		log.Println("ERROR, Decode16 JPEG failed!!")
	}
	return flag
}

// Encode RAW File to JPEG
func Encode16(raw_data []uint8, width int, height int, samples int, out_data *[]byte) bool {
	flag := false
	var jpeg_data *C.uchar
	var jpegSize C.int
	if C.encode16((*C.ushort)(unsafe.Pointer(&raw_data[0])), C.ushort(width), C.ushort(height), C.ushort(samples), &jpeg_data, &jpegSize, C.int(0)) == 1 {
		if jpegSize > 0 {
			*out_data = C.GoBytes(unsafe.Pointer(jpeg_data), jpegSize)
			C.free(unsafe.Pointer(jpeg_data))
			flag=true
		}
	}
	return flag
}
