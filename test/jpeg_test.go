package main

import (
	"fmt"
	"os"
	"testing"

	"git.onebytedata.com/odb/go-dicom/jpeglib"
)

func Test_JPEGLibEIJG8decode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should decode j2k image",
			args:    args{fileName: "./images/test.jpg"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jpegData []byte
			var outData []byte

			if LoadFromFile(tt.args.fileName, &jpegData) {
				outSize := 1576 * 1134 * 3 // Image Size, have to know in advance.
				outData = make([]byte, outSize)

				if err := jpeglib.DIJG8decode(jpegData, uint32(len(jpegData)), outData, uint32(outSize)); (err != nil) != tt.wantErr {
					t.Errorf("jpeglib.DIJG8decode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func Test_JPEGLibEIJG8encode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should encode jpeg 8 image",
			args:    args{fileName: "./images/test.raw"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jpegData []byte
			outData := make([]byte, 0)
			var jpegSize int

			if LoadFromFile(tt.args.fileName, &outData) {
				if err := jpeglib.EIJG8encode(outData, 1576, 1134, 3, &jpegData, &jpegSize, 4); (err != nil) != tt.wantErr {
					t.Errorf("jpeglib.EIJG8encode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

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
