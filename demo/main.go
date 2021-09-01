package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/metrumresearchgroup/filter"
)

func main() {
	var fl filter.Funcs = []filter.Func{bytes.TrimSpace, bytes.ToLower, bytes.Title}
	f, err := fl.AsChain(os.Stdout, os.Stdin)
	if err != nil {
		panic(err)
	}

	errs := f.Wait()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}
}
