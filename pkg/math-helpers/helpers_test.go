package math_helpers

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestDtr(t *testing.T) {
	tests := []struct {
		deg, want float64
	}{
		{0, 0},
		{180, math.Pi},
		{360, 2 * math.Pi},
		{90, math.Pi / 2},
		{-90, -math.Pi / 2},
		{45, math.Pi / 4},
	}
	for _, tt := range tests {
		got := Dtr(tt.deg)
		if !almostEqual(got, tt.want, eps) {
			t.Errorf("Dtr(%v) = %v, want %v", tt.deg, got, tt.want)
		}
	}
}

func TestDsin(t *testing.T) {
	tests := []struct {
		deg, want float64
	}{
		{0, 0},
		{30, 0.5},
		{90, 1},
		{180, 0},
		{270, -1},
		{360, 0},
	}
	for _, tt := range tests {
		got := Dsin(tt.deg)
		if !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("Dsin(%v) = %v, want %v", tt.deg, got, tt.want)
		}
	}
}

func TestDcos(t *testing.T) {
	tests := []struct {
		deg, want float64
	}{
		{0, 1},
		{60, 0.5},
		{90, 0},
		{180, -1},
		{270, 0},
		{360, 1},
	}
	for _, tt := range tests {
		got := Dcos(tt.deg)
		if !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("Dcos(%v) = %v, want %v", tt.deg, got, tt.want)
		}
	}
}

func TestConstrain(t *testing.T) {
	tests := []struct {
		in, want float64
	}{
		{0, 0},
		{180, 180},
		{360, 0},
		{370, 10},
		{720, 0},
		{-10, 350},
		{-370, 350},
		{45.5, 45.5},
	}
	for _, tt := range tests {
		got := Constrain(tt.in)
		if !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("Constrain(%v) = %v, want %v", tt.in, got, tt.want)
		}
		// invariant: result must be in [0, 360)
		if got < 0 || got >= 360 {
			t.Errorf("Constrain(%v) = %v out of [0,360)", tt.in, got)
		}
	}
}

func TestGetSignPrefix(t *testing.T) {
	tests := []struct {
		sign int
		want string
	}{
		{0, "+"},
		{1, "+"},
		{42, "+"},
		{-1, "-"},
		{-100, "-"},
	}
	for _, tt := range tests {
		got := GetSignPrefix(tt.sign)
		if got != tt.want {
			t.Errorf("GetSignPrefix(%d) = %q, want %q", tt.sign, got, tt.want)
		}
	}
}
