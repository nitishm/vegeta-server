package vegeta

import (
	"bytes"
	"fmt"
	"io"
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

func TestCreateReportFromReader(t *testing.T) {
	type args struct {
		reader io.Reader
		id     string
		format Format
	}
	type want struct {
		byt []byte
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Unsupported format",
			args: args{
				reader: bytes.NewBuffer([]byte{}),
				id:     "id",
				format: &HistogramFormat{repr: "unsupported format", meta: make(MetaInfo)},
			},
			want: want{
				byt: nil,
				err: fmt.Errorf("format unsupported format not supported"),
			},
		},
		{
			name: "JSONFormat",
			args: args{
				reader: bytes.NewBuffer([]byte{
					119, 255, 129, 3, 1, 1, 6, 82, 101, 115, 117, 108, 116, 1, 255, 130, 0, 1, 9, 1, 6, 65, 116, 116, 97, 99, 107, 1, 12, 0, 1, 3, 83, 101, 113, 1, 6, 0, 1, 4, 67, 111, 100, 101, 1, 6, 0, 1, 9, 84, 105, 109, 101, 115, 116, 97, 109, 112, 1, 255, 132, 0, 1, 7, 76, 97, 116, 101, 110, 99, 121, 1, 4, 0, 1, 8, 66, 121, 116, 101, 115, 79, 117, 116, 1, 6, 0, 1, 7, 66, 121, 116, 101, 115, 73, 110, 1, 6, 0, 1, 5, 69, 114, 114, 111, 114, 1, 12, 0, 1, 4, 66, 111, 100, 121, 1, 10, 0, 0, 0, 16, 255, 131, 5, 1, 1, 4, 84, 105, 109, 101, 1, 255, 132, 0, 0, 0, 68, 255, 130, 1, 36, 56, 56, 97, 100, 97, 52, 100, 100, 45, 101, 51, 97, 98, 45, 52, 50, 49, 56, 45, 98, 101, 102, 100, 45, 98, 98, 100, 48, 97, 52, 98, 50, 99, 101, 99, 102, 2, 255, 200, 1, 15, 1, 0, 0, 0, 14, 212, 12, 177, 127, 57, 196, 111, 216, 1, 74, 1, 251, 3, 5, 4, 97, 190, 0, // nolint: lll
				}),
				id:     "id",
				format: NewJSONFormat(),
			},
			want: want{
				byt: []byte(`{"id":"id","latencies":{"total":6484537567,"mean":6484537567,"max":6484537567,"50th":6484537567,"95th":6484537567,"99th":6484537567},"bytes_in":{"total":0,"mean":0},"bytes_out":{"total":0,"mean":0},"earliest":"2019-03-02T22:46:47.969175+05:30","latest":"2019-03-02T22:46:47.969175+05:30","end":"2019-03-02T22:46:54.453712567+05:30","duration":0,"wait":6484537567,"requests":1,"rate":1,"success":1,"status_codes":{"200":1},"errors":[]}`), // nolint: lll
				err: nil,
			},
		},
		{
			name: "TextFormat",
			args: args{
				reader: bytes.NewBuffer([]byte{
					119, 255, 129, 3, 1, 1, 6, 82, 101, 115, 117, 108, 116, 1, 255, 130, 0, 1, 9, 1, 6, 65, 116, 116, 97, 99, 107, 1, 12, 0, 1, 3, 83, 101, 113, 1, 6, 0, 1, 4, 67, 111, 100, 101, 1, 6, 0, 1, 9, 84, 105, 109, 101, 115, 116, 97, 109, 112, 1, 255, 132, 0, 1, 7, 76, 97, 116, 101, 110, 99, 121, 1, 4, 0, 1, 8, 66, 121, 116, 101, 115, 79, 117, 116, 1, 6, 0, 1, 7, 66, 121, 116, 101, 115, 73, 110, 1, 6, 0, 1, 5, 69, 114, 114, 111, 114, 1, 12, 0, 1, 4, 66, 111, 100, 121, 1, 10, 0, 0, 0, 16, 255, 131, 5, 1, 1, 4, 84, 105, 109, 101, 1, 255, 132, 0, 0, 0, 68, 255, 130, 1, 36, 56, 56, 97, 100, 97, 52, 100, 100, 45, 101, 51, 97, 98, 45, 52, 50, 49, 56, 45, 98, 101, 102, 100, 45, 98, 98, 100, 48, 97, 52, 98, 50, 99, 101, 99, 102, 2, 255, 200, 1, 15, 1, 0, 0, 0, 14, 212, 12, 177, 127, 57, 196, 111, 216, 1, 74, 1, 251, 3, 5, 4, 97, 190, 0, // nolint: lll
				}),
				id:     "id",
				format: NewTextFormat(),
			},
			want: want{
				byt: []byte(`ID id
Requests      [total, rate]            1, 1.00
Duration      [total, attack, wait]    6.484537567s, 0s, 6.484537567s
Latencies     [mean, 50, 95, 99, max]  6.484537567s, 6.484537567s, 6.484537567s, 6.484537567s, 6.484537567s
Bytes In      [total, mean]            0, 0.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:1  
Error Set:`),
				err: nil,
			},
		},
		{
			name: "HistogramFormat",
			args: args{
				reader: bytes.NewBuffer([]byte{
					119, 255, 129, 3, 1, 1, 6, 82, 101, 115, 117, 108, 116, 1, 255, 130, 0, 1, 9, 1, 6, 65, 116, 116, 97, 99, 107, 1, 12, 0, 1, 3, 83, 101, 113, 1, 6, 0, 1, 4, 67, 111, 100, 101, 1, 6, 0, 1, 9, 84, 105, 109, 101, 115, 116, 97, 109, 112, 1, 255, 132, 0, 1, 7, 76, 97, 116, 101, 110, 99, 121, 1, 4, 0, 1, 8, 66, 121, 116, 101, 115, 79, 117, 116, 1, 6, 0, 1, 7, 66, 121, 116, 101, 115, 73, 110, 1, 6, 0, 1, 5, 69, 114, 114, 111, 114, 1, 12, 0, 1, 4, 66, 111, 100, 121, 1, 10, 0, 0, 0, 16, 255, 131, 5, 1, 1, 4, 84, 105, 109, 101, 1, 255, 132, 0, 0, 0, 68, 255, 130, 1, 36, 56, 56, 97, 100, 97, 52, 100, 100, 45, 101, 51, 97, 98, 45, 52, 50, 49, 56, 45, 98, 101, 102, 100, 45, 98, 98, 100, 48, 97, 52, 98, 50, 99, 101, 99, 102, 2, 255, 200, 1, 15, 1, 0, 0, 0, 14, 212, 12, 177, 127, 57, 196, 111, 216, 1, 74, 1, 251, 3, 5, 4, 97, 190, 0, // nolint: lll
				}),
				id:     "id",
				format: NewHistogramFormat(),
			},
			want: want{
				byt: []byte(`ID id
Bucket           #  %        Histogram
[0s,     500ms]  0  0.00%    
[500ms,  1s]     0  0.00%    
[1s,     1.5s]   0  0.00%    
[1.5s,   2s]     0  0.00%    
[2s,     2.5s]   0  0.00%    
[2.5s,   3s]     0  0.00%    
[3s,     +Inf]   1  100.00%  ###########################################################################`),
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.format.SetMeta("bucket", DefaultBucketString)
			b, e := CreateReportFromReader(tt.args.reader, tt.args.id, tt.args.format)
			if !reflect.DeepEqual(b, tt.want.byt) && !reflect.DeepEqual(e, tt.want.err) {
				t.Errorf("CreateReportFromReader() = %v, %v, want %v %v", b, e, tt.want.byt, tt.want.err)
			}
		})
	}
}
