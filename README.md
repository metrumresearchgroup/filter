# filter

Filters on pipelines between readers and writers.

## Purpose

This library provides a method to wrap filters in the necessary input and output pipes to pull from readers and write to
writers. It is intended to work with line-based input, and requires newline characters to call the wrapped functions.

## Usage

The best way to use this library is to create a reusable list of function references of the type `func([]byte) []byte`
and wrap that in the `filter.FuncList` type, then calling its `AsChain()` function to create a new filter chain.

```go
var fl filter.FuncList = []filter.Func{bytes.TrimSpace, bytes.ToLower, bytes.Title}

// We're using Stdout and Stdin, because it's convenient for demonstration.
// If you're going to pipe in data, use os.Pipe instead of io.Pipe
f, _ := fl.AsChain(os.Stdout, os.Stdin)

// if you have reason to close out (done doing the process, found what you wanted to, etc.)
// normally you'd put this behind a timeout or other condition.
_ = f.Close()

// perform interactions, ctrl-d or ctrl-c to interrupt.
_ = f.Wait()
```

## Dependencies

This package depends upon [github.com/metrumresearchgroup/wrapt](https://github.com/metrumresearchgroup/wrapt), our
general testing library.