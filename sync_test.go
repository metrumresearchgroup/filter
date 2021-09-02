package filter

import (
	"bytes"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestFuncs_Apply(tt *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name    string
		fs      Funcs
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "base",
			fs:   Funcs([]Func{bytes.Title, bytes.TrimSpace}),
			args: args{s: []byte("hello, world\ntesting   \nremember me")},
			want: []byte("Hello, World\nTesting\nRemember Me"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			got, err := test.fs.Apply(test.args.s)

			t.R.WantError(test.wantErr, err)
			t.R.Equal(test.want, got)
		})
	}
}
