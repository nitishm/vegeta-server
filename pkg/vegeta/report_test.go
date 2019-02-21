package vegeta

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_addId(t *testing.T) {
	type args struct {
		report *bytes.Buffer
		id     string
	}
	tests := []struct {
		name string
		args args
		want *bytes.Buffer
	}{
		{"add id short", args{bytes.NewBufferString("test string"), "123"}, bytes.NewBufferString("ID 123\ntest string")},
		{"add id long", args{bytes.NewBufferString("test string"), "283orcniouq3hnq8hcqn3f8ohuicfbhn"},
			bytes.NewBufferString("ID 283orcniouq3hnq8hcqn3f8ohuicfbhn\ntest string")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addId(tt.args.report, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addId() = %v, want %v", got, tt.want)
			}
		})
	}
}
