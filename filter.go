package filter

import (
	"bufio"
	"io"
	"sync"
)

// A Filter holds everything needed to maintain a long-running piped filter
// processing a slice of bytes at a time.
type Filter struct {
	Err error

	funcs Funcs

	readCloser io.ReadCloser
	writer     io.Writer

	// synchronization
	wg        *sync.WaitGroup
	startOnce *sync.Once
}

// NewFilter creates a filter and binds it to a function of type
// func([]byte) []byte.
func NewFilter(writer io.Writer, readCloser io.ReadCloser, ff ...Func) *Filter {
	f := &Filter{
		readCloser: readCloser,
		writer:     writer,
		funcs:      ff,

		// synchronization
		wg:        &sync.WaitGroup{},
		startOnce: &sync.Once{},
	}

	f.start()

	return f
}

// start kicks off the background process to monitor incoming new lines
// of text and sends them to the bound function to process. It will only
// run once, so it's prepared in a sync.Once block.
func (f *Filter) start() {
	f.startOnce.Do(func() {
		scanner := bufio.NewScanner(f.readCloser)

		f.wg.Add(1)

		// note on debugging tests: they have a very short window of operation
		// and they will time out if you create a breakpoint inside this gofunc.
		go func() {
			defer f.wg.Done()

			var err error

			for scanner.Scan() {
				if _, err = f.writer.Write(f.funcs.applyRow(scanner.Bytes())); err != nil {
					break
				}

				if _, err = f.writer.Write([]byte{'\n'}); err != nil {
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
	if err := f.readCloser.Close(); err != nil {
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
	return NewFilter(writer, readCloser, ff)
}
