package network

import (
	"reflect"
	"testing"

	"git.onebytedata.com/odb/go-dicom/dictionary/sopclass"
)

func TestNewUIDitem(t *testing.T) {
	type args struct {
		uid      string
		itemType byte
	}
	tests := []struct {
		name string
		args args
		want UIDItem
	}{
		{
			name: "Should create UIDItem",
			args: args{
				uid:      sopclass.MediaStorageDirectoryStorage.UID,
				itemType: 0x00,
			},
			want: &uidItem{
				itemType:  0x00,
				reserved1: 0x00,
				length:    uint16(len(sopclass.MediaStorageDirectoryStorage.UID)),
				uid:       sopclass.MediaStorageDirectoryStorage.UID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUIDItem(tt.args.uid, tt.args.itemType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUIDitem() = %v, want %v", got, tt.want)
			}
		})
	}
}
