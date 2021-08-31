package filter_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/metrumresearchgroup/wrapt"

	"scratch/filter"
)

func TestFilter(tt *testing.T) {
	createDefaultChain := func(t *wrapt.T, w io.Writer, r io.ReadCloser) (*filter.Chain, error) {
		passthroughFunc := func(bs []byte) []byte {
			return bs
		}
		// double passthrough to test chaining.
		fl := filter.FuncList([]filter.Func{passthroughFunc, passthroughFunc})
		return fl.AsChain(w, r)
	}
	tests := []struct {
		name                  string
		createFunc            func(t *wrapt.T, w io.Writer, r io.ReadCloser) (*filter.Chain, error)
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
				test.createFunc = createDefaultChain
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

			f, err := test.createFunc(t, outputWriter, inputReader)
			t.R.WantError(test.newWantErr, err)

			f.Start()

			n, err := inputWriter.Write(test.input)
			t.R.NoError(err)
			t.R.Equal(n, len(test.input))

			time.Sleep(100 * time.Millisecond)

			t.A.Equal(test.expectedOutput, buf.Bytes())

			t.R.Empty(f.Close())
			t.R.Empty(f.Wait())
		})
	}
}
