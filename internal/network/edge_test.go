package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_edge_isResidual(t *testing.T) {
	tests := []struct {
		name string
		e    edge
		want bool
	}{
		{
			name: "residual edge default zero initial capacity",
			e:    edge{},
			want: true,
		},
		{
			name: "non-residual edge with positive initial capacity",
			e:    edge{initialCapacity: 10},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.e.isResidual())
		})
	}
}
