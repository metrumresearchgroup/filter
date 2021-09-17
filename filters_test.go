package filter

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestDropEmpty(tt *testing.T) {
	tests := []struct {
		name string
		s    []byte
		want []byte
	}{
		{
			name: "not empty",
			s:    []byte("not empty"),
			want: []byte("not empty"),
		},
		{
			name: "not empty",
			s:    []byte(""),
			want: []byte(nil),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			got := DropEmpty(test.s)

			t.R.Equal(test.want, got)
		})
	}
}
