package main

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/metrumresearchgroup/filter"
)

func main() {
	f := NewFilter(os.Stdout, os.Stdin, bytes.TrimSpace, bytes.ToLower, bytes.Title)

	if err := f.Wait(); err != nil {
		fmt.Println(err)
	}
}
