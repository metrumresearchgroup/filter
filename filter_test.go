package filter_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/metrumresearchgroup/wrapt"

	"github.com/metrumresearchgroup/filter"
)

func TestFilter(tt *testing.T) {
	createDefaultFilter := func(t *wrapt.T, w io.Writer, r io.ReadCloser) *filter.Filter {
		passthroughFunc := func(bs []byte) []byte {
			return bs
		}
		// double passthrough to test chaining.
		return filter.Funcs([]filter.Func{passthroughFunc}).AsFilter(w, r)
	}
	tests := []struct {
		name                  string
		createFunc            func(t *wrapt.T, w io.Writer, r io.ReadCloser) *filter.Filter
		newWantErr            bool
		input, expectedOutput []byte
	}{
		{
			name:           "pass-through nothing",
			expectedOutput: []byte{},
		},
		{
			name:  "apply DropEmpty",
			input: []byte("hello\n\nworld\n"),
			createFunc: func(t *wrapt.T, w io.Writer, r io.ReadCloser) *filter.Filter {
				f := createDefaultFilter(t, w, r)
				f.Funcs = append(f.Funcs, filter.DropEmpty)
				return f
			},
			expectedOutput: []byte("hello\nworld\n"),
		},
		{
			name:           "pass-through unterminated row",
			input:          []byte("hello world"),
			expectedOutput: []byte{},
		},
		{
			name:           "pass-through terminated row",
			input:          []byte("hello world\n"),
			expectedOutput: []byte("hello world\n"),
		},
		{
			name:           "pass-through two-line row, second unterminated",
			input:          []byte("hello world\nhow are you?"),
			expectedOutput: []byte("hello world\n"),
		},
		{
			name:           "pass-through two-line row, second terminated",
			input:          []byte("hello world\nhow are you?\n"),
			expectedOutput: []byte("hello world\nhow are you?\n"),
		},
	}

	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			if test.createFunc == nil {
				test.createFunc = createDefaultFilter
			}

			inputReader, inputWriter, err := os.Pipe()
			t.R.NoError(err)

			outputReader, outputWriter, err := os.Pipe()
			t.R.NoError(err)

			buf := &bytes.Buffer{}

			errCh := make(chan error, 1)
			go func() {
				_, copyErr := io.Copy(buf, outputReader)
				if copyErr != nil {
					errCh <- copyErr
				}
			}()

			f := test.createFunc(t, outputWriter, inputReader)

			n, err := inputWriter.Write(test.input)
			t.R.NoError(err)
			t.R.Equal(n, len(test.input))

			time.Sleep(100 * time.Millisecond)
			t.A.WantError(test.newWantErr, f.Err)
			t.A.Equal(test.expectedOutput, buf.Bytes())

			t.R.Empty(f.Close())
			t.R.Empty(f.Wait())
		})
	}
}

func TestFuncs_Apply(tt *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		fs   filter.Funcs
		args args
		want []byte
	}{
		{
			name: "base",
			fs:   filter.Funcs([]filter.Func{bytes.Title, bytes.TrimSpace}),
			args: args{s: []byte("hello, world\ntesting   \nremember me")},
			want: []byte("Hello, World\nTesting\nRemember Me"),
		},
		{
			name: "short single row",
			fs:   filter.Funcs([]filter.Func{bytes.Title, bytes.TrimSpace}),
			args: args{s: []byte("hello, world")},
			want: []byte("Hello, World"),
		},
		{
			name: "long single row",
			fs:   filter.Funcs([]filter.Func{}),
			args: args{s: []byte(Lipsum)},
			want: []byte(Lipsum),
		},
		{
			name: "empty row removed by DropEmpty",
			fs:   filter.Funcs([]filter.Func{filter.DropEmpty}),
			args: args{s: []byte("hello\n\nworld\n")},
			want: []byte("hello\nworld\n"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			got := test.fs.Apply(test.args.s)

			t.R.Equal(test.want, got)
		})
	}
}

const Lipsum = "Nunc a mi tincidunt, aliquam augue id, fermentum risus. Pellentesque ultricies bibendum ante, nec blandit odio eleifend sed. Aliquam egestas viverra ante ac sollicitudin. Vestibulum sit amet sapien id tortor eleifend pharetra eget non erat. Phasellus feugiat turpis at urna tristique, quis pulvinar neque auctor. Mauris at fermentum velit, nec rhoncus urna. Aliquam id mauris quam. Pellentesque sagittis leo orci, in dictum tortor rhoncus vel. Morbi in ultricies ipsum, non eleifend diam. In pharetra aliquam varius. Interdum et malesuada fames ac ante ipsum primis in faucibus. Vivamus pulvinar est enim, in tincidunt ligula volutpat ultrices. Nullam sit amet tincidunt odio, eu consequat dolor. Nam vitae vulputate turpis. Morbi porttitor leo et pulvinar ultrices. Fusce imperdiet, nulla vel molestie venenatis, tortor lorem semper libero, laoreet bibendum mi nibh eu diam.\n\nUt molestie leo id vehicula elementum. Integer ut pulvinar enim, et ultrices enim. Sed ut nunc nunc. Sed sit amet sem non diam porta lacinia id proin."
