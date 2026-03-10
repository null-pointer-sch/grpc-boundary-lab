package core

import (
	"testing"
)

func TestGeneratePong(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "standard message",
			message:  "hello",
			expected: "pong: hello",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "pong: ping from frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePong(tt.message)
			if got != tt.expected {
				t.Errorf("GeneratePong() = %v, want %v", got, tt.expected)
			}
		})
	}
}
