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
