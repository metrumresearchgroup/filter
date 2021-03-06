package filter

import (
	"bytes"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestFilter_Apply(tt *testing.T) {
	tests := []struct {
		name string
		f    Filter
		s    []byte
		want []byte
	}{
		{
			name: "apply",
			f:    bytes.ToLower,
			s:    []byte("ARG"),
			want: []byte("arg"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			got := test.f.ApplyRow(test.s)

			t.R.Equal(test.want, got)
		})
	}
}
