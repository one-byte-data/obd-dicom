package media

import (
	"testing"
)

func TestNewDCMObjFromFile(t *testing.T) {
	InitDict()

	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should load DICOM file",
			args:    args{fileName: "../test/images/test.dcm"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dcmObj, err := NewDCMObjFromFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDCMObjFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			dcmObj.DumpTags()
		})
	}
}
