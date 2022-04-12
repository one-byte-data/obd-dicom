package main

import (
	"testing"

	"git.onebytedata.com/odb/go-dicom/openjpeg"
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
			args:    args{fileName: "./images/test.j2k"},
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

				if err := openjpeg.J2Kdecode(jpegData, uint32(len(jpegData)), outData); (err != nil) != tt.wantErr {
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
				if err := openjpeg.J2Kencode(outData, 1576, 1134, 3, 8, &jpegData, &jpegSize, 10); (err != nil) != tt.wantErr {
					t.Errorf("openjpeg.J2Kencode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
