package filter

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// A Filter holds everything needed to maintain a long-running piped filter
// processing a slice of bytes at a time.
type Filter struct {
	Err error

	Funcs Funcs

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
		Funcs: ff,

		// IO
		readCloser: readCloser,
		writer:     writer,

		// Sync
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
				row := f.Funcs.applyRow(scanner.Bytes())
				if row == nil {
					continue
				}

				if _, err = f.writer.Write(row); err != nil {
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

// Func is any simple transform of a byte slice.
type Func func([]byte) []byte

// AsFilter creates a Filter from a single func.
func (ff Func) AsFilter(writer io.Writer, readCloser io.ReadCloser) *Filter {
	return NewFilter(writer, readCloser, ff)
}

// Funcs is a slice of Func type for the purpose of ordered execution.
type Funcs []Func

// AsFilter converts Funcs into a discrete, concurrent filter.
func (fs Funcs) AsFilter(writer io.Writer, readCloser io.ReadCloser) *Filter {
	return NewFilter(writer, readCloser, fs...)
}

// Apply process a slice of byte for rows of input and applies all Func
// entries to each row.
func (fs Funcs) Apply(s []byte) []byte {
	// a shortcut for single rows; trying to keep it efficient for
	// large datasets by skipping this search if the line is > 1k
	if len(s) < 1024 && !bytes.Contains(s, []byte{'\n'}) {
		return fs.applyRow(s)
	}

	// we'll take a final newline off, but this makes for less
	// convoluted work.
	var addedNewline bool
	if !bytes.HasSuffix(s, []byte{'\n'}) {
		s = append(s, '\n')
		addedNewline = true
	}

	inScan := bufio.NewScanner(bytes.NewReader(s))
	outBuffer := &bytes.Buffer{}

	// scanning rows, and applying a row at a time to all filters.
	// This prevents lots of unnecessary construction.
	for inScan.Scan() {
		res := fs.applyRow(inScan.Bytes())

		// since applyRow can return nil for empty lines, check
		// for nil before writing.
		if res != nil {
			// We're not checking err since bytes.Buffer never
			// returns error.
			outBuffer.Write(res)
			outBuffer.WriteByte('\n')
		}
	}

	res := outBuffer.Bytes()
	if addedNewline {
		res = bytes.TrimSuffix(res, []byte{'\n'})
	}

	return res
}

func (fs Funcs) applyRow(s []byte) []byte {
	work := s
	for _, fn := range fs {
		work = fn(work)

		// Nothing left to work with. This allows a Func to drop rows.
		if work == nil {
			break
		}
	}

	return work
}

func DropEmpty(s []byte) []byte {
	if len(s) == 0 {
		return nil
	}

	return s
}
