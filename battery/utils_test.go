package battery

import (
	"testing"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		sec      float64
		expected string
	}{
		{0, "00h 00m 00s"},
		{61, "00h 01m 01s"},
		{3661, "01h 01m 01s"},
		{7384, "02h 03m 04s"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatDuration(tt.sec)
			if result != tt.expected {
				t.Errorf("FormatDuration(%f) = %v, want %v", tt.sec, result, tt.expected)
			}
		})
	}
}
