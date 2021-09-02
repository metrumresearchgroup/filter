package filter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"sync"
)

// A Filter holds everything needed to maintain a long-running piped filter
// processing a slice of bytes at a time.
type Filter struct {
	Err error

	reader io.ReadCloser
	writer io.Writer

	ff        Func
	wg        *sync.WaitGroup
	startOnce *sync.Once
}

// NewFilter creates a filter and binds it to a function of type
// func([]byte) []byte.
func NewFilter(writer io.Writer, readCLoser io.ReadCloser, ff Func) (*Filter, error) {
	if ff == nil {
		return nil, errors.New("passed nil filter func")
	}
	f := &Filter{
		reader:    readCLoser,
		writer:    writer,
		ff:        ff,
		wg:        &sync.WaitGroup{},
		startOnce: &sync.Once{},
	}
	f.start()

	return f, nil
}

// start kicks off the background process to monitor incoming new lines
// of text and sends them to the bound function to process. It will only
// run once, so it's prepared in a sync.Once block.
func (f *Filter) start() {
	f.startOnce.Do(func() {
		scanner := bufio.NewScanner(f.reader)

		f.wg.Add(1)

		go func() {
			defer f.wg.Done()

			var err error
			for scanner.Scan() {
				buf := bytes.Buffer{}
				buf.Write(f.ff(scanner.Bytes()))
				buf.WriteByte('\n')

				_, err = f.writer.Write(buf.Bytes())
				if err != nil {
					break
				}
			}

			f.Err = err
		}()
	})
}

// Close closes the input of the Filter, which terminates the main scan
// loop. This allows us to exit. It waits for the reader to close and
// scan work to complete before exiting.
func (f *Filter) Close() error {
	if err := f.reader.Close(); err != nil {
		return err
	}

	return f.Wait()
}

// Wait waits until the workgroup is terminated. This is useful in situations
// where a user would manually ^D on input at a terminal to exit. Wait
// would hold indefinitely until the workgroup is finished.
func (f *Filter) Wait() error {
	f.wg.Wait()

	return f.Err
}

type Func func([]byte) []byte

func (ff Func) AsFilter(writer io.Writer, readCloser io.ReadCloser) *Filter {
	f, err := NewFilter(writer, readCloser, ff)
	if err != nil {
		panic(errors.New("impossible situation"))
	}

	return f
}
