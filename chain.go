package filter

// Chain is a slice of Applier type for the purpose of ordered execution.
type Chain []Filter

// NewChain creates a chain of Filter type.
func NewChain(fs ...Filter) Chain {
	return fs
}

// ApplyRow iterates []Filter in the Chain, applying them to the
// input data.
// This operation is not newline-aware.
func (c Chain) ApplyRow(s []byte) []byte {
	for _, f := range c {
		s = f.ApplyRow(s)
	}

	return s
}
