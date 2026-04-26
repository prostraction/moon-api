package julian_time

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	ma "moon/pkg/math-helpers"
)

func TestGetTimeFromLocation(t *testing.T) {
	tests := []struct {
		name          string
		location      *time.Location
		expectedHours int
		expectedMins  int
		expectedErr   error
		expectErr     bool
	}{
		// Error cases
		{
			name:        "nil location",
			location:    nil,
			expectedErr: errors.New("loc is nil"),
			expectErr:   true,
		},
		{
			name:        "empty timezone string",
			location:    time.FixedZone("", 0),
			expectedErr: errors.New("empty timezone string"),
			expectErr:   true,
		},
		{
			name:        "invalid timezone format",
			location:    time.FixedZone("invalid", 0),
			expectedErr: errors.New("invalid timezone format"),
			expectErr:   true,
		},

		// UTC/GMT cases
		{
			name:          "UTC",
			location:      time.UTC,
			expectedHours: 0,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "GMT",
			location:      time.FixedZone("GMT", 0),
			expectedHours: 0,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "UTC+0",
			location:      time.FixedZone("UTC+0", 0),
			expectedHours: 0,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "UTC-0",
			location:      time.FixedZone("UTC-0", 0),
			expectedHours: 0,
			expectedMins:  0,
			expectErr:     false,
		},

		// Positive offsets with colon
		{
			name:          "UTC+5:30",
			location:      time.FixedZone("UTC+5:30", 5*3600+30*60),
			expectedHours: 5,
			expectedMins:  30,
			expectErr:     false,
		},
		{
			name:          "UTC+05:30",
			location:      time.FixedZone("UTC+05:30", 5*3600+30*60),
			expectedHours: 5,
			expectedMins:  30,
			expectErr:     false,
		},
		{
			name:          "UTC+9:45",
			location:      time.FixedZone("UTC+9:45", 9*3600+45*60),
			expectedHours: 9,
			expectedMins:  45,
			expectErr:     false,
		},

		// Negative offsets with colon
		{
			name:          "UTC-5:30",
			location:      time.FixedZone("UTC-5:30", -5*3600-30*60),
			expectedHours: -5,
			expectedMins:  30,
			expectErr:     false,
		},
		{
			name:          "UTC-08:00",
			location:      time.FixedZone("UTC-08:00", -8*3600),
			expectedHours: -8,
			expectedMins:  0,
			expectErr:     false,
		},

		// Positive offsets without colon (4 digits)
		{
			name:          "UTC+0530",
			location:      time.FixedZone("UTC+0530", 5*3600+30*60),
			expectedHours: 5,
			expectedMins:  30,
			expectErr:     false,
		},
		{
			name:          "UTC+1230",
			location:      time.FixedZone("UTC+1230", 12*3600+30*60),
			expectedHours: 12,
			expectedMins:  30,
			expectErr:     false,
		},

		// Negative offsets without colon (4 digits)
		{
			name:          "UTC-0530",
			location:      time.FixedZone("UTC-0530", -5*3600-30*60),
			expectedHours: -5,
			expectedMins:  30,
			expectErr:     false,
		},

		// 3-digit offsets (hours + minutes)
		{
			name:          "UTC+530",
			location:      time.FixedZone("UTC+530", 5*3600+30*60),
			expectedHours: 5,
			expectedMins:  30,
			expectErr:     false,
		},
		{
			name:          "UTC-530",
			location:      time.FixedZone("UTC-530", -5*3600-30*60),
			expectedHours: -5,
			expectedMins:  30,
			expectErr:     false,
		},

		// Simple hour offsets
		{
			name:          "UTC+5",
			location:      time.FixedZone("UTC+5", 5*3600),
			expectedHours: 5,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "UTC-3",
			location:      time.FixedZone("UTC-3", -3*3600),
			expectedHours: -3,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "UTC+12",
			location:      time.FixedZone("UTC+12", 12*3600),
			expectedHours: 12,
			expectedMins:  0,
			expectErr:     false,
		},

		// Edge cases
		{
			name:          "UTC+14:00",
			location:      time.FixedZone("UTC+14:00", 14*3600),
			expectedHours: 14,
			expectedMins:  0,
			expectErr:     false,
		},
		{
			name:          "UTC-12:00",
			location:      time.FixedZone("UTC-12:00", -12*3600),
			expectedHours: -12,
			expectedMins:  0,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hours, minutes, err := GetTimeFromLocation(tt.location)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if hours != tt.expectedHours {
				t.Errorf("Expected hours %d, got %d", tt.expectedHours, hours)
			}

			if minutes != tt.expectedMins {
				t.Errorf("Expected minutes %d, got %d", tt.expectedMins, minutes)
			}
		})
	}
}

