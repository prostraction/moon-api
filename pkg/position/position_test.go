package position

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestParseLocation(t *testing.T) {
	t.Run("two values", func(t *testing.T) {
		// Callers pass (longitude, latitude); parseLocation returns lat, lon.
		lat, lon, err := parseLocation([]float64{71.43, 51.05}) // lon, lat for Astana
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lat != 51.05 || lon != 71.43 {
			t.Errorf("got lat=%v lon=%v", lat, lon)
		}
	})

	t.Run("empty", func(t *testing.T) {
		_, _, err := parseLocation(nil)
		if err == nil {
			t.Error("expected error for empty location")
		}
	})

	t.Run("one value", func(t *testing.T) {
		_, _, err := parseLocation([]float64{1})
		if err == nil {
			t.Error("expected error for single value")
		}
	})

	t.Run("three values", func(t *testing.T) {
		_, _, err := parseLocation([]float64{1, 2, 3})
		if err == nil {
			t.Error("expected error for >2 values")
		}
	})
}

func TestTimestampToGoTime(t *testing.T) {
	ts := int64(1735689600) // 2025-01-01 00:00:00 UTC

	t.Run("nil input", func(t *testing.T) {
		got := timestampToGoTime(nil, "ISO", time.UTC)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("timestamp format", func(t *testing.T) {
		got := timestampToGoTime(&ts, "timestamp", time.UTC)
		if got == nil {
			t.Fatal("nil result")
		}
		v, ok := (*got).(int64)
		if !ok {
			t.Fatalf("expected int64, got %T", *got)
		}
		if v != ts {
			t.Errorf("got %d, want %d", v, ts)
		}
	})

	t.Run("ISO format", func(t *testing.T) {
		got := timestampToGoTime(&ts, "ISO", time.UTC)
		if got == nil {
			t.Fatal("nil result")
		}
		_, ok := (*got).(time.Time)
		if !ok {
			t.Fatalf("expected time.Time, got %T", *got)
		}
	})

	t.Run("custom format", func(t *testing.T) {
		got := timestampToGoTime(&ts, "2006-01-02", time.UTC)
		if got == nil {
			t.Fatal("nil result")
		}
		s, ok := (*got).(string)
		if !ok {
			t.Fatalf("expected string, got %T", *got)
		}
		if s != "2025-01-01" {
			t.Errorf("got %q, want 2025-01-01", s)
		}
	})

	t.Run("nil location falls back to UTC", func(t *testing.T) {
		got := timestampToGoTime(&ts, "2006-01-02 15:04", nil)
		if got == nil {
			t.Fatal("nil result")
		}
		s := (*got).(string)
		if !strings.HasPrefix(s, "2025-01-01") {
			t.Errorf("got %q, want 2025-01-01 prefix", s)
		}
	})

	t.Run("timezone shifts wall clock", func(t *testing.T) {
		loc := time.FixedZone("UTC+5", 5*3600)
		got := timestampToGoTime(&ts, "2006-01-02 15:04", loc)
		s := (*got).(string)
		if s != "2025-01-01 05:00" {
			t.Errorf("got %q, want 2025-01-01 05:00", s)
		}
	})

	t.Run("case insensitive iso", func(t *testing.T) {
		got := timestampToGoTime(&ts, "iso", time.UTC)
		if got == nil {
			t.Fatal("nil result")
		}
		_, ok := (*got).(time.Time)
		if !ok {
			t.Fatalf("expected time.Time, got %T", *got)
		}
	})
}

// Cache must be safe for concurrent reads/writes — without a mutex this would
// reproduce "concurrent map writes" panic.
func TestCacheConcurrentAccess(t *testing.T) {
	c := &Cache{}
	const goroutines = 16
	const iterations = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				key := "k"
				c.mu.Lock()
				if c.CacheDaily == nil {
					c.CacheDaily = make(map[string]*DayData)
				}
				c.CacheDaily[key] = &DayData{}
				c.mu.Unlock()

				c.mu.RLock()
				_ = c.CacheDaily[key]
				c.mu.RUnlock()
			}
		}(g)
	}
	wg.Wait()
}
