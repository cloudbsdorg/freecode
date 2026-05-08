package util

import "testing"

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		secs     int
		expected string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"negative large", -100, ""},
		{"one second", 1, "1s"},
		{"45 seconds", 45, "45s"},
		{"59 seconds", 59, "59s"},
		{"one minute", 60, "1m"},
		{"one minute 30 seconds", 90, "1m 30s"},
		{"5 minutes", 300, "5m"},
		{"59 minutes 59 seconds", 3599, "59m 59s"},
		{"one hour", 3600, "1h"},
		{"one hour one minute", 3660, "1h 1m"},
		{"one hour one minute one second", 3661, "1h 1m"},
		{"2 hours 30 minutes", 9000, "2h 30m"},
		{"23 hours 59 minutes", 86340, "23h 59m"},
		{"one day", 86400, "~1 day"},
		{"2 days", 172800, "~2 days"},
		{"6 days 23 hours 59 minutes", 604739, "~6 days"},
		{"one week", 604800, "~1 week"},
		{"2 weeks", 1209600, "~2 weeks"},
		{"52 weeks", 31449600, "~52 weeks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.secs)
			if result != tt.expected {
				t.Errorf("FormatDuration(%d) = %q, want %q", tt.secs, result, tt.expected)
			}
		})
	}
}

func TestFormatDurationExamples(t *testing.T) {
	examples := []struct {
		secs     int
		expected string
	}{
		{45, "45s"},
		{90, "1m 30s"},
		{3661, "1h 1m"},
		{90061, "~1 day"},
		{604801, "~1 week"},
	}

	for _, ex := range examples {
		result := FormatDuration(ex.secs)
		if result != ex.expected {
			t.Errorf("FormatDuration(%d) = %q, want %q", ex.secs, result, ex.expected)
		}
	}
}