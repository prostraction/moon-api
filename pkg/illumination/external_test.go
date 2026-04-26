package illumination

import (
	"math"
	"testing"
	"time"

	jt "moon/pkg/julian-time"
)

// External reference data from timeanddate.com/moon/phases.
//
// The algorithm here uses Meeus's illuminated-fraction formula. Compared
// against timeanddate.com it agrees to within ~1-2 percentage points, with
// the largest residual at quarter-phase boundaries.

// At the exact UTC moment of a New Moon, illumination must be near zero (<2%).
func TestIllumination_NewMoonsNearZero(t *testing.T) {
	events := []struct {
		name string
		when time.Time
	}{
		{"2023-01-21 New Moon", time.Date(2023, 1, 21, 20, 53, 0, 0, time.UTC)},
		{"2023-02-20 New Moon", time.Date(2023, 2, 20, 7, 6, 0, 0, time.UTC)},
		{"2024-01-11 New Moon", time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC)},
		{"2024-02-09 New Moon", time.Date(2024, 2, 9, 22, 59, 0, 0, time.UTC)},
		{"2024-04-08 Total Solar Eclipse", time.Date(2024, 4, 8, 18, 21, 0, 0, time.UTC)},
		{"2025-01-29 New Moon", time.Date(2025, 1, 29, 12, 36, 0, 0, time.UTC)},
		{"2025-03-29 Partial Solar Eclipse", time.Date(2025, 3, 29, 10, 58, 0, 0, time.UTC)},
	}
	for _, ev := range events {
		t.Run(ev.name, func(t *testing.T) {
			illum := GetCurrentMoonIllumination(ev.when, time.UTC)
			if illum > 0.02 {
				t.Errorf("New Moon %v: illumination = %.4f, expected < 0.02", ev.when, illum)
			}
		})
	}
}

// At the exact UTC moment of a Full Moon, illumination must be very high (>0.98).
func TestIllumination_FullMoonsNearOne(t *testing.T) {
	events := []struct {
		name string
		when time.Time
	}{
		{"2023-01-06 Wolf Moon", time.Date(2023, 1, 6, 23, 8, 0, 0, time.UTC)},
		{"2023-02-05 Snow Moon", time.Date(2023, 2, 5, 18, 29, 0, 0, time.UTC)},
		{"2023-07-03 Buck Moon", time.Date(2023, 7, 3, 11, 39, 0, 0, time.UTC)},
		{"2024-01-25 Wolf Moon", time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC)},
		{"2024-04-23 Pink Moon", time.Date(2024, 4, 23, 23, 48, 0, 0, time.UTC)},
		{"2024-09-18 Harvest Moon", time.Date(2024, 9, 18, 2, 34, 0, 0, time.UTC)},
		{"2025-01-13 Wolf Moon", time.Date(2025, 1, 13, 22, 27, 0, 0, time.UTC)},
		{"2025-03-14 Lunar Eclipse", time.Date(2025, 3, 14, 6, 55, 0, 0, time.UTC)},
	}
	for _, ev := range events {
		t.Run(ev.name, func(t *testing.T) {
			illum := GetCurrentMoonIllumination(ev.when, time.UTC)
			if illum < 0.98 {
				t.Errorf("Full Moon %v: illumination = %.4f, expected > 0.98", ev.when, illum)
			}
		})
	}
}

// At quarter phases, illumination must be approximately 50% (±5pp).
func TestIllumination_QuarterPhasesNearHalf(t *testing.T) {
	events := []struct {
		name string
		when time.Time
	}{
		{"2024-01-18 First Quarter", time.Date(2024, 1, 18, 3, 53, 0, 0, time.UTC)},
		{"2024-02-02 Last Quarter", time.Date(2024, 2, 2, 23, 18, 0, 0, time.UTC)},
		{"2024-06-14 First Quarter", time.Date(2024, 6, 14, 5, 18, 0, 0, time.UTC)},
		{"2024-06-28 Last Quarter", time.Date(2024, 6, 28, 21, 53, 0, 0, time.UTC)},
		{"2025-02-05 First Quarter", time.Date(2025, 2, 5, 8, 2, 0, 0, time.UTC)},
		{"2025-02-20 Last Quarter", time.Date(2025, 2, 20, 17, 32, 0, 0, time.UTC)},
	}
	for _, ev := range events {
		t.Run(ev.name, func(t *testing.T) {
			illum := GetCurrentMoonIllumination(ev.when, time.UTC)
			if illum < 0.45 || illum > 0.55 {
				t.Errorf("Quarter %v: illumination = %.4f, expected ~0.50", ev.when, illum)
			}
		})
	}
}

