package filter

import (
	"bufio"
	"bytes"
)

// Flow provides the Apply function that takes an input and runs each
// row through the RowApplier.
type Flow struct {
	RowApplier RowApplier
}

// NewFlow takes any RowApplier and provides a flow from it.
func NewFlow(ra RowApplier) Flow {
	return Flow{
		RowApplier: ra,
	}
}

// Apply scans a slice of byte for rows of input and repeatedly
// calls RowApplier.ApplyRow() on each row.
func (f Flow) Apply(s []byte) []byte {
	// a shortcut for single rows; trying to keep it efficient for
	// large datasets by skipping this search if the line is > 1k
	if len(s) < 1024 && !bytes.Contains(s, []byte{'\n'}) {
		return f.RowApplier.ApplyRow(s)
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
		res := f.RowApplier.ApplyRow(inScan.Bytes())

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
