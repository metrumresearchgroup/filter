package filter

import (
	"io"
	"os"
)

// Funcs is a slice of Func type for the purpose of ordered execution.
type Funcs []Func

// AsChain converts Funcs into discrete concurrent filters.
func (ffl Funcs) AsChain(writer io.Writer, readCloser io.ReadCloser) (*Chain, error) {
	return NewChain(writer, readCloser, ffl...)
}

// A Chain is simply a structure listing completed filters. They are wired
// to one another via pipes, and will process concurrently.
type Chain struct {
	ffs []*Filter
}

// NewChain creates a filter chain, or a sequence of filters that
// concurrently shuffle data from input to output when they encounter
// a scan on newline.
func NewChain(writer io.Writer, readCloser io.ReadCloser, filterFns ...Func) (*Chain, error) {
	fc := &Chain{}
	var nextReader *os.File

	for i, filterFn := range filterFns {
		var chainWriter io.Writer
		var chainReader io.ReadCloser

		if i == 0 {
			chainReader = readCloser
		} else {
			chainReader = nextReader
		}

		if i == len(filterFns)-1 {
			chainWriter = writer
		} else {
			pipeReader, pipeWriter, err := os.Pipe()
			if err != nil {
				return nil, err
			}
			chainWriter = pipeWriter
			nextReader = pipeReader
		}

		filter := filterFn.Filter(chainWriter, chainReader)
		fc.ffs = append(fc.ffs, filter)
	}

	return fc, nil
}

// Close shuts down all filters via their close method, in order.
func (fc *Chain) Close() []error {
	errs := make([]error, 0, len(fc.ffs))
	for _, filter := range fc.ffs {
		err := filter.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// Wait holds further processing until all filters cleanly exit.
// Usually called after calling Close.
func (fc *Chain) Wait() []error {
	errs := make([]error, 0, len(fc.ffs))
	for _, filter := range fc.ffs {
		err := filter.Wait()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
