package gateway_test

import (
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/gateway"
	"github.com/stretchr/testify/assert"
)

func TestStatsProvider_GetStats(t *testing.T) {
	provider := gateway.NewStatsProvider()

	tests := []struct {
		protocol string
		tls      bool
		expected bool
	}{
		{"rest", false, true},
		{"rest", true, true},
		{"grpc", false, true},
		{"grpc", true, true},
		{"unknown", false, false},
	}

	for _, tt := range tests {
		data, ok := provider.GetStats(tt.protocol, tt.tls)
		assert.Equal(t, tt.expected, ok)
		if ok {
			assert.Equal(t, tt.protocol, data.Protocol)
			assert.Equal(t, tt.tls, data.TLS)
		}
	}
}
