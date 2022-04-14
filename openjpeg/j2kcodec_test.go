package openjpeg

import (
	"os"
	"testing"
)

func Test_J2Kdecode(t *testing.T) {
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
			args:    args{fileName: "../samples/test.j2k"},
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

				if err := J2Kdecode(jpegData, uint32(len(jpegData)), outData); (err != nil) != tt.wantErr {
					t.Errorf("openjpeg.J2Kdecode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func Test_J2Kencode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should encode j2k image",
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
				if err := J2Kencode(outData, 1576, 1134, 3, 8, &jpegData, &jpegSize, 10); (err != nil) != tt.wantErr {
					t.Errorf("openjpeg.J2Kencode() error = %v, wantErr %v", err, tt.wantErr)
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
