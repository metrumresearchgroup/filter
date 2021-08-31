package filters

import "bytes"

func Capitalize(v []byte) []byte {
	return bytes.ToUpper(v)
}
