package moon

import (
	"math"
	"testing"
	"time"
)

// External reference data for moon phase events, sourced from
//   - timeanddate.com/moon/phases (canonical reference)
//   - NASA Daily Moon Guide (science.nasa.gov/moon/daily-moon-guide)
//
// The Meeus algorithm used in Truephase has typical accuracy ~30 min for
// 20th–21st century dates, so we allow a 2-hour tolerance.
const phaseTolerance = 2 * time.Hour

type phaseEvent struct {
	when time.Time
	kind EnumPhase
	name string // human readable for diagnostics
}

// findClosestPhase returns the closest stored phase of `kind` to `target`.
// Returns the time.Time and a flag indicating whether anything was found.
func findClosestPhase(table *MoonTable, target time.Time, kind EnumPhase) (time.Time, bool) {
	if table == nil {
		return time.Time{}, false
	}
	var best time.Time
	bestDiff := time.Duration(math.MaxInt64)
	for _, e := range table.Elems {
		var candidate time.Time
		switch kind {
		case NewMoon:
			candidate = e.NewMoon
		case FirstQuarter:
			candidate = e.FirstQuarter
		case FullMoon:
			candidate = e.FullMoon
		case LastQuarter:
			candidate = e.LastQuarter
		}
		d := candidate.Sub(target)
		if d < 0 {
			d = -d
		}
		if d < bestDiff {
			bestDiff = d
			best = candidate
		}
	}
	return best, bestDiff != time.Duration(math.MaxInt64)
}

func TestCreateMoonTable_PhasesAgainstReferences(t *testing.T) {
	// All times are UTC, sourced from timeanddate.com/moon/phases.
	events := []phaseEvent{
		// 2023
		{time.Date(2023, 1, 6, 23, 8, 0, 0, time.UTC), FullMoon, "2023-01 Wolf Moon"},
		{time.Date(2023, 1, 21, 20, 53, 0, 0, time.UTC), NewMoon, "2023-01 New Moon"},
		{time.Date(2023, 2, 5, 18, 29, 0, 0, time.UTC), FullMoon, "2023-02 Snow Moon"},
		{time.Date(2023, 2, 20, 7, 6, 0, 0, time.UTC), NewMoon, "2023-02 New Moon"},
		{time.Date(2023, 3, 7, 12, 40, 0, 0, time.UTC), FullMoon, "2023-03 Worm Moon"},
		{time.Date(2023, 3, 21, 17, 23, 0, 0, time.UTC), NewMoon, "2023-03 New Moon"},
		{time.Date(2023, 7, 3, 11, 39, 0, 0, time.UTC), FullMoon, "2023-07 Buck Moon"},
		{time.Date(2023, 8, 1, 18, 31, 0, 0, time.UTC), FullMoon, "2023-08 Sturgeon Moon"},
		{time.Date(2023, 12, 27, 0, 33, 0, 0, time.UTC), FullMoon, "2023-12 Cold Moon"},

		// 2024
		{time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC), NewMoon, "2024-01 New Moon"},
		{time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC), FullMoon, "2024-01 Wolf Moon"},
		{time.Date(2024, 2, 9, 22, 59, 0, 0, time.UTC), NewMoon, "2024-02 New Moon"},
		{time.Date(2024, 2, 24, 12, 30, 0, 0, time.UTC), FullMoon, "2024-02 Snow Moon"},
		{time.Date(2024, 4, 8, 18, 21, 0, 0, time.UTC), NewMoon, "2024-04 Total Solar Eclipse"},
		{time.Date(2024, 4, 23, 23, 48, 0, 0, time.UTC), FullMoon, "2024-04 Pink Moon"},
		{time.Date(2024, 8, 19, 18, 26, 0, 0, time.UTC), FullMoon, "2024-08 Sturgeon Moon"},
		{time.Date(2024, 9, 18, 2, 34, 0, 0, time.UTC), FullMoon, "2024-09 Harvest Moon"},

		// 2025
		{time.Date(2025, 1, 13, 22, 27, 0, 0, time.UTC), FullMoon, "2025-01 Wolf Moon"},
		{time.Date(2025, 1, 29, 12, 36, 0, 0, time.UTC), NewMoon, "2025-01 New Moon"},
		{time.Date(2025, 2, 12, 13, 53, 0, 0, time.UTC), FullMoon, "2025-02 Snow Moon"},
		{time.Date(2025, 3, 14, 6, 55, 0, 0, time.UTC), FullMoon, "2025-03 Lunar Eclipse"},
		{time.Date(2025, 3, 29, 10, 58, 0, 0, time.UTC), NewMoon, "2025-03 Partial Solar Eclipse"},
	}

	for _, ev := range events {
		t.Run(ev.name, func(t *testing.T) {
			table := CreateMoonTable(ev.when)
			got, ok := findClosestPhase(table, ev.when, ev.kind)
			if !ok {
				t.Fatalf("CreateMoonTable produced no usable phase for %v", ev.when)
			}
			diff := got.Sub(ev.when)
			if diff < 0 {
				diff = -diff
			}
			if diff > phaseTolerance {
				t.Errorf("%s: predicted %v, expected ~%v, diff %v > %v",
					ev.name, got.UTC(), ev.when, diff, phaseTolerance)
			}
		})
	}
}

