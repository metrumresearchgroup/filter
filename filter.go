package filter

// Filter is any simple transform of a byte slice.
type Filter func([]byte) []byte

// ApplyRow the Filter to a []byte, returning the result.
// This operation is not newline-aware.
func (f Filter) ApplyRow(s []byte) []byte {
	return f(s)
}

// RowApplier is any type that has a ApplyRow function.
type RowApplier interface {
	ApplyRow(s []byte) []byte
}
