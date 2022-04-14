package codingscheme

import (
	"reflect"
	"testing"
)

func TestGetCodingSchemeFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *CodingScheme
	}{
		{
			name: "Should get DCM scheme",
			args: args{name: "DCM"},
			want: &CodingScheme{
				UID:         "1.2.840.10008.2.16.4",
				Name:        "DCM",
				Description: "DICOM Controlled Terminology",
				Type:        "Coding Scheme",
			},
		},
		{
			name: "Should get nil from invlid name",
			args: args{name: "Not valid"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCodingSchemeFromName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCodingSchemeFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCodingSchemeFromUID(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
		want *CodingScheme
	}{
		{
			name: "Should get DCM scheme",
			args: args{uid: "1.2.840.10008.2.16.4"},
			want: &CodingScheme{
				UID:         "1.2.840.10008.2.16.4",
				Name:        "DCM",
				Description: "DICOM Controlled Terminology",
				Type:        "Coding Scheme",
			},
		},
		{
			name: "Should get nil from invalid UID",
			args: args{uid: "1.2.84.1.1"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCodingSchemeFromUID(tt.args.uid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCodingSchemeFromUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