func TestCreateMoonTable_QuarterPhases(t *testing.T) {
	// First/Last quarter validation. timeanddate.com data.
	// Quarters use BinarySearchIllumination(50%); tolerance is wider than
	// Truephase (binary search step size + boundary effects).
	const quarterTol = 3 * time.Hour

	events := []phaseEvent{
		{time.Date(2024, 1, 18, 3, 53, 0, 0, time.UTC), FirstQuarter, "2024-01 first quarter"},
		{time.Date(2024, 2, 2, 23, 18, 0, 0, time.UTC), LastQuarter, "2024-02 last quarter"},
		{time.Date(2024, 6, 14, 5, 18, 0, 0, time.UTC), FirstQuarter, "2024-06 first quarter"},
		{time.Date(2024, 6, 28, 21, 53, 0, 0, time.UTC), LastQuarter, "2024-06 last quarter"},
		{time.Date(2025, 2, 5, 8, 2, 0, 0, time.UTC), FirstQuarter, "2025-02 first quarter"},
		{time.Date(2025, 2, 20, 17, 32, 0, 0, time.UTC), LastQuarter, "2025-02 last quarter"},
	}

	for _, ev := range events {
		t.Run(ev.name, func(t *testing.T) {
			table := CreateMoonTable(ev.when)
			got, ok := findClosestPhase(table, ev.when, ev.kind)
			if !ok {
				t.Fatalf("no phase found")
			}
			diff := got.Sub(ev.when)
			if diff < 0 {
				diff = -diff
			}
			if diff > quarterTol {
				t.Errorf("%s: predicted %v, expected ~%v, diff %v > %v",
					ev.name, got.UTC(), ev.when, diff, quarterTol)
			}
		})
	}
}

func TestCreateMoonTable_StructuralInvariants(t *testing.T) {
	for year := 2023; year <= 2026; year++ {
		t.Run(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"), func(t *testing.T) {
			tg := time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC)
			table := CreateMoonTable(tg)
			if table == nil || len(table.Elems) == 0 {
				t.Fatal("empty table")
			}
			// Each element: NewMoon < FirstQuarter < FullMoon < LastQuarter, and
			// NewMoons are strictly increasing across elements.
			for i, e := range table.Elems {
				if !e.NewMoon.Before(e.FirstQuarter) {
					t.Errorf("elem %d: NewMoon %v not before FirstQuarter %v", i, e.NewMoon, e.FirstQuarter)
				}
				if !e.FirstQuarter.Before(e.FullMoon) {
					t.Errorf("elem %d: FirstQuarter %v not before FullMoon %v", i, e.FirstQuarter, e.FullMoon)
				}
				if !e.FullMoon.Before(e.LastQuarter) {
					t.Errorf("elem %d: FullMoon %v not before LastQuarter %v", i, e.FullMoon, e.LastQuarter)
				}
				if i > 0 && !table.Elems[i-1].NewMoon.Before(e.NewMoon) {
					t.Errorf("elem %d: NewMoon not increasing", i)
				}
			}
			// Synodic month is ~29.53 days. There should be 12-14 NewMoons per year.
			count := len(table.Elems)
			if count < 11 || count > 15 {
				t.Errorf("unexpected count of phases: %d", count)
			}
		})
	}
}