// Monotonicity check: between exact New Moon and the next Full Moon,
// illumination must strictly grow.
func TestIllumination_MonotonicGrowthNewToFull(t *testing.T) {
	newMoon := time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC)
	fullMoon := time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC)

	prev := GetCurrentMoonIllumination(newMoon, time.UTC)
	step := time.Hour * 12
	for t1 := newMoon.Add(step); t1.Before(fullMoon); t1 = t1.Add(step) {
		cur := GetCurrentMoonIllumination(t1, time.UTC)
		if cur < prev-0.001 { // tiny float tolerance
			t.Errorf("not monotonic at %v: prev=%.4f cur=%.4f", t1, prev, cur)
		}
		prev = cur
	}
}

// Monotonicity check: between Full Moon and the next New Moon, illumination
// must strictly decrease.
func TestIllumination_MonotonicDecayFullToNew(t *testing.T) {
	fullMoon := time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC)
	newMoon := time.Date(2024, 2, 9, 22, 59, 0, 0, time.UTC)

	prev := GetCurrentMoonIllumination(fullMoon, time.UTC)
	step := time.Hour * 12
	for t1 := fullMoon.Add(step); t1.Before(newMoon); t1 = t1.Add(step) {
		cur := GetCurrentMoonIllumination(t1, time.UTC)
		if cur > prev+0.001 {
			t.Errorf("not monotonic at %v: prev=%.4f cur=%.4f", t1, prev, cur)
		}
		prev = cur
	}
}

func TestLocOffsetSeconds_NilFallsBackToZero(t *testing.T) {
	got := locOffsetSeconds(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if got != 0 {
		t.Errorf("nil loc: got %d, want 0", got)
	}
}

// BinarySearchIllumination must converge to a JD near a known quarter event.
func TestBinarySearchIllumination_Quarters(t *testing.T) {
	newMoon := time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC)
	fullMoon := time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC)
	wantQ1 := time.Date(2024, 1, 18, 3, 53, 0, 0, time.UTC)

	jdNew := jt.ToJulianDate(newMoon)
	jdFull := jt.ToJulianDate(fullMoon)

	jdQ1 := BinarySearchIllumination(jdNew, jdFull, time.UTC, true)
	gotQ1 := jt.FromJulianDate(jdQ1, time.UTC)

	diff := gotQ1.Sub(wantQ1)
	if diff < 0 {
		diff = -diff
	}
	if diff > 6*time.Hour {
		t.Errorf("first quarter: got %v, want ~%v, diff %v", gotQ1, wantQ1, diff)
	}

	illum := GetCurrentMoonIllumination(gotQ1, time.UTC)
	if math.Abs(illum-0.5) > 0.05 {
		t.Errorf("first quarter illumination = %v, expected ~0.5", illum)
	}
}

func TestBinarySearchIllumination_NilLocFallsBack(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked with nil loc: %v", r)
		}
	}()
	jdStart := jt.ToJulianDate(time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC))
	jdEnd := jt.ToJulianDate(time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC))
	_ = BinarySearchIllumination(jdStart, jdEnd, nil, true)
}

// Astana, Kazakhstan (lat 51.13°N lon 71.43°E, UTC+5, no DST).
// timeanddate.com/moon/kazakstan/astana lists illumination at local midnight.
// We test only points where the value is unambiguous: dates where a New Moon
// or Full Moon falls inside the local day (so midnight value is dominated by
// the underlying conjunction/opposition).
func TestIllumination_DailyAstanaFullAndNew(t *testing.T) {
	astana := time.FixedZone("Astana", 5*3600)

	cases := []struct {
		date    time.Time
		minPct  float64
		maxPct  float64
		comment string
	}{
		// 2024-01-25 22:54 Astana = Full Moon. Local midnight is ~23h before.
		{time.Date(2024, 1, 25, 0, 0, 0, 0, astana), 97, 100, "Wolf Moon 2024 day"},
		// 2024-02-09 22:59 UTC = Feb 10 03:59 Astana — New Moon falls early next day.
		// midnight Feb 10 Astana is right after new moon → near 0.
		{time.Date(2024, 2, 10, 0, 0, 0, 0, astana), 0, 3, "New Moon 2024-02-10 Astana"},
		// 2025-01-14 03:27 Astana = Full Moon. midnight Jan 14 ~3h before.
		{time.Date(2025, 1, 14, 0, 0, 0, 0, astana), 99, 100, "Wolf Moon 2025"},
		// 2025-01-29 17:36 Astana = New Moon. midnight Jan 29 ~17h before.
		{time.Date(2025, 1, 29, 0, 0, 0, 0, astana), 0, 3, "New Moon 2025"},
		// 2023-01-07 04:08 Astana = Full Moon (Jan 6 23:08 UTC).
		{time.Date(2023, 1, 7, 0, 0, 0, 0, astana), 99, 100, "Wolf Moon 2023"},
	}

	for _, c := range cases {
		t.Run(c.comment, func(t *testing.T) {
			got := GetDailyMoonIllumination(c.date, astana) * 100
			if got < c.minPct || got > c.maxPct {
				t.Errorf("Astana %v: got %.2f%%, expected [%.0f,%.0f]%% (%s)",
					c.date, got, c.minPct, c.maxPct, c.comment)
			}
		})
	}
}
