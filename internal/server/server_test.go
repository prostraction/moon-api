package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"moon/pkg/position"
)

// newTestApp builds a Fiber app with all routes attached but without Prefork
// or other production-only knobs. Each test gets its own Server (and thus its
// own positionCache) for isolation.
func newTestApp(t *testing.T) *fiber.App {
	t.Helper()
	s := &Server{}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	s.RegisterRoutes(app)
	return app
}

// doGet sends a GET request through the in-memory test transport and decodes
// the JSON body. Returns the HTTP status and the raw body so individual tests
// can also inspect text.
func doGet(t *testing.T, app *fiber.App, path string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	resp, err := app.Test(req, -1) // -1 disables the test timeout
	if err != nil {
		t.Fatalf("app.Test(%s): %v", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, body
}

func TestVersionRoute(t *testing.T) {
	app := newTestApp(t)

	for _, path := range []string{"/v1/version", "/api/v1/version"} {
		status, body := doGet(t, app, path)
		if status != 200 {
			t.Errorf("%s: status %d, want 200", path, status)
		}
		var s string
		if err := json.Unmarshal(body, &s); err != nil {
			t.Errorf("%s: bad JSON: %v (%s)", path, err, body)
		}
		if s == "" {
			t.Errorf("%s: empty version string", path)
		}
	}
}

func TestMoonTableYearRoute(t *testing.T) {
	app := newTestApp(t)

	t.Run("ISO format", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonTableYear?year=2024")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		// Body is a JSON array of MoonTableElement.
		var arr []map[string]any
		if err := json.Unmarshal(body, &arr); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(arr) < 11 || len(arr) > 14 {
			t.Errorf("len=%d, want 11..14 phases per year", len(arr))
		}
		for i, e := range arr {
			for _, k := range []string{"NewMoon", "FirstQuarter", "FullMoon", "LastQuarter"} {
				if e[k] == nil {
					t.Errorf("elem %d: missing %s", i, k)
				}
			}
		}
	})

	t.Run("timestamp format", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonTableYear?year=2024&timeFormat=timestamp")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var arr []map[string]int64
		if err := json.Unmarshal(body, &arr); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// Sanity: LastQuarter must NOT equal FullMoon (the bug we fixed).
		for i, e := range arr {
			if e["LastQuarter"] == e["FullMoon"] {
				t.Errorf("elem %d: LastQuarter == FullMoon (regression)", i)
			}
			if e["NewMoon"] >= e["FirstQuarter"] {
				t.Errorf("elem %d: NewMoon not before FirstQuarter", i)
			}
		}
	})

	t.Run("custom format", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonTableYear?year=2024&timeFormat=2006-01-02")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var arr []map[string]string
		if err := json.Unmarshal(body, &arr); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(arr) == 0 {
			t.Fatal("empty result")
		}
		// LastQuarter format string must look like a date; in particular it
		// must NOT equal FullMoon (regression check).
		for i, e := range arr {
			if e["LastQuarter"] == e["FullMoon"] {
				t.Errorf("elem %d: LastQuarter == FullMoon (regression)", i)
			}
		}
	})
}

func TestMoonTableCurrentRoute(t *testing.T) {
	app := newTestApp(t)
	status, body := doGet(t, app, "/v1/moonTableCurrent")
	if status != 200 {
		t.Fatalf("status %d: %s", status, body)
	}
}

func TestNextMoonPhaseRoute(t *testing.T) {
	app := newTestApp(t)

	t.Run("ISO", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/nextMoonPhase")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var resp map[string]any
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		for _, k := range []string{"NewMoon", "FirstQuarter", "FullMoon", "LastQuarter"} {
			if resp[k] == nil {
				t.Errorf("missing %s", k)
			}
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/nextMoonPhase?timeFormat=timestamp")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var resp map[string]int64
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// All four should be future timestamps.
		for k, v := range resp {
			if v <= 0 {
				t.Errorf("%s = %d (expected > 0)", k, v)
			}
		}
	})
}

func TestNextMoonDayRoute(t *testing.T) {
	app := newTestApp(t)

	t.Run("missing day", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/nextMoonDay")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("non-int day", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/nextMoonDay?day=abc")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("day out of range", func(t *testing.T) {
		status, _ := doGet(t, app, "/v1/nextMoonDay?day=99")
		if status != 400 {
			t.Errorf("status %d, want 400", status)
		}
	})

	t.Run("valid day", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/nextMoonDay?day=5")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
	})
}

