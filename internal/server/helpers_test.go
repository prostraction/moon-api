package server

import (
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	tests := []struct {
		in   float64
		want int
	}{
		{0.0, 0},
		{0.4, 0},
		{0.5, 1},
		{0.6, 1},
		{1.5, 2},
		{-0.5, -1},
		{-0.4, 0},
		{-1.5, -2},
	}
	for _, tt := range tests {
		got := round(tt.in)
		if got != tt.want {
			t.Errorf("round(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestToFixed(t *testing.T) {
	tests := []struct {
		num       float64
		precision int
		want      float64
	}{
		{1.234567, 2, 1.23},
		{1.235, 2, 1.24},
		{1.2345, 3, 1.235},
		{0, 5, 0},
		{-1.235, 2, -1.24},
		{100.0, 0, 100},
	}
	for _, tt := range tests {
		got := ToFixed(tt.num, tt.precision)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("ToFixed(%v, %d) = %v, want %v", tt.num, tt.precision, got, tt.want)
		}
	}
}

func TestStrToInt(t *testing.T) {
	tests := []struct {
		val      string
		fallback int
		min, max int
		want     int
	}{
		{"5", 0, 0, 10, 5},
		{"-3", 0, 0, 10, 0},     // below min, clamped to min
		{"100", 0, 0, 10, 10},   // above max, clamped to max
		{"abc", 7, 0, 10, 7},    // invalid → fallback (in range)
		{"abc", 100, 0, 10, 10}, // invalid → fallback then clamp
		{"abc", -5, 0, 10, 0},   // invalid → fallback then clamp
		{"", 5, 0, 10, 5},       // empty → fallback
		{"0", 0, 0, 0, 0},       // single value range
		{"5", 0, -10, 10, 5},
	}
	for _, tt := range tests {
		got := StrToInt(tt.val, tt.fallback, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("StrToInt(%q, %d, %d, %d) = %d, want %d",
				tt.val, tt.fallback, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestParseCoords(t *testing.T) {
	t.Run("valid coords", func(t *testing.T) {
		c := parseCoords("51.05", "71.43")
		if !c.IsValid {
			t.Fatal("expected valid")
		}
		if math.Abs(c.Latitude-51.05) > 1e-9 {
			t.Errorf("Latitude = %v", c.Latitude)
		}
		if math.Abs(c.Longitude-71.43) > 1e-9 {
			t.Errorf("Longitude = %v", c.Longitude)
		}
	})

	t.Run("missing values", func(t *testing.T) {
		c := parseCoords("no-value", "71.43")
		if c.IsValid {
			t.Error("expected invalid when latitude is no-value")
		}
		c = parseCoords("51.05", "no-value")
		if c.IsValid {
			t.Error("expected invalid when longitude is no-value")
		}
	})

	t.Run("malformed values", func(t *testing.T) {
		c := parseCoords("xyz", "71.43")
		if c.IsValid {
			t.Error("expected invalid for non-numeric latitude")
		}
		c = parseCoords("51.05", "xyz")
		if c.IsValid {
			t.Error("expected invalid for non-numeric longitude")
		}
	})

	t.Run("zero coords are valid", func(t *testing.T) {
		c := parseCoords("0", "0")
		if !c.IsValid {
			t.Error("expected valid for 0,0")
		}
	})

	t.Run("negative coords", func(t *testing.T) {
		c := parseCoords("-33.86", "-151.21")
		if !c.IsValid {
			t.Fatal("expected valid")
		}
		if c.Latitude != -33.86 || c.Longitude != -151.21 {
			t.Errorf("got lat=%v lon=%v", c.Latitude, c.Longitude)
		}
	})
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		name           string
		year, mon, day int
		wantErr        bool
	}{
		{"valid date", 2024, 6, 15, false},
		{"feb 29 leap year", 2024, 2, 29, false},
		{"feb 29 non-leap", 2023, 2, 29, true},
		{"feb 29 century non-leap", 1900, 2, 29, true},
		{"feb 29 400-year leap", 2000, 2, 29, false},
		{"april 31", 2024, 4, 31, true},
		{"day 0", 2024, 6, 0, true},
		{"month 0", 2024, 0, 15, true},
		{"month 13", 2024, 13, 15, true},
		{"year negative", -1, 6, 15, true},
		{"year 10000", 10000, 6, 15, true},
		{"year 0", 0, 6, 15, false},
		{"year 9999", 9999, 12, 31, false},
		{"jan 31", 2024, 1, 31, false},
		{"feb 28 non-leap", 2023, 2, 28, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidDate(tt.year, tt.mon, tt.day)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for %d-%d-%d", tt.year, tt.mon, tt.day)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error for %d-%d-%d: %v", tt.year, tt.mon, tt.day, err)
			}
		})
	}
}
