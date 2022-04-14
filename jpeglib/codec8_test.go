package jpeglib

import (
	"os"
	"testing"
)

func TestDIJG8decode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should decode jpeg 8 image",
			args:    args{fileName: "../samples/test8.jpg"},
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

				if err := DIJG8decode(jpegData, uint32(len(jpegData)), outData, uint32(outSize)); (err != nil) != tt.wantErr {
					t.Errorf("DIJG8decode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func Test_DIJG8encode(t *testing.T) {
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
			args:    args{fileName: "../samples/test.raw"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jpegData []byte
			outData := make([]byte, 0)
			var jpegSize int

			if LoadFromFile(tt.args.fileName, &outData) {
				if err := EIJG8encode(outData, 1576, 1134, 3, &jpegData, &jpegSize, 4); (err != nil) != tt.wantErr {
					t.Errorf("EIJG8encode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func LoadFromFile(FileName string, buffer *[]byte) bool {
	file, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	size := int(stat.Size())
	bs := make([]byte, size)
	_, err = file.Read(bs)
	if err != nil {
		panic(err)
	}
	*buffer = append(*buffer, bs...)
	return true
}
