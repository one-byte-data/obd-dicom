package network

import (
	"reflect"
	"testing"
)

func TestNewUserInformation(t *testing.T) {
	tests := []struct {
		name string
		want UserInformation
	}{
		{
			name: "Should get UserInformation",
			want: &userInformation{
				ItemType:      0x50,
				MaxSubLength:  NewMaximumSubLength(),
				AsyncOpWindow: NewAsyncOperationWindow(),
				SCPSCURole:    NewRoleSelect(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserInformation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}
