package sopclass

import (
	"reflect"
	"testing"
)

func TestGetSOPClassFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *SOPClass
	}{
		{
			name: "Should get Verification SOP class",
			args: args{name: "Verification"},
			want: &SOPClass{
				UID:         "1.2.840.10008.1.1",
				Name:        "Verification",
				Description: "Verification SOP Class",
				Type:        "SOP Class",
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
			if got := GetSOPClassFromName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOPClassFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSOPClassFromUID(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
		want *SOPClass
	}{
		{
			name: "Should get Verification SOP class",
			args: args{uid: "1.2.840.10008.1.1"},
			want: &SOPClass{
				UID:         "1.2.840.10008.1.1",
				Name:        "Verification",
				Description: "Verification SOP Class",
				Type:        "SOP Class",
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
			if got := GetSOPClassFromUID(tt.args.uid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOPClassFromUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
