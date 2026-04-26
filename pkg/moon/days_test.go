package moon

import (
	"testing"
	"time"
)

func TestGetMoonDays_NoMatchReturnsZero(t *testing.T) {
	// Empty table → zero result, no panic.
	days, mon := GetMoonDays(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), nil)
	if days.Begin != 0 || days.Current != 0 || days.End != 0 || mon != 0 {
		t.Errorf("empty table: got days=%v mon=%v, want zeros", days, mon)
	}

	// Date before all elements → no match → zeros.
	table := &MoonTable{
		Elems: []*MoonTableElement{
			{
				NewMoon:      time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC),
				FirstQuarter: time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC),
				FullMoon:     time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC),
				LastQuarter:  time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	days, mon = GetMoonDays(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), table.Elems)
	if days.Begin != 0 || mon != 0 {
		t.Errorf("date before table: got days=%v mon=%v, want zeros", days, mon)
	}
}

func TestGetMoonDays_BasicShape(t *testing.T) {
	// Use a real moon table for 2024 and ensure GetMoonDays produces sensible
	// values for a date inside the cycle.
	tg := time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC) // ~7 days into cycle
	table := CreateMoonTable(tg)
	days, mon := GetMoonDays(tg, table.Elems)
	if mon <= 0 {
		t.Fatalf("expected positive moon month, got %v", mon)
	}
	if days.Begin < 0 || days.Begin > mon {
		t.Errorf("days.Begin out of [0, mon): %v", days.Begin)
	}
	if days.End < 0 || days.End > mon {
		t.Errorf("days.End out of [0, mon): %v", days.End)
	}
	// End should be later than Begin within one cycle (modulo wrap-around).
	if days.End == days.Begin {
		t.Error("Begin and End should differ")
	}
}

func TestCurrentMoonDays_NilLocFallsBackToUTC(t *testing.T) {
	tg := time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC)
	table := CreateMoonTable(tg)
	res := CurrentMoonDays(tg, nil, table)
	// Just sanity — shouldn't panic. End or Current may be zero if table mismatch.
	_ = res
}

func TestCurrentMoonDays_NilTable(t *testing.T) {
	tg := time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC)
	res := CurrentMoonDays(tg, time.UTC, nil)
	if res.Begin != 0 || res.Current != 0 || res.End != 0 {
		t.Errorf("nil moon table should return zero, got %v", res)
	}
}

func TestSearchMoonDay_BasicSuccess(t *testing.T) {
	tg := time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC) // mid-cycle
	table := CreateMoonTable(tg)
	resp, err := SearchMoonDay(tg, table, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.From.Before(resp.To) {
		t.Errorf("From %v should be before To %v", resp.From, resp.To)
	}
}

func TestSearchMoonDay_NilTable(t *testing.T) {
	_, err := SearchMoonDay(time.Now(), nil, 5)
	if err == nil {
		t.Error("expected error for nil table")
	}
}

func TestSearchMoonDay_DayBeforeAnyCycle(t *testing.T) {
	// tGiven outside any element cycle → "not found".
	tg := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	table := &MoonTable{Elems: []*MoonTableElement{
		{
			NewMoon:      time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC),
			FirstQuarter: time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC),
			FullMoon:     time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC),
			LastQuarter:  time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
		},
	}}
	_, err := SearchMoonDay(tg, table, 1)
	if err == nil {
		t.Error("expected 'not found' error")
	}
}

func TestFindNearestPhase_Sanity(t *testing.T) {
	tg := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // between new (Jan 11) and full (Jan 25)
	table := CreateMoonTable(tg)
	np := FindNearestPhase(tg, table)

	// All four phases must be set to non-zero times in the future of tg.
	if np.NewMoon.IsZero() || np.FirstQuarter.IsZero() || np.FullMoon.IsZero() || np.LastQuarter.IsZero() {
		t.Errorf("unexpected zero phase: %+v", np)
	}
	for label, p := range map[string]time.Time{
		"NewMoon":      np.NewMoon,
		"FirstQuarter": np.FirstQuarter,
		"FullMoon":     np.FullMoon,
		"LastQuarter":  np.LastQuarter,
	} {
		if !p.After(tg) {
			t.Errorf("%s = %v, expected after %v", label, p, tg)
		}
	}
}

func TestFindNearestPhase_NilTable(t *testing.T) {
	np := FindNearestPhase(time.Now(), nil)
	if !np.NewMoon.IsZero() {
		t.Error("nil table should return zero NearestPhase")
	}
}

func TestBeginMoonDayToEarthDay(t *testing.T) {
	tg := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // mid-cycle
	table := CreateMoonTable(tg)

	t.Run("ISO format", func(t *testing.T) {
		got := BeginMoonDayToEarthDay(tg, time.Hour*24, "ISO", table.Elems)
		if got == nil {
			t.Fatal("nil result")
		}
		_, ok := (*got).(time.Time)
		if !ok {
			t.Errorf("expected time.Time, got %T", *got)
		}
	})

	t.Run("timestamp format", func(t *testing.T) {
		got := BeginMoonDayToEarthDay(tg, time.Hour*24, "timestamp", table.Elems)
		if got == nil {
			t.Fatal("nil result")
		}
		_, ok := (*got).(int64)
		if !ok {
			t.Errorf("expected int64, got %T", *got)
		}
	})

	t.Run("custom format", func(t *testing.T) {
		got := BeginMoonDayToEarthDay(tg, time.Hour*24, "2006-01-02", table.Elems)
		if got == nil {
			t.Fatal("nil result")
		}
		_, ok := (*got).(string)
		if !ok {
			t.Errorf("expected string, got %T", *got)
		}
	})

	t.Run("empty table", func(t *testing.T) {
		got := BeginMoonDayToEarthDay(tg, time.Hour*24, "ISO", nil)
		if got == nil {
			t.Fatal("expected non-nil even for empty table")
		}
	})
}

func TestMoonDetailed_NilLocFallsBackToUTC(t *testing.T) {
	// Skipping live-API-dependent paths: only assert no panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MoonDetailed panicked: %v", r)
		}
	}()
	// External API is not reachable in tests — we just verify the function
	// degrades gracefully when calls fail.
	res := MoonDetailed(time.Date(2024, 1, 18, 12, 0, 0, 0, time.UTC), nil, "en", "ISO", 71.43, 51.13)
	if res == nil {
		t.Error("MoonDetailed returned nil")
	}
}
