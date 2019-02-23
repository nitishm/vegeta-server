package vegeta

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_addID(t *testing.T) {
	type args struct {
		report *bytes.Buffer
		id     string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"add id short", args{bytes.NewBufferString("test string"), "123"}, []byte("ID 123\ntest string")},
		{"add id long", args{bytes.NewBufferString("test string"), "283orcniouq3hnq8hcqn3f8ohuicfbhn"},
			[]byte("ID 283orcniouq3hnq8hcqn3f8ohuicfbhn\ntest string")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addID(tt.args.report, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addID() = %v, want %v", got, tt.want)
			}
		})
	}
}
