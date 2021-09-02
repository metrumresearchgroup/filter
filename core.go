package filter

import (
	"bufio"
	"bytes"
)

func (fs Funcs) Apply(s []byte) ([]byte, error) {
	// a shortcut for single rows; trying to keep it efficient for
	// large datasets by skipping this step if the line is > 1000 bytes
	if len(s) < 1024 && !bytes.Contains(s, []byte{'\n'}) {
		return fs.applyRow(s), nil
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
		res := fs.applyRow(inScan.Bytes())

		if _, err := outBuffer.Write(res); err != nil {
			return nil, err
		}

		if err := outBuffer.WriteByte('\n'); err != nil {
			return nil, err
		}
	}

	res := outBuffer.Bytes()
	if addedNewline {
		res = bytes.TrimSuffix(res, []byte{'\n'})
	}

	return res, nil
}

func (fs Funcs) applyRow(s []byte) []byte {
	res := s
	for _, fn := range fs {
		out := fn(res)
		res = out
	}
	return res
}
