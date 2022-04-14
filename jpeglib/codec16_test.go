package jpeglib

import (
	"testing"
)

// func TestDIJG16decode(t *testing.T) {
// 	type args struct {
// 		fileName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Should decode jpeg 16 image",
// 			args:    args{fileName: "../samples/test.jpg"},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var jpegData []byte
// 			var outData []byte

// 			if LoadFromFile(tt.args.fileName, &jpegData) {
// 				outSize := 1576 * 1134 * 3 // Image Size, have to know in advance.
// 				outData = make([]byte, outSize)

// 				if err := DIJG16decode(jpegData, uint32(len(jpegData)), outData, uint32(outSize)); (err != nil) != tt.wantErr {
// 					t.Errorf("DIJG16decode() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			}
// 		})
// 	}
// }

func Test_DIJG16encode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should encode jpeg 16 image",
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
				if err := EIJG16encode(outData, 1576, 1134, 3, &jpegData, &jpegSize, 4); (err != nil) != tt.wantErr {
					t.Errorf("EIJG16encode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
