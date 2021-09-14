package filter_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/metrumresearchgroup/wrapt"

	. "github.com/metrumresearchgroup/filter"
)

func TestStream(tt *testing.T) {
	passFunc := func(bs []byte) []byte {
		return bs
	}
	defaultApplier := Filter(passFunc)

	tests := []struct {
		name                  string
		rowApplier            RowApplier
		newWantErr            bool
		input, expectedOutput []byte
	}{
		{
			name:           "pass-through nothing",
			expectedOutput: []byte{},
		},
		{
			name:           "apply DropEmpty",
			input:          []byte("hello\n\nworld\n"),
			rowApplier:     NewChain(passFunc, DropEmpty),
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

			if test.rowApplier == nil {
				test.rowApplier = defaultApplier
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

			f := NewStream(outputWriter, inputReader, test.rowApplier)

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
