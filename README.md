# filter

## Purpose

This library provides a method to wrap filter functions in the necessary input and output pipes to pull from readers and
write to writers. It is intended to work with line-based input, and processes lines one at a time.

## Types

* `RowApplier` - an `interface` that represents anything with an `ApplyRow([]byte) []byte` function.
* `Filter` - a `func([]byte) []byte` for performing a transform of bytes. Typing it in a `Filter()` will give you
  the `ApplyRow` function, making it a `RowApplier`.
* `Chain` - an alias to`[]Filter`, operated on in sequence in order. Also Provides `ApplyRow`.
* `Stream` - an asynchronous type that takes a WriteCloser, a Reader, and a RowApplier, and scans input to output on
  newline.
* `Flow` - a synchronous type that takes a `RowApplier` and provides `Apply([]byte) []byte` and processes all input
  row-by-row until input is exhausted.

## Usage

The simplest case is using Filter to process a single row:

```go
var f := Filter(bytes.Trim)
r := f.ApplyRow([]byte("hello world"))
// r = "hello world"
```

To string filters together:

```go
c := NewChain(bytes.Trim, bytes.Title)
r := c.ApplyRow([]byte("  hello world"))
// r = "Hello World"
```

To apply to multiple rows, use Flow:

```go
// using compact chaining, you can safely save it as a variable.
f := NewFlow(NewChain(bytes.Trim, bytes.Title))
r := f.Apply([]byte("hello world\nhow are you?  "))
// r = "Hello World\nHow Are You?"
```

For live interaction, you can bind inputs and outputs to a `Stream`:

```go
// You can store the and re-use it to make multiple streams.
c := NewChain(bytes.TrimSpace, bytes.ToLower, bytes.Title)
f := NewStream(os.Stdout, os.Stdin, c)

// if you have reason to close out (done running the process, found what you wanted to, etc.)
// normally you'd put this behind a timeout or other condition.
_ = f.Close()

// perform interactions, ctrl-d to interrupt.
_ = f.Wait()
```

## Dependencies

This package depends upon [github.com/metrumresearchgroup/wrapt](https://github.com/metrumresearchgroup/wrapt), our
general testing library.