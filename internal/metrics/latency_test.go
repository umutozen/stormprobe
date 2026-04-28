package metrics

import "testing"

func TestPercentile(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		p        float64
		expected float64
	}{
		{"empty", nil, 0.95, 0},
		{"single", []float64{100}, 0.95, 100},
		{"p50", []float64{10, 20, 30, 40, 50}, 0.5, 30},
		{"p95", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 0.95, 19},
		{"p99 single element", []float64{42}, 0.99, 42},
		{"unsorted input", []float64{50, 10, 30, 20, 40}, 0.5, 30},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Percentile(tc.data, tc.p)
			if got != tc.expected {
				t.Errorf("Percentile(%v, %.2f) = %.0f, want %.0f", tc.data, tc.p, got, tc.expected)
			}
		})
	}
}

func TestPercentileSorted(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		p        float64
		expected float64
	}{
		{"empty", nil, 0.95, 0},
		{"p50 sorted", []float64{10, 20, 30, 40, 50}, 0.5, 30},
		{"p99 sorted", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0.99, 9},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PercentileSorted(tc.data, tc.p)
			if got != tc.expected {
				t.Errorf("PercentileSorted(%v, %.2f) = %.0f, want %.0f", tc.data, tc.p, got, tc.expected)
			}
		})
	}
}
