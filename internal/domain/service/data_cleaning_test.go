package service

import (
	"testing"
	"time"
)

func TestTruncateToDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "truncates to midnight",
			input:    time.Date(2026, 2, 19, 14, 35, 22, 123456789, time.UTC),
			expected: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "already at midnight stays same",
			input:    time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "end of day truncates to start",
			input:    time.Date(2026, 2, 19, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "preserves timezone",
			input:    time.Date(2026, 2, 19, 14, 35, 22, 0, time.FixedZone("CST", 8*3600)),
			expected: time.Date(2026, 2, 19, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToDay(tt.input)

			if result.Year() != tt.expected.Year() {
				t.Errorf("Year mismatch: got %d, want %d", result.Year(), tt.expected.Year())
			}
			if result.Month() != tt.expected.Month() {
				t.Errorf("Month mismatch: got %d, want %d", result.Month(), tt.expected.Month())
			}
			if result.Day() != tt.expected.Day() {
				t.Errorf("Day mismatch: got %d, want %d", result.Day(), tt.expected.Day())
			}
			if result.Hour() != 0 {
				t.Errorf("Hour should be 0: got %d", result.Hour())
			}
			if result.Minute() != 0 {
				t.Errorf("Minute should be 0: got %d", result.Minute())
			}
			if result.Second() != 0 {
				t.Errorf("Second should be 0: got %d", result.Second())
			}
			if result.Nanosecond() != 0 {
				t.Errorf("Nanosecond should be 0: got %d", result.Nanosecond())
			}
		})
	}
}

func TestTruncateToDay_Consistency(t *testing.T) {
	// Test that multiple calls with same day produce same result
	base := time.Date(2026, 2, 19, 14, 35, 22, 0, time.UTC)

	results := make([]time.Time, 10)
	for i := 0; i < 10; i++ {
		// Different times on same day
		input := base.Add(time.Duration(i) * time.Hour)
		results[i] = truncateToDay(input)
	}

	// All results should be identical
	for i := 1; i < 10; i++ {
		if !results[i].Equal(results[0]) {
			t.Errorf("Results[%d] = %v, want %v", i, results[i], results[0])
		}
	}
}
