package filter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"sync"
)

type Filter struct {
	Err error

	reader io.ReadCloser
	writer io.Writer

	ff Func
	wg *sync.WaitGroup
}

func NewFilter(writer io.Writer, readCLoser io.ReadCloser, ff Func) (*Filter, error) {
	if ff == nil {
		return nil, errors.New("passed nil filter func")
	}

	return &Filter{
		reader: readCLoser,
		writer: writer,
		ff:     ff,
		wg:     &sync.WaitGroup{},
	}, nil
}

func (f *Filter) Start() {
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
}

func (f *Filter) Wait() error {
	f.wg.Wait()

	return f.Err
}

func (f *Filter) Close() error {
	if err := f.reader.Close(); err != nil {
		return err
	}
	return f.Wait()
}

type Func func([]byte) []byte

func (ff Func) Filter(writer io.Writer, readCloser io.ReadCloser) *Filter {
	f, err := NewFilter(writer, readCloser, ff)
	if err != nil {
		panic(errors.New("impossible situation"))
	}

	return f
}
