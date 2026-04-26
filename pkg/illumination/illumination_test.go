package illumination

import (
	"math"
	"testing"
	"time"
)

func TestGetDailyMoonIllumination(t *testing.T) {
	tests := []struct {
		name     string
		tGiven   time.Time
		loc      *time.Location
		expected float64
	}{
		{
			name:     "Full moon reference date",
			tGiven:   time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 1.0,
		},
		{
			name:     "New moon (half cycle after full)",
			tGiven:   time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 0.0,
		},
		{
			name:     "First quarter (quarter cycle after full)",
			tGiven:   time.Date(2025, 9, 15, 5, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 0.43, // nasa = 49%
		},
		{
			name:     "With timezone offset",
			tGiven:   time.Date(2025, 9, 8, 0, 0, 0, 0, time.FixedZone("TEST", 3600)), // UTC+1
			loc:      time.FixedZone("TEST", 3600),
			expected: 1.0,
		},
		{
			name:     "Midnight UTC",
			tGiven:   time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDailyMoonIllumination(tt.tGiven, tt.loc)

			// Allow small floating point tolerance
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("%s GetDailyMoonIllumination() = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetCurrentMoonIllumination(t *testing.T) {
	tests := []struct {
		name     string
		tGiven   time.Time
		loc      *time.Location
		expected float64
	}{
		{
			name:     "Exact full moon time",
			tGiven:   time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 1.0,
		},
		{
			name:     "Specific time on full moon day",
			tGiven:   time.Date(2025, 9, 7, 15, 30, 45, 0, time.UTC),
			loc:      time.UTC,
			expected: 0.999,
		},
		{
			name:     "New moon exact time",
			tGiven:   time.Date(2025, 9, 21, 18, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 0.0,
		},
		{
			name:     "With positive timezone offset",
			tGiven:   time.Date(2025, 9, 8, 1, 0, 0, 0, time.FixedZone("UTC+1", 3600)),
			loc:      time.FixedZone("UTC+1", 3600),
			expected: 1.0,
		},
		{
			name:     "With negative timezone offset",
			tGiven:   time.Date(2025, 9, 7, 23, 0, 0, 0, time.FixedZone("UTC-1", -3600)),
			loc:      time.FixedZone("UTC-1", -3600),
			expected: 1.0,
		},
		{
			name:     "Random time between phases",
			tGiven:   time.Date(2025, 9, 12, 13, 0, 0, 0, time.UTC),
			loc:      time.UTC,
			expected: 0.71,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCurrentMoonIllumination(tt.tGiven, tt.loc)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("%s GetCurrentMoonIllumination() = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestIlluminationFunctionsComparison(t *testing.T) {
	testTime := time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC)

	daily := GetDailyMoonIllumination(testTime, time.UTC)
	current := GetCurrentMoonIllumination(testTime, time.UTC)

	if math.Abs(daily-current) > 0.001 {
		t.Errorf("Functions should return same value for midnight. Daily: %v, Current: %v", daily, current)
	}
}

func TestIlluminationRange(t *testing.T) {
	testTimes := []time.Time{
		time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),   // Full
		time.Date(2025, 9, 21, 18, 0, 0, 0, time.UTC), // New
		time.Date(2025, 9, 14, 9, 0, 0, 0, time.UTC),  // First quarter
		time.Date(2025, 9, 28, 9, 0, 0, 0, time.UTC),  // Last quarter
		time.Now(), // Current time
	}

	for _, testTime := range testTimes {
		daily := GetDailyMoonIllumination(testTime, time.UTC)
		current := GetCurrentMoonIllumination(testTime, time.UTC)

		if daily < 0 || daily > 1 {
			t.Errorf("GetDailyMoonIllumination() returned out of range value: %v", daily)
		}

		if current < 0 || current > 1 {
			t.Errorf("GetCurrentMoonIllumination() returned out of range value: %v", current)
		}
	}
}