func TestGetTimeFromLocation_ErrorCases(t *testing.T) {
	errorTests := []struct {
		name     string
		location *time.Location
		errMsg   string
	}{
		{
			name:     "invalid minutes in colon format",
			location: time.FixedZone("UTC+5:60", 5*3600+60*60),
			errMsg:   "invalid minutes",
		},
		{
			name:     "negative minutes",
			location: time.FixedZone("UTC+5:-30", 5*3600-30*60),
			errMsg:   "invalid minutes",
		},
		{
			name:     "invalid hours in colon format",
			location: time.FixedZone("UTC+ab:30", 0),
			errMsg:   "invalid hours",
		},
		{
			name:     "invalid 4-digit format",
			location: time.FixedZone("UTC+ab30", 0),
			errMsg:   "invalid hours",
		},
		{
			name:     "invalid 3-digit format",
			location: time.FixedZone("UTC+a30", 0),
			errMsg:   "invalid hours",
		},
		{
			name:     "malformed colon format",
			location: time.FixedZone("UTC+5:30:45", 0),
			errMsg:   "invalid timezone format",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetTimeFromLocation(tt.location)
			if err == nil {
				t.Errorf("Expected error, got nil")
			} else if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

///////////////////////////////////////////////////////////////////////////////////

func TestSetTimezoneLocFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLoc *time.Location
		wantErr bool
		errMsg  string
	}{
		// Empty and UTC cases
		{
			name:    "empty string",
			input:   "",
			wantLoc: time.UTC,
			wantErr: false,
		},
		{
			name:    "UTC",
			input:   "UTC",
			wantLoc: time.UTC,
			wantErr: false,
		},
		{
			name:    "UTC+0",
			input:   "UTC+0",
			wantLoc: time.UTC,
			wantErr: false,
		},
		{
			name:    "UTC-0",
			input:   "UTC-0",
			wantLoc: time.UTC,
			wantErr: false,
		},
		{
			name:    "UTC0",
			input:   "UTC0",
			wantLoc: time.UTC,
			wantErr: false,
		},

		// Positive offsets with colon
		{
			name:    "UTC+5:30",
			input:   "UTC+5:30",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},
		{
			name:    "UTC+12:00",
			input:   "UTC+12:00",
			wantLoc: time.FixedZone("UTC+12", 12*3600),
			wantErr: false,
		},
		{
			name:    "UTC+23:59 rejected (out of IANA range)",
			input:   "UTC+23:59",
			wantErr: true,
			errMsg:  "hours out of range",
		},
		{
			name:    "UTC+14:00 (IANA max)",
			input:   "UTC+14:00",
			wantLoc: time.FixedZone("UTC+14", 14*3600),
			wantErr: false,
		},

		// Negative offsets with colon
		{
			name:    "UTC-5:30",
			input:   "UTC-5:30",
			wantLoc: time.FixedZone("UTC-5:30", -5*3600-30*60),
			wantErr: false,
		},
		{
			name:    "UTC-12:00",
			input:   "UTC-12:00",
			wantLoc: time.FixedZone("UTC-12", -12*3600),
			wantErr: false,
		},

		// Positive offsets without colon
		{
			name:    "UTC+5",
			input:   "UTC+5",
			wantLoc: time.FixedZone("UTC+5", 5*3600),
			wantErr: false,
		},
		{
			name:    "UTC+0530",
			input:   "UTC+0530",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},
		{
			name:    "UTC+530",
			input:   "UTC+530",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},
		{
			name:    "UTC+12",
			input:   "UTC+12",
			wantLoc: time.FixedZone("UTC+12", 12*3600),
			wantErr: false,
		},

		// Negative offsets without colon
		{
			name:    "UTC-5",
			input:   "UTC-5",
			wantLoc: time.FixedZone("UTC-5", -5*3600),
			wantErr: false,
		},
		{
			name:    "UTC-0530",
			input:   "UTC-0530",
			wantLoc: time.FixedZone("UTC-5:30", -5*3600-30*60),
			wantErr: false,
		},
		{
			name:    "UTC-530",
			input:   "UTC-530",
			wantLoc: time.FixedZone("UTC-5:30", -5*3600-30*60),
			wantErr: false,
		},

		// Without UTC prefix
		{
			name:    "+5:30 without UTC",
			input:   "+5:30",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},
		{
			name:    "-5:30 without UTC",
			input:   "-5:30",
			wantLoc: time.FixedZone("UTC-5:30", -5*3600-30*60),
			wantErr: false,
		},
		{
			name:    "5:30 without UTC or sign",
			input:   "5:30",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},

		// Edge cases and whitespace
		{
			name:    "with whitespace",
			input:   " UTC+5:30 ",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},
		{
			name:    "case insensitive UTC",
			input:   "utc+5:30",
			wantLoc: time.FixedZone("UTC+5:30", 5*3600+30*60),
			wantErr: false,
		},

		// Error cases
		{
			name:    "invalid hours",
			input:   "UTC+abc:30",
			wantErr: true,
			errMsg:  "invalid hours",
		},
		{
			name:    "invalid minutes",
			input:   "UTC+5:abc",
			wantErr: true,
			errMsg:  "invalid minutes",
		},
		{
			name:    "minutes out of range",
			input:   "UTC+5:60",
			wantErr: true,
			errMsg:  "invalid minutes",
		},
		{
			name:    "hours out of range positive",
			input:   "UTC+15:00",
			wantErr: true,
			errMsg:  "hours out of range",
		},
		{
			name:    "hours out of range negative",
			input:   "UTC-15:00",
			wantErr: true,
			errMsg:  "hours out of range",
		},
		{
			name:    "invalid format",
			input:   "UTC+5:30:15",
			wantErr: true,
			errMsg:  "invalid timezone format",
		},
		{
			name:    "invalid format without colon",
			input:   "UTC+12345",
			wantErr: true,
			errMsg:  "invalid timezone format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLoc, err := SetTimezoneLocFromString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetTimezoneLocFromString(%q) expected error, got nil", tt.input)
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("SetTimezoneLocFromString(%q) error = %v, expected to contain %q", tt.input, err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("SetTimezoneLocFromString(%q) unexpected error: %v", tt.input, err)
				return
			}

			if gotLoc.String() != tt.wantLoc.String() {
				t.Errorf("SetTimezoneLocFromString(%q) = %v, want %v", tt.input, gotLoc, tt.wantLoc)
			}

			testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, gotLoc)
			expectedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, tt.wantLoc)

			if testTime.UTC() != expectedTime.UTC() {
				t.Errorf("SetTimezoneLocFromString(%q) offset mismatch: got %v, want %v",
					tt.input, testTime.UTC(), expectedTime.UTC())
			}
		})
	}
}

