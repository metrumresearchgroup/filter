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
		{
			name: "short single row",
			fs:   Funcs([]Func{bytes.Title, bytes.TrimSpace}),
			args: args{s: []byte("hello, world")},
			want: []byte("Hello, World"),
		},
		{
			name: "long single row",
			fs:   Funcs([]Func{}),
			args: args{s: []byte(Lipsum)},
			want: []byte(Lipsum),
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

const Lipsum = "Nunc a mi tincidunt, aliquam augue id, fermentum risus. Pellentesque ultricies bibendum ante, nec blandit odio eleifend sed. Aliquam egestas viverra ante ac sollicitudin. Vestibulum sit amet sapien id tortor eleifend pharetra eget non erat. Phasellus feugiat turpis at urna tristique, quis pulvinar neque auctor. Mauris at fermentum velit, nec rhoncus urna. Aliquam id mauris quam. Pellentesque sagittis leo orci, in dictum tortor rhoncus vel. Morbi in ultricies ipsum, non eleifend diam. In pharetra aliquam varius. Interdum et malesuada fames ac ante ipsum primis in faucibus. Vivamus pulvinar est enim, in tincidunt ligula volutpat ultrices. Nullam sit amet tincidunt odio, eu consequat dolor. Nam vitae vulputate turpis. Morbi porttitor leo et pulvinar ultrices. Fusce imperdiet, nulla vel molestie venenatis, tortor lorem semper libero, laoreet bibendum mi nibh eu diam.\n\nUt molestie leo id vehicula elementum. Integer ut pulvinar enim, et ultrices enim. Sed ut nunc nunc. Sed sit amet sem non diam porta lacinia id proin."
