package position

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// withMockUpstream replaces baseURL with an httptest server for the duration
// of the test, restoring the original on cleanup. It returns the server so
// the caller can install handlers.
func withMockUpstream(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	prev := baseURL
	baseURL = srv.URL + "/"
	t.Cleanup(func() {
		baseURL = prev
		srv.Close()
	})
	return srv
}

func TestGetRisesDay_Success(t *testing.T) {
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/daily") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Echo back the parameters we expect from the Go layer.
		if r.URL.Query().Get("lat") == "" || r.URL.Query().Get("lon") == "" {
			t.Errorf("missing lat/lon: %v", r.URL.Query())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"Status": "success",
			"Range": "single_day",
			"Parameters": {"Latitude": 51.13, "Longitude": 71.43, "Year": 2024, "Month": 1, "Day": 15},
			"Data": {
				"IsMoonRise": true,
				"IsMoonSet": true,
				"IsMeridian": true,
				"Moonrise": {"Timestamp": 1705320000, "AzimuthDegrees": 90.0, "AltitudeDegrees": 0.0, "Direction": "E", "DistanceKm": 384400.0},
				"Moonset":  {"Timestamp": 1705370000, "AzimuthDegrees": 270.0, "AltitudeDegrees": 0.0, "Direction": "W", "DistanceKm": 384500.0},
				"Meridian": {"Timestamp": 1705345000, "AzimuthDegrees": 180.0, "AltitudeDegrees": 60.0, "Direction": "S", "DistanceKm": 384450.0}
			}
		}`))
	})

	d, err := GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.IsMoonRise || !d.IsMoonSet || !d.IsMeridian {
		t.Error("expected all flags true")
	}
	if d.Moonrise.Time == nil {
		t.Error("Moonrise.Time should be set")
	}
	if d.Moonrise.Timestamp != nil {
		t.Error("Moonrise.Timestamp should be cleared after conversion")
	}
}

// Regression for the panic when upstream returns Data: null.
func TestGetRisesDay_NullDataReturnsError(t *testing.T) {
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Status":"success","Data":null}`))
	})

	d, err := GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
	if d != nil {
		t.Error("expected nil DayData")
	}
	if err == nil {
		t.Fatal("expected error for null Data")
	}
}

// Regression for the panic in `fmt.Errorf("[%s]...", resp.Status, err)` when
// the request itself fails (resp is nil).
func TestGetRisesDay_NetworkErrorReturnsError(t *testing.T) {
	prev := baseURL
	baseURL = "http://127.0.0.1:1/" // refused
	t.Cleanup(func() { baseURL = prev })

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked instead of returning error: %v", r)
		}
	}()
	d, err := GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
	if d != nil {
		t.Error("expected nil DayData on network error")
	}
	if err == nil {
		t.Fatal("expected error on connection refused")
	}
}

func TestGetRisesDay_NonOKStatusReturnsError(t *testing.T) {
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"Status":"error","Message":"month must be between 1 and 12"}`))
	})

	d, err := GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
	if d != nil {
		t.Error("expected nil DayData")
	}
	if err == nil {
		t.Fatal("expected error for 400")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("error should mention status: %v", err)
	}
}

func TestCacheGetRisesDay_CachesResults(t *testing.T) {
	var hits int64
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Status":"success","Data":{"IsMoonRise":true,"IsMoonSet":true,"IsMeridian":true}}`))
	})

	c := &Cache{}
	for i := 0; i < 5; i++ {
		_, err := c.GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Errorf("expected 1 upstream hit, got %d", got)
	}
}

func TestCacheGetRisesMonthly_CachesResults(t *testing.T) {
	var hits int64
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"Status":"success",
			"Range":"Full_month",
			"DaysCount":2,
			"Parameters":{"Latitude":51.13,"Longitude":71.43,"Year":2024,"Month":1},
			"Data":[
				{"IsMoonRise":true,"IsMoonSet":true,"IsMeridian":true},
				{"IsMoonRise":true,"IsMoonSet":true,"IsMeridian":false}
			]
		}`))
	})

	c := &Cache{}
	for i := 0; i < 5; i++ {
		out, err := c.GetRisesMonthly(2024, 1, time.UTC, 2, "ISO", 71.43, 51.13)
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if out == nil || len(*out) != 2 {
			t.Errorf("call %d: unexpected data", i)
		}
	}
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Errorf("expected 1 upstream hit, got %d", got)
	}
}

func TestCacheGetRisesDay_NoLocationReturnsError(t *testing.T) {
	c := &Cache{}
	_, err := c.GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO")
	if err == nil {
		t.Error("expected error for missing location")
	}
}

func TestGetMoonPosition_Success(t *testing.T) {
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"Status":"success",
			"Timestamp":1705345000,
			"AzimuthDegrees":180.0,
			"AltitudeDegrees":60.0,
			"Direction":"S",
			"DistanceKm":384450.0
		}`))
	})

	p, err := GetMoonPosition(time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), time.UTC, 2, "ISO", 71.43, 51.13)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("nil position")
	}
	if p.Time == nil {
		t.Error("Time should be populated")
	}
	if p.Timestamp != nil {
		t.Error("Timestamp should be cleared after conversion")
	}
	if p.Direction != "S" {
		t.Errorf("Direction = %q, want S", p.Direction)
	}
}

func TestGetMoonPosition_NullResponseHandled(t *testing.T) {
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	p, err := GetMoonPosition(time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), time.UTC, 2, "ISO", 71.43, 51.13)
	// Either nil position or graceful nil — but no panic.
	_ = p
	_ = err
}

func TestGetRisesDay_MoonsetWithoutTimestamp(t *testing.T) {
	// Regression: Moonset != nil but Moonset.Timestamp == nil must not panic
	// and must not set garbage in Time.
	withMockUpstream(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"Status":"success",
			"Data": {
				"IsMoonRise": false,
				"IsMoonSet": true,
				"IsMeridian": false,
				"Moonset": {"AzimuthDegrees": 270.0, "AltitudeDegrees": 0.0, "Direction": "W", "DistanceKm": 384500.0}
			}
		}`))
	})
	d, err := GetRisesDay(2024, 1, 15, time.UTC, 2, "ISO", 71.43, 51.13)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Moonset == nil {
		t.Fatal("Moonset should be present")
	}
	if d.Moonset.Time != nil {
		t.Error("Moonset.Time should remain nil when Timestamp was nil")
	}
}