func TestGetSignPrefix(t *testing.T) {
	tests := []struct {
		sign     int
		expected string
	}{
		{1, "+"},
		{-1, "-"},
		{0, "+"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("sign%d", tt.sign), func(t *testing.T) {
			result := ma.GetSignPrefix(tt.sign)
			if result != tt.expected {
				t.Errorf("GetSignPrefix(%d) = %q, want %q", tt.sign, result, tt.expected)
			}
		})
	}
}

// Benchmark test for performance
func BenchmarkSetTimezoneLocFromString(b *testing.B) {
	testCases := []string{
		"UTC+5:30",
		"UTC-8:00",
		"UTC+0",
		"UTC+1230",
		"+9:00",
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			SetTimezoneLocFromString(tc)
		}
	}
}

/////////////////////////////////////////////////////

func TestToJulianDate(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		expected  float64
		tolerance float64
	}{
		{
			name:      "January 1, 2000 12:00:00 UTC",
			input:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  2451545.0,
			tolerance: 0.0001,
		},
		{
			name:      "January 1, 2000 00:00:00 UTC",
			input:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  2451544.5,
			tolerance: 0.0001,
		},
		{
			name:      "July 4, 1976 14:30:00 UTC",
			input:     time.Date(1976, 7, 4, 14, 30, 0, 0, time.UTC),
			expected:  2442964.10417,
			tolerance: 0.0001,
		},
		{
			name:      "November 17, 1858 00:00:00 UTC",
			input:     time.Date(1858, 11, 17, 0, 0, 0, 0, time.UTC),
			expected:  2400000.5,
			tolerance: 0.0001,
		},
		{
			name:      "February 28, 2023 23:59:59 UTC",
			input:     time.Date(2023, 2, 28, 23, 59, 59, 0, time.UTC),
			expected:  2460004.49999,
			tolerance: 0.0001,
		},
		{
			name:      "December 31, 1999 23:59:59 UTC",
			input:     time.Date(1999, 12, 31, 23, 59, 59, 0, time.UTC),
			expected:  2451544.49999,
			tolerance: 0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToJulianDate(tt.input)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("ToJulianDate(%v) = %v, diff = %v",
					tt.input, result, (result - tt.expected))
			}
		})
	}
}

