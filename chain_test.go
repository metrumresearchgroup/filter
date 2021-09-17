package filter_test

import (
	"bytes"
	"testing"

	"github.com/metrumresearchgroup/wrapt"

	. "github.com/metrumresearchgroup/filter"
)

func TestChain_ApplyRow(tt *testing.T) {
	tests := []struct {
		name string
		c    Chain
		s    []byte
		want []byte
	}{
		{
			name: "base",
			c:    NewChain(bytes.Title, bytes.TrimSpace),
			s:    []byte("hello, world  "),
			want: []byte("Hello, World"),
		},
		{
			name: "gigo",
			c:    NewChain(bytes.Title, bytes.TrimSpace),
			s:    []byte("hello\n world  "),
			want: []byte("Hello\n World"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			got := test.c.ApplyRow(test.s)

			t.R.Equal(test.want, got)
		})
	}
}
