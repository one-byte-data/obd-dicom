package transfersyntax

import (
	"reflect"
	"testing"
)

func TestGetTransferSyntaxFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *TransferSyntax
	}{
		{
			name: "Should get ILE transfer syntax",
			args: args{name: "ImplicitVRLittleEndian"},
			want: &TransferSyntax{
				UID:         "1.2.840.10008.1.2",
				Name:        "ImplicitVRLittleEndian",
				Description: "Implicit VR Little Endian",
				Type:        "Transfer Syntax",
			},
		},
		{
			name: "Should get nil from invalid transfer syntax name",
			args: args{name: "DERP"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTransferSyntaxFromName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransferSyntaxFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransferSyntaxFromUID(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
		want *TransferSyntax
	}{
		{
			name: "Should get ILE transfer syntax",
			args: args{uid: "1.2.840.10008.1.2"},
			want: &TransferSyntax{
				UID:         "1.2.840.10008.1.2",
				Name:        "ImplicitVRLittleEndian",
				Description: "Implicit VR Little Endian",
				Type:        "Transfer Syntax",
			},
		},
		{
			name: "Should get nil from invalid transfer syntax UID",
			args: args{uid: "1.2.840.10008.1.2.00000"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTransferSyntaxFromUID(tt.args.uid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransferSyntaxFromUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
