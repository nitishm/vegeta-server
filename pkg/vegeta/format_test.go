package vegeta

import (
	"reflect"
	"testing"
)

func TestNewFormat(t *testing.T) {
	type args struct {
		typ string
		key string
		val string
	}
	type want struct {
		typ string
		mta MetaInfo
		frm Format
	}
	tests := []struct {
		name  string
		args  args
		want  want
		setup func(w *want)
	}{
		{
			name: "type: json",
			args: args{
				typ: "json",
				key: "bucket",
				val: DefaultBucketString,
			},
			want: want{
				typ: "json",
				mta: nil,
			},
			setup: func(w *want) {
				x := JSONFormat("json")
				w.frm = &x
			},
		},
		{
			name: "type: text",
			args: args{
				typ: "text",
				key: "bucket",
				val: DefaultBucketString,
			},
			want: want{
				typ: "text",
				mta: nil,
			},
			setup: func(w *want) {
				x := TextFormat("text")
				w.frm = &x
			},
		},
		{
			name: "type: binary",
			args: args{
				typ: "binary",
				key: "bucket",
				val: DefaultBucketString,
			},
			want: want{
				typ: "binary",
				mta: nil,
			},
			setup: func(w *want) {
				x := BinaryFormat("binary")
				w.frm = &x
			},
		},
		{
			name: "type: histogram",
			args: args{
				typ: "histogram",
				key: "bucket",
				val: DefaultBucketString,
			},
			want: want{
				typ: "histogram",
				mta: MetaInfo{"bucket": DefaultBucketString},
			},
			setup: func(w *want) {
				w.frm = &HistogramFormat{
					repr: "histogram",
					meta: make(MetaInfo),
				}
			},
		},
		{
			name: "type: unknown",
			args: args{
				typ: "unknown",
				key: "bucket",
				val: DefaultBucketString,
			},
			want: want{
				typ: "json",
				mta: nil,
			},
			setup: func(w *want) {
				x := JSONFormat("json")
				w.frm = &x
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(&tt.want)

			f := NewFormat(tt.args.typ)
			if !reflect.DeepEqual(f, tt.want.frm) {
				t.Errorf("NewFormat() = %v, want %v", f, tt.want.frm)
			}
			if x := f.String(); !reflect.DeepEqual(x, tt.want.typ) {
				t.Errorf("String() = %v, want %v", x, tt.want.typ)
			}
			f.SetMeta(tt.args.key, tt.args.val)
			if m := f.Meta(); !reflect.DeepEqual(m, tt.want.mta) {
				t.Errorf("Meta() = %v, want %v", m, tt.want.mta)
			}
		})
	}
}
