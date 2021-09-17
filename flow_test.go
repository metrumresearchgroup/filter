package filter

import (
	"bytes"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestFlow_Apply(tt *testing.T) {
	tests := []struct {
		name string
		flow Flow
		s    []byte
		want []byte
	}{
		{
			name: "base",
			flow: NewFlow(NewChain(bytes.ToLower, bytes.Title)),
			s:    []byte("hEllO, WoRld\nI aM aWare"),
			want: []byte("Hello, World\nI Am Aware"),
		},
		{
			name: "shortcut flow",
			flow: NewFlow(NewChain(bytes.ToLower, bytes.Title)),
			s:    []byte("hEllO, WoRld"),
			want: []byte("Hello, World"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			got := test.flow.Apply(test.s)

			t.R.Equal(string(test.want), string(got))
		})
	}
}
