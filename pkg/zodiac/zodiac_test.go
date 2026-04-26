package zodiac

import (
	"moon/pkg/moon"
	"testing"
	"time"
)

func TestGetZodiacResp(t *testing.T) {
	tests := []struct {
		pos       int
		wantName  string
		wantEmoji string
	}{
		{0, "Virgo", "♍"},
		{1, "Libra", "♎"},
		{6, "Pisces", "♓"},
		{11, "Leo", "♌"},
	}
	for _, tt := range tests {
		name, emoji := getZodiacResp(tt.pos)
		if name != tt.wantName {
			t.Errorf("pos %d: name = %q, want %q", tt.pos, name, tt.wantName)
		}
		if emoji != tt.wantEmoji {
			t.Errorf("pos %d: emoji = %q, want %q", tt.pos, emoji, tt.wantEmoji)
		}
	}
}

func TestGetZodiacResp_OutOfRange(t *testing.T) {
	for _, pos := range []int{-1, 12, 100, -100} {
		name, emoji := getZodiacResp(pos)
		if name != "" || emoji != "" {
			t.Errorf("pos %d: expected empty, got %q %q", pos, name, emoji)
		}
	}
}

func TestGetZodiacRespLocalized(t *testing.T) {
	tests := []struct {
		pos  int
		lang string
		want string
	}{
		{0, "en", "Virgo"},
		{0, "ru", "Дева"},
		{0, "es", "Virgo"},
		{0, "de", "Jungfrau"},
		{0, "fr", "Vierge"},
		{0, "jp", "おとめ座"},
		{0, "xx", "Virgo"}, // unknown lang → fallback en
		{11, "ru", "Лев"},
		{6, "fr", "Poissons"},
	}
	for _, tt := range tests {
		got := getZodiacRespLocalized(tt.pos, tt.lang)
		if got != tt.want {
			t.Errorf("pos %d lang %q: got %q, want %q", tt.pos, tt.lang, got, tt.want)
		}
	}
}

func TestGetZodiacRespLocalized_OutOfRange(t *testing.T) {
	langs := []string{"en", "ru", "es", "de", "fr", "jp", "unknown"}
	for _, lang := range langs {
		for _, pos := range []int{-1, 12, 100} {
			got := getZodiacRespLocalized(pos, lang)
			if got != "" {
				t.Errorf("lang %s pos %d: expected empty, got %q", lang, pos, got)
			}
		}
	}
}

// CurrentZodiacs must not panic when loc is nil.
func TestCurrentZodiacs_NilLocDoesNotPanic(t *testing.T) {
	tGiven := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	moonTable := moon.CreateMoonTable(tGiven)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CurrentZodiacs panicked with nil loc: %v", r)
		}
	}()
	CurrentZodiacs(tGiven, nil, "en", "ISO", moonTable.Elems)
}

// Sanity check: zodiac sign names returned for a real moonTable are non-empty.
func TestCurrentZodiacs_ReturnsValidNames(t *testing.T) {
	tGiven := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	moonTable := moon.CreateMoonTable(tGiven)

	zods, zBegin, zCur, zEnd := CurrentZodiacs(tGiven, time.UTC, "en", "ISO", moonTable.Elems)

	if zods == nil {
		t.Fatal("zods is nil")
	}
	if zods.Count < 1 || zods.Count > 2 {
		t.Errorf("Count = %d, want 1 or 2", zods.Count)
	}
	if len(zods.Zodiac) != zods.Count {
		t.Errorf("len(Zodiac) = %d, Count = %d", len(zods.Zodiac), zods.Count)
	}
	for i, z := range zods.Zodiac {
		if z.Name == "" {
			t.Errorf("Zodiac[%d].Name is empty", i)
		}
	}
	if zBegin.Name == "" || zCur.Name == "" || zEnd.Name == "" {
		t.Errorf("zodiac begin/current/end name should be non-empty: %q/%q/%q", zBegin.Name, zCur.Name, zEnd.Name)
	}
}

func TestCurrentZodiacs_Localization(t *testing.T) {
	tGiven := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	moonTable := moon.CreateMoonTable(tGiven)

	for _, lang := range []string{"en", "ru", "es", "de", "fr", "jp"} {
		_, zBegin, _, _ := CurrentZodiacs(tGiven, time.UTC, lang, "ISO", moonTable.Elems)
		if zBegin.NameLocalized == "" {
			t.Errorf("lang %q: NameLocalized is empty", lang)
		}
	}
}
