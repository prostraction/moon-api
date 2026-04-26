package phase

import (
	"testing"
	"time"
)

// CurrentMoonPhase against external references — for known full/new moons,
// the returned phase name must match the expected category.
func TestCurrentMoonPhase_KnownEvents(t *testing.T) {
	cases := []struct {
		when         time.Time
		wantName     string
		wantWaxing   bool
		checkWaxing  bool
	}{
		// Around full moons → "Full Moon"
		{time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC), "Full Moon", false, false},
		{time.Date(2025, 1, 13, 22, 27, 0, 0, time.UTC), "Full Moon", false, false},
		{time.Date(2023, 1, 6, 23, 8, 0, 0, time.UTC), "Full Moon", false, false},

		// Around new moons → "New Moon"
		{time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC), "New Moon", false, false},
		{time.Date(2025, 1, 29, 12, 36, 0, 0, time.UTC), "New Moon", false, false},

		// Mid-cycle between new and full 2024 → must be waxing.
		{time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), "Waxing Crescent", true, true},
		// Three days after Full Moon Jan 25 → Waning Gibbous, waning.
		{time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC), "Waning Gibbous", false, true},
	}

	for _, c := range cases {
		t.Run(c.when.Format("2006-01-02"), func(t *testing.T) {
			pr := CurrentMoonPhase(c.when, "en")
			if pr == nil || pr.Current == nil {
				t.Fatal("nil phase response")
			}
			if pr.Current.Name != c.wantName {
				t.Errorf("at %v: phase = %q, want %q (illum=%.4f)",
					c.when, pr.Current.Name, c.wantName, pr.Illumination.Current)
			}
			if c.checkWaxing && pr.Current.IsWaxing != c.wantWaxing {
				t.Errorf("at %v: IsWaxing = %v, want %v", c.when, pr.Current.IsWaxing, c.wantWaxing)
			}
		})
	}
}

// CurrentMoonPhase must produce well-formed Phase responses (non-nil pointers,
// non-empty names, illumination in [0,1]).
func TestCurrentMoonPhase_StructuralIntegrity(t *testing.T) {
	tests := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		time.Date(2025, 7, 4, 6, 30, 0, 0, time.UTC),
	}
	for _, when := range tests {
		t.Run(when.Format("2006-01-02"), func(t *testing.T) {
			pr := CurrentMoonPhase(when, "en")
			if pr == nil {
				t.Fatal("nil response")
			}
			for label, p := range map[string]*Phase{
				"BeginDay": pr.BeginDay,
				"Current":  pr.Current,
				"EndDay":   pr.EndDay,
			} {
				if p == nil {
					t.Errorf("%s phase is nil", label)
					continue
				}
				if p.Name == "" {
					t.Errorf("%s phase name empty", label)
				}
				if p.Emoji == "" {
					t.Errorf("%s phase emoji empty", label)
				}
			}
			i := pr.Illumination
			for name, v := range map[string]float64{
				"BeginDay": i.BeginDay,
				"Current":  i.Current,
				"EndDay":   i.EndDay,
			} {
				if v < 0 || v > 1 {
					t.Errorf("%s illumination out of [0,1]: %v", name, v)
				}
			}
		})
	}
}

// CurrentMoonPhase must respect the lang parameter: localized name should
// differ between languages where possible.
func TestCurrentMoonPhase_Localization(t *testing.T) {
	when := time.Date(2024, 1, 25, 17, 54, 0, 0, time.UTC) // Full Moon

	pr := CurrentMoonPhase(when, "ru")
	if pr.Current.NameLocalized != "Полнолуние" {
		t.Errorf("ru: got %q, want Полнолуние", pr.Current.NameLocalized)
	}

	pr = CurrentMoonPhase(when, "es")
	if pr.Current.NameLocalized != "Luna llena" {
		t.Errorf("es: got %q, want Luna llena", pr.Current.NameLocalized)
	}

	pr = CurrentMoonPhase(when, "jp")
	if pr.Current.NameLocalized != "満月" {
		t.Errorf("jp: got %q, want 満月", pr.Current.NameLocalized)
	}

	pr = CurrentMoonPhase(when, "unknown")
	if pr.Current.NameLocalized != "Full Moon" {
		t.Errorf("unknown lang: got %q, want fallback to Full Moon", pr.Current.NameLocalized)
	}
}

// Truephase output must be a finite Julian day (sanity).
func TestTruephase_FiniteOutput(t *testing.T) {
	for _, k := range []float64{-1000, 0, 1000, 2000} {
		for _, p := range []float64{0, 0.25, 0.5, 0.75} {
			got := Truephase(k, p)
			if got <= 0 {
				t.Errorf("Truephase(%v,%v) = %v, expected positive JD", k, p, got)
			}
		}
	}
}