func TestFromJulianDate(t *testing.T) {
	tests := []struct {
		name     string
		julian   float64
		loc      *time.Location
		expected time.Time
	}{
		{
			name:     "JD 2451545.0 (January 1, 2000 12:00:00 UTC)",
			julian:   2451545.0,
			loc:      time.UTC,
			expected: time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "JD 2451544.5 (January 1, 2000 00:00:00 UTC)",
			julian:   2451544.5,
			loc:      time.UTC,
			expected: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "JD 2442964.10417 (July 4, 1976 14:30:00 UTC)",
			julian:   2442964.10417,
			loc:      time.UTC,
			expected: time.Date(1976, 7, 4, 14, 30, 0, 0, time.UTC),
		},
		{
			name:     "JD 2400000.5 (November 17, 1858 00:00:00 UTC)",
			julian:   2400000.5,
			loc:      time.UTC,
			expected: time.Date(1858, 11, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "With different timezone",
			julian:   2451545.0,
			loc:      time.FixedZone("EST", -5*3600),
			expected: time.Date(2000, 1, 1, 7, 0, 0, 0, time.FixedZone("EST", -5*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromJulianDate(tt.julian, tt.loc)

			// Compare individual components to avoid nanosecond precision issues
			if result.Year() != tt.expected.Year() ||
				result.Month() != tt.expected.Month() ||
				result.Day() != tt.expected.Day() ||
				result.Hour() != tt.expected.Hour() ||
				result.Minute() != tt.expected.Minute() ||
				result.Second() != tt.expected.Second() ||
				result.Location().String() != tt.expected.Location().String() {
				t.Errorf("FromJulianDate(%v, %v) = %v, expected %v",
					tt.julian, tt.loc, result, tt.expected)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	testCases := []struct {
		name  string
		input time.Time
	}{
		{"Modern date", time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)},
		{"Leap year date", time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC)},
		{"Historical date", time.Date(1776, 7, 4, 9, 0, 0, 0, time.UTC)},
		{"Future date", time.Date(2050, 6, 15, 18, 45, 30, 0, time.UTC)},
		//{"With timezone", time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*3600))}, unsupported
		{"Midnight", time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)},
		{"Noon", time.Date(2001, 9, 11, 12, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			julian := ToJulianDate(tc.input)
			converted := FromJulianDate(julian, tc.input.Location())

			// Allow for some floating point precision loss in time conversion
			if tc.input.Year() != converted.Year() ||
				tc.input.Month() != converted.Month() ||
				tc.input.Day() != converted.Day() ||
				math.Abs(float64(tc.input.Hour()-converted.Hour())) > 1 ||
				math.Abs(float64(tc.input.Minute()-converted.Minute())) > 1 {
				t.Errorf("Round trip conversion failed for %v\nOriginal: %v\nConverted: %v\nJulian: %v",
					tc.name, tc.input, converted, julian)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Zero Julian date", func(t *testing.T) {
		// This should handle very old dates gracefully
		result := FromJulianDate(0, time.UTC)
		if result.Year() == 0 {
			t.Error("FromJulianDate(0) returned zero time")
		}
	})

	t.Run("Negative Julian date", func(t *testing.T) {
		result := FromJulianDate(-1000, time.UTC)
		if result.Year() == 0 {
			t.Error("FromJulianDate(-1000) returned zero time")
		}
	})

	t.Run("Very large Julian date", func(t *testing.T) {
		result := FromJulianDate(10000000, time.UTC)
		if result.Year() < 10000 { // Should be far in the future
			t.Errorf("FromJulianDate(10000000) returned unexpected year: %v", result.Year())
		}
	})
}

func TestConsistencyWithKnownValues(t *testing.T) {
	// Test against some known Julian date conversions
	knownDates := []struct {
		time     time.Time
		expected float64
	}{
		{
			time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			2451545.0,
		},
		{
			time.Date(1987, 1, 27, 0, 0, 0, 0, time.UTC),
			2446822.5,
		},
	}

	for _, known := range knownDates {
		t.Run(known.time.String(), func(t *testing.T) {
			result := ToJulianDate(known.time)
			if math.Abs(result-known.expected) > 0.0001 {
				t.Errorf("Known date conversion failed: got %v, expected %v", result, known.expected)
			}

			// Test round trip
			converted := FromJulianDate(result, known.time.Location())
			if converted.Year() != known.time.Year() ||
				converted.Month() != known.time.Month() ||
				converted.Day() != known.time.Day() {
				t.Errorf("Round trip failed for known date: %v -> %v", known.time, converted)
			}
		})
	}
}
