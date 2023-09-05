package storage

import (
	"testing"
)

func TestCalculateHash(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"empty", args{""}, 0, true},
		{"simple word", args{"hello"}, 0, false},
		{"only letters", args{"asdasdasdas"}, 1, false},
		{"combined", args{"12341234#$@!#asasdads"}, 2, false},
		{"long and combined", args{"!@$#!@#!@#122121asdadsadsadsdadsads"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pickNodeNumber(tt.args.data, 3)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("pickNodeNumber() = %v, want %v", got, tt.want)
			}

		})
	}
}
