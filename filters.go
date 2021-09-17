package filter

// DropEmpty returns nil on an empty slice. This allows us
// to eliminate rows that have no meaning in output elsewhere.
func DropEmpty(s []byte) []byte {
	if len(s) == 0 {
		return nil
	}

	return s
}
