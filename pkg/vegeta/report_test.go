package vegeta

import (
	"bytes"
	"github.com/pkg/errors"
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

func TestFormat(t *testing.T) {
	type args struct {
		f   Format
		sli []string
	}
	type want struct {
		f   Format
		ff  Format
		fs  []Format
		buk []byte
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Formats only Check",
			args: args{
				f:   Format("json"),
				sli: []string{"json", "text"},
			},
			want: want{
				f:   Format("json"),
				ff:  Format("json-text"),
				fs:  []Format{Format("json"), Format("text")},
				buk: []byte{},
				err: errors.New("time bucket can be obtained by histogram format only"),
			},
		},
		{
			name: "Histogram Check",
			args: args{
				f:   Format("histogram"),
				sli: []string{"histogram", "0,1s,2s,3s"},
			},
			want: want{
				f:   Format("histogram"),
				ff:  Format("histogram-0,1s,2s,3s"),
				fs:  []Format{Format("histogram"), Format("0,1s,2s,3s")},
				buk: []byte("[0,1s,2s,3s]"),
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if f := tt.args.f.GetFormat(); f != tt.want.f {
				t.Errorf("GetFormat() = %v, want %v", f, tt.want.f)
			}
			ff := StringsToFormat(tt.args.sli...)
			if ff != tt.want.ff {
				t.Errorf("StringsToFormat() = %v, want %v", ff, tt.want.ff)
			}
			if g := ff.SplitFormat(); !reflect.DeepEqual(g, tt.want.fs) {
				t.Errorf("SplitFormat() = %v, want %v", g, tt.want.fs)
			}
			if b, err := ff.GetTimeBucketOfHistogram(); !reflect.DeepEqual(string(b), string(tt.want.buk)) && !reflect.DeepEqual(err.Error(), tt.want.err.Error()) {
				t.Errorf("GetTimeBucketOfHistogram() Bucket = %v, want %v", b, tt.want.buk)
				t.Errorf("GetTimeBucketOfHistogram() Error = %v, want %v", err, tt.want.err)
			}
		})
	}
}
