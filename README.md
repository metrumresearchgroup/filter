# filter

Filters on pipelines between readers and writers.

## Purpose

This library provides a method to wrap filters in the necessary input and output pipes to pull from readers and write to
writers. It is intended to work with line-based input, and requires newline characters to call the wrapped functions.

## Types

* `Func` - a `func([]byte) []byte` for performing a transform of bytes with newlines.
* `Funcs` - a `[]Func`, also capable of performing transforms of bytes with or without newlines.
* `Filter` - an asynchronous type that takes a WriteCloser and a Reader, scans input to output on newline.

## Usage

The easiest way to use this library is to create a list of []Func and re-type as `filter.Funcs`. You can then
run `Funcs.Apply([]byte) []byte` on your data.

```go
var fs filter.Funcs = []filter.Func{bytes.Trim, bytes.ToUpper}
r := fs.Apply([]byte("hello world  "))
// r = "HELLO WORLD"
```

For live interaction, you can bind inputs and outputs to a `Filter`. The best way do this is to create a reusable list
of function references of the type `func([]byte) []byte` and wrap that in the `filter.Funcs` type, then calling
its `AsChain()` function to create a new filter chain.

You can create a Filter to process items concurrently.

The short form:

```go
f := filter.NewFilter(os.Stdout, os.Stdin, bytes.TrimSpace, bytes.ToLower, bytes.Title)
// if you have reason to close out (done doing the process, found what you wanted to, etc.)
// normally you'd put this behind a timeout or other condition.
_ = f.Close()

// perform interactions, ctrl-d or ctrl-c to interrupt.
_ = f.Wait()
```

The reusable form:
```go
var fs filter.Funcs = []filter.Func{bytes.TrimSpace, bytes.ToLower, bytes.Title}

// We're using Stdout and Stdin, because it's convenient for demonstration.
// If you're going to pipe in data, use os.Pipe instead of io.Pipe.
f, _ := fs.AsFilter(os.Stdout, os.Stdin)

// if you have reason to close out (done doing the process, found what you wanted to, etc.)
// normally you'd put this behind a timeout or other condition.
_ = f.Close()

// perform interactions, ctrl-d or ctrl-c to interrupt.
_ = f.Wait()
```

## Dependencies

This package depends upon [github.com/metrumresearchgroup/wrapt](https://github.com/metrumresearchgroup/wrapt), our
general testing library.