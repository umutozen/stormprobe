package report

import "testing"

func TestPickSuccessColor(t *testing.T) {
	tests := []struct {
		rate     float64
		expected string
	}{
		{100.0, colorSuccess},
		{99.0, colorSuccess},
		{95.0, colorSuccess},
		{94.9, colorWarning},
		{80.0, colorWarning},
		{79.9, colorDanger},
		{0.0, colorDanger},
		{40.7, colorDanger},
	}

	for _, tc := range tests {
		got := pickSuccessColor(tc.rate)
		if got != tc.expected {
			t.Errorf("pickSuccessColor(%.1f) = %q, want %q", tc.rate, got, tc.expected)
		}
	}
}