func TestJulianTimeRoutes(t *testing.T) {
	app := newTestApp(t)

	t.Run("toJulianTimeByDate valid", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/toJulianTimeByDate?year=2024&month=1&day=1")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var resp map[string]any
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		jd, ok := resp["JulianDate"].(float64)
		if !ok {
			t.Fatalf("JulianDate missing or wrong type: %v", resp)
		}
		// 2024-01-01 00:00 UTC is JD ≈ 2460310.5
		if jd < 2460000 || jd > 2461000 {
			t.Errorf("JD=%v outside expected range for 2024", jd)
		}
	})

	t.Run("toJulianTimeByDate invalid day in month", func(t *testing.T) {
		// Feb 29 in non-leap year — IsValidDate flags this even though
		// StrToInt accepts day=29 (within clamp [1,31]).
		status, body := doGet(t, app, "/v1/toJulianTimeByDate?year=2023&month=2&day=29")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("toJulianTimeByDate Apr 31", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/toJulianTimeByDate?year=2024&month=4&day=31")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("toJulianTimeByTimestamp", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/toJulianTimeByTimestamp?timestamp=1735689600")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
	})

	t.Run("fromJulianTime missing param", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/fromJulianTime")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("fromJulianTime valid", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/fromJulianTime?jtime=2460310.5")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
	})

	t.Run("round-trip", func(t *testing.T) {
		_, body := doGet(t, app, "/v1/toJulianTimeByDate?year=2024&month=6&day=15&hour=12&minute=0&second=0&precision=8")
		var first map[string]any
		_ = json.Unmarshal(body, &first)
		jd := first["JulianDate"].(float64)

		_, body2 := doGet(t, app, "/v1/fromJulianTime?jtime="+strings.TrimSpace(jsonNumber(jd)))
		var second map[string]any
		if err := json.Unmarshal(body2, &second); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if second["CivilDate"] == nil {
			t.Errorf("CivilDate missing in round-trip: %s", body2)
		}
	})
}

func TestMoonPhaseRoutes(t *testing.T) {
	app := newTestApp(t)

	t.Run("moonPhaseCurrent", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonPhaseCurrent")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		var resp map[string]any
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp["CurrentState"] == nil {
			t.Error("missing CurrentState")
		}
	})

	t.Run("moonPhaseDate Wolf Moon 2024", func(t *testing.T) {
		// 2024-01-25 is a Full Moon day.
		status, body := doGet(t, app, "/v1/moonPhaseDate?year=2024&month=1&day=25&hour=18")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		// Should report "Full Moon" as current phase.
		if !strings.Contains(string(body), "Full Moon") {
			t.Errorf("expected 'Full Moon' in response: %s", body)
		}
	})

	t.Run("moonPhaseDate New Moon 2024-01-11", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonPhaseDate?year=2024&month=1&day=11&hour=12")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		if !strings.Contains(string(body), "New Moon") {
			t.Errorf("expected 'New Moon' in response: %s", body)
		}
	})

	t.Run("moonPhaseDate invalid date", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonPhaseDate?year=2023&month=2&day=29")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("moonPhaseTimestamp", func(t *testing.T) {
		// 2025-01-13 22:27 UTC = Wolf Moon 2025
		status, body := doGet(t, app, "/v1/moonPhaseTimestamp?timestamp=1736807220")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
	})

	t.Run("moonPhase localized ru", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonPhaseDate?year=2024&month=1&day=25&hour=18&lang=ru")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
		if !strings.Contains(string(body), "Полнолуние") {
			t.Errorf("expected 'Полнолуние' in ru response: %s", body)
		}
	})
}

func TestMoonPositionMonthlyRoute(t *testing.T) {
	app := newTestApp(t)

	t.Run("missing coords", func(t *testing.T) {
		status, body := doGet(t, app, "/v1/moonPositionMonthly?year=2024&month=1")
		if status != 400 {
			t.Errorf("status %d, want 400. body=%s", status, body)
		}
	})

	t.Run("with coords + mocked upstream", func(t *testing.T) {
		mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"Status":"success",
				"Range":"Full_month",
				"DaysCount":1,
				"Parameters":{"Latitude":51.13,"Longitude":71.43,"Year":2024,"Month":1},
				"Data":[
					{"IsMoonRise":true,"IsMoonSet":true,"IsMeridian":true}
				]
			}`))
		}))
		defer mock.Close()
		restore := position.SetBaseURLForTesting(mock.URL + "/")
		defer restore()

		status, body := doGet(t, app, "/v1/moonPositionMonthly?year=2024&month=1&latitude=51.13&longitude=71.43")
		if status != 200 {
			t.Fatalf("status %d: %s", status, body)
		}
	})

	t.Run("upstream 400 propagates as 400", func(t *testing.T) {
		mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Status":"error","Message":"month must be between 1 and 12"}`))
		}))
		defer mock.Close()
		restore := position.SetBaseURLForTesting(mock.URL + "/")
		defer restore()

		// Use different coords to avoid hitting the previous sub-test's cache.
		status, _ := doGet(t, app, "/v1/moonPositionMonthly?year=2024&month=2&latitude=52.52&longitude=13.4")
		if status != 400 {
			t.Errorf("status %d, want 400", status)
		}
	})
}

// helper: format a float as a plain JSON number (no exponential notation).
func jsonNumber(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}
