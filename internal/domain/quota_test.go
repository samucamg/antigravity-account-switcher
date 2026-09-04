package domain

import "testing"

func TestQuotaBucket_UsageFraction(t *testing.T) {
	cases := []struct {
		name      string
		remaining float64
		expected  float64
	}{
		{"full quota", 1.0, 0.0},
		{"half quota", 0.5, 0.5},
		{"quarter remaining", 0.25, 0.75},
		{"nearly exhausted", 0.05, 0.95},
		{"fully exhausted", 0.0, 1.0},
		{"negative remaining overshoot clamps to full usage", -0.5, 1.0},
		{"remaining above 1.0 clamps usage to zero", 2.0, 0.0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := &QuotaBucket{RemainingFraction: c.remaining}
			got := b.UsageFraction()
			if got != c.expected {
				t.Errorf("UsageFraction() with remaining %.2f = %.4f; want %.4f", c.remaining, got, c.expected)
			}
		})
	}
}

func TestQuotaBucket_IsUsageAboveThreshold(t *testing.T) {
	cases := []struct {
		name      string
		remaining float64
		threshold float64
		expected  bool
	}{
		// usage = 1.0 - remaining
		{"usage exactly at 80% threshold (remaining 0.20)", 0.20, 0.80, true},
		{"usage above 85% switch threshold (remaining 0.15)", 0.15, 0.85, true},
		{"usage below threshold (remaining 0.80 -> 20%)", 0.80, 0.80, false},
		{"full quota never warns", 1.0, 0.80, false},
		{"exhausted quota always warns", 0.0, 0.80, true},
		{"zero threshold never warns", 0.0, 0.0, false},
		{"negative threshold never warns", 0.0, -1.0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := &QuotaBucket{RemainingFraction: c.remaining}
			got := b.IsUsageAboveThreshold(c.threshold)
			if got != c.expected {
				t.Errorf("IsUsageAboveThreshold(%.2f) with remaining %.2f = %v; want %v",
					c.threshold, c.remaining, got, c.expected)
			}
		})
	}
}