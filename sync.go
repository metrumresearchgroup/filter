package filter

import (
	"bufio"
	"bytes"
)

func (ff Func) Apply(s []byte) ([]byte, error) {
	var addedNewline bool
	if !bytes.HasSuffix(s, []byte{'\n'}) {
		s = append(s, '\n')
		addedNewline = true
	}

	inScan := bufio.NewScanner(bytes.NewReader(s))
	outBuffer := &bytes.Buffer{}

	for inScan.Scan() {
		res := ff(inScan.Bytes())
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

func (fs Funcs) Apply(s []byte) ([]byte, error) {
	res := s
	for _, fn := range fs {
		out, err := fn.Apply(res)
		if err != nil {
			return nil, err
		}
		res = out
	}

	return res, nil
}
