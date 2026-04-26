package metrics

import (
	"errors"
	"testing"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		msg      string
		expected string
	}{
		{"context deadline exceeded", "timeout"},
		{"request Timeout exceeded", "timeout"},
		{"connection reset by peer", "reset"},
		{"forcibly closed the existing connection", "reset"},
		{"connection refused", "refused"},
		{"actively refused it", "refused"},
		{"unexpected EOF", "other"},
		{"no such host", "other"},
	}

	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) {
			got := ClassifyError(errors.New(tc.msg))
			if got != tc.expected {
				t.Errorf("ClassifyError(%q) = %q, want %q", tc.msg, got, tc.expected)
			}
		})
	}
}
