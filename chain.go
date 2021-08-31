package filter

import (
	"io"
	"os"
)

type FuncList []Func

func (ffl FuncList) AsChain(writer io.Writer, readCloser io.ReadCloser) (*Chain, error) {
	return NewChain(writer, readCloser, ffl...)
}

type Chain struct {
	ffs []*Filter
}

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

func (fc *Chain) Start() {
	for _, filter := range fc.ffs {
		filter.Start()
	}
}

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
