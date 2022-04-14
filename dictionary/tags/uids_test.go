package tags

import (
	"reflect"
	"testing"
)

func TestGetTagFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Tag
	}{
		{
			name: "Should get CommandGroupLength tag",
			args: args{name: "CommandGroupLength"},
			want: &Tag{
				Group:       0x0000,
				Element:     0x0000,
				VR:          "UL",
				VM:          "1",
				Name:        "CommandGroupLength",
				Description: "Command Group Length",
			},
		},
		{
			name: "Should get empty tag from invlid name",
			args: args{name: "Not valid"},
			want: &Tag{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTagFromName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTag(t *testing.T) {
	type args struct {
		group   uint16
		element uint16
	}
	tests := []struct {
		name string
		args args
		want *Tag
	}{
		{
			name: "Should get CommandGroupLength tag",
			args: args{group: 0x0000, element: 0x0000},
			want: &Tag{
				Group:       0x0000,
				Element:     0x0000,
				VR:          "UL",
				VM:          "1",
				Name:        "CommandGroupLength",
				Description: "Command Group Length",
			},
		},
		{
			name: "Should get empty tag from invalid group and element",
			args: args{group: 0xffff, element: 0xffff},
			want: &Tag{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTag(tt.args.group, tt.args.element); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagFromUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
