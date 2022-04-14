package jpeglib

// import (
// 	"testing"
// )

// func TestDIJG12decode(t *testing.T) {
// 	type args struct {
// 		fileName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Should decode jpeg 12 image",
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

// 				if err := DIJG12decode(jpegData, uint32(len(jpegData)), outData, uint32(outSize)); (err != nil) != tt.wantErr {
// 					t.Errorf("DIJG12decode() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			}
// 		})
// 	}
// }

// func Test_DIJG12encode(t *testing.T) {
// 	type args struct {
// 		fileName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Should encode jpeg 12 image",
// 			args:    args{fileName: "../samples/test.raw"},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var jpegData []byte
// 			outData := make([]byte, 0)
// 			var jpegSize int

// 			if LoadFromFile(tt.args.fileName, &outData) {
// 				if err := EIJG12encode(outData, 1576, 1134, 3, &jpegData, &jpegSize, 4); (err != nil) != tt.wantErr {
// 					t.Errorf("EIJG12encode() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			}
// 		})
// 	}
// }
