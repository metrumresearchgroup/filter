package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"scratch/filter"
	"scratch/filter/filters"
)

func main() {
	in := os.Stdin
	out := os.Stdout

	inReader, inWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(inWriter, in)

	var ffl filter.FuncList = []filter.Func{filters.Capitalize, CrazyCaps}
	f, err := ffl.AsChain(out, inReader)
	if err != nil {
		panic(err)
	}

	f.Start()
	errs := f.Wait()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}
}

// CrazyCaps proves you can define one of these functions everywhere.
func CrazyCaps(v []byte) []byte {
	buf := &bytes.Buffer{}
	for n, b := range v {
		if n%2 == 0 {
			buf.Write(bytes.ToUpper([]byte{b}))
		} else {
			buf.Write(bytes.ToLower([]byte{b}))
		}
	}
	return buf.Bytes()
}
