package moon

import (
	"errors"
	"testing"
	"time"
)

func TestSearchPhase_NextPhase(t *testing.T) {
	location := time.UTC

	tests := []struct {
		name      string
		tGiven    time.Time
		moonTable *MoonTable
		phase     EnumPhase
		wantTime  time.Time
		wantError error
	}{
		{
			name:      "nil moon table returns error",
			tGiven:    time.Now(),
			moonTable: nil,
			phase:     NewMoon,
			wantError: errors.New("passed empty moonTable to SearchNewMoon"),
		},
		{
			name:   "empty moon table returns not found",
			tGiven: time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{},
			},
			phase:     NewMoon,
			wantError: errors.New("not found"),
		},
		{
			name:   "next new moon after current time - found in current element",
			tGiven: time.Date(2023, 1, 5, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
		},
		{
			name:   "next new moon when current phase already passed - get from next element",
			tGiven: time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    NewMoon,
			wantTime: time.Date(2023, 2, 1, 0, 0, 0, 0, location),
		},
		{
			name:   "next full moon between phases - found in next element",
			tGiven: time.Date(2023, 1, 25, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    FullMoon,
			wantTime: time.Date(2023, 2, 15, 0, 0, 0, 0, location),
		},
		{
			name:   "next phase when time is exactly on phase - should return next occurrence",
			tGiven: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
		},
		{
			name:   "time before all phases in table - return first phase",
			tGiven: time.Date(2022, 12, 31, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
				},
			},
			phase:    NewMoon,
			wantTime: time.Date(2023, 1, 1, 0, 0, 0, 0, location),
		},
		{
			name:   "element with t1 == t2 should be skipped",
			tGiven: time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           1.0,
					},
				},
			},
			phase:     NewMoon,
			wantError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{}

			gotTime, err := cache.SearchPhase(tt.tGiven, tt.moonTable, tt.phase)
			if tt.wantError != nil {
				if err == nil {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				} else if err.Error() != tt.wantError.Error() {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Errorf("Test: %v, SearchPhase() unexpected error = %v", tt.name, err)
				return
			}

			if !gotTime.Equal(tt.wantTime) {
				t.Errorf("Test: %v, SearchPhase() gotTime = %v, want %v", tt.name, gotTime, tt.wantTime)
			}
		})
	}
}

func TestSearchPhase_AllPhaseTypes(t *testing.T) {
	location := time.UTC

	moonTable := &MoonTable{
		Elems: []*MoonTableElement{
			{
				NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
				FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
				FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
				LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
				t1:           1.0,
				t2:           2.0,
			},
			{
				NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
				FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
				FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
				LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
				t1:           3.0,
				t2:           4.0,
			},
		},
	}

	tests := []struct {
		name     string
		tGiven   time.Time
		phase    EnumPhase
		wantTime time.Time
	}{
		{
			name:     "next new moon",
			tGiven:   time.Date(2023, 1, 5, 0, 0, 0, 0, location),
			phase:    NewMoon,
			wantTime: time.Date(2023, 2, 1, 0, 0, 0, 0, location),
		},
		{
			name:     "next first quarter",
			tGiven:   time.Date(2023, 1, 5, 0, 0, 0, 0, location),
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
		},
		{
			name:     "next full moon",
			tGiven:   time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			phase:    FullMoon,
			wantTime: time.Date(2023, 1, 15, 0, 0, 0, 0, location),
		},
		{
			name:     "next last quarter",
			tGiven:   time.Date(2023, 1, 20, 0, 0, 0, 0, location),
			phase:    LastQuarter,
			wantTime: time.Date(2023, 1, 22, 0, 0, 0, 0, location),
		},
	}

	cache := &Cache{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cache.SearchPhase(tt.tGiven, moonTable, tt.phase)

			if err != nil {
				t.Errorf("SearchPhase() for %s failed with error: %v", tt.name, err)
				return
			}

			if !result.Equal(tt.wantTime) {
				t.Errorf("SearchPhase() for %s gotTime = %v, want %v", tt.name, result, tt.wantTime)
			}
		})
	}
}

func TestSearchPhase_ComplexScenarios(t *testing.T) {
	location := time.UTC

	tests := []struct {
		name      string
		tGiven    time.Time
		moonTable *MoonTable
		phase     EnumPhase
		wantTime  time.Time
	}{
		{
			name:   "multiple elements - find in third element",
			tGiven: time.Date(2023, 3, 1, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
					{
						NewMoon:      time.Date(2023, 3, 3, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 3, 10, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 3, 17, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 3, 24, 0, 0, 0, 0, location),
						t1:           5.0,
						t2:           6.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 3, 10, 0, 0, 0, 0, location),
		},
		{
			name:   "phase in gap between elements",
			tGiven: time.Date(2023, 1, 25, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    NewMoon,
			wantTime: time.Date(2023, 2, 1, 0, 0, 0, 0, location),
		},
	}

	cache := &Cache{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cache.SearchPhase(tt.tGiven, tt.moonTable, tt.phase)

			if err != nil {
				t.Errorf("SearchPhase() for %s failed with error: %v", tt.name, err)
				return
			}

			if !result.Equal(tt.wantTime) {
				t.Errorf("SearchPhase() for %s gotTime = %v, want %v", tt.name, result, tt.wantTime)
			}
		})
	}
}

func TestSearchPhase_EdgeCases(t *testing.T) {
	location := time.UTC

	tests := []struct {
		name      string
		tGiven    time.Time
		moonTable *MoonTable
		phase     EnumPhase
		wantTime  time.Time
		wantError error
	}{
		{
			name:   "time exactly at last quarter - should return next phase from next element",
			tGiven: time.Date(2023, 1, 22, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    LastQuarter,
			wantTime: time.Date(2023, 2, 22, 0, 0, 0, 0, location),
		},
		{
			name:   "time between last quarter and next new moon - should return next new moon",
			tGiven: time.Date(2023, 1, 25, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           2.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           3.0,
						t2:           4.0,
					},
				},
			},
			phase:    NewMoon,
			wantTime: time.Date(2023, 2, 1, 0, 0, 0, 0, location),
		},
		{
			name:   "all elements skipped due to t1 == t2",
			tGiven: time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           1.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           2.0,
						t2:           2.0,
					},
				},
			},
			phase:     NewMoon,
			wantError: errors.New("not found"),
		},
		{
			name:   "mixed valid and invalid elements",
			tGiven: time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
						t1:           1.0,
						t2:           1.0,
					},
					{
						NewMoon:      time.Date(2023, 2, 1, 0, 0, 0, 0, location),
						FirstQuarter: time.Date(2023, 2, 8, 0, 0, 0, 0, location),
						FullMoon:     time.Date(2023, 2, 15, 0, 0, 0, 0, location),
						LastQuarter:  time.Date(2023, 2, 22, 0, 0, 0, 0, location),
						t1:           2.0,
						t2:           3.0,
					},
				},
			},
			phase:    NewMoon,
			wantTime: time.Date(2023, 2, 1, 0, 0, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{}

			gotTime, err := cache.SearchPhase(tt.tGiven, tt.moonTable, tt.phase)

			if tt.wantError != nil {
				if err == nil {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				} else if err.Error() != tt.wantError.Error() {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Errorf("Test: %v, SearchPhase() unexpected error = %v", tt.name, err)
				return
			}

			if !gotTime.Equal(tt.wantTime) {
				t.Errorf("Test: %v, SearchPhase() gotTime = %v, want %v", tt.name, gotTime, tt.wantTime)
			}
		})
	}
}

func TestSearchPhase_TimeLocation(t *testing.T) {
	utcLocation := time.UTC
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	newYorkLocation, _ := time.LoadLocation("America/New_York")

	tests := []struct {
		name      string
		tGiven    time.Time
		moonTable *MoonTable
		phase     EnumPhase
		wantTime  time.Time
	}{
		{
			name:   "Moscow timezone",
			tGiven: time.Date(2023, 1, 5, 0, 0, 0, 0, moscowLocation),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, moscowLocation),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, moscowLocation),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, moscowLocation),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, moscowLocation),
						t1:           1.0,
						t2:           2.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, moscowLocation),
		},
		{
			name:   "New York timezone",
			tGiven: time.Date(2023, 1, 5, 0, 0, 0, 0, newYorkLocation),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, newYorkLocation),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, newYorkLocation),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, newYorkLocation),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, newYorkLocation),
						t1:           1.0,
						t2:           2.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, newYorkLocation),
		},
		{
			name:   "mixed timezones - should handle correctly",
			tGiven: time.Date(2023, 1, 5, 0, 0, 0, 0, moscowLocation),
			moonTable: &MoonTable{
				Elems: []*MoonTableElement{
					{
						NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, utcLocation),
						FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, utcLocation),
						FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, utcLocation),
						LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, utcLocation),
						t1:           1.0,
						t2:           2.0,
					},
				},
			},
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, utcLocation),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{}

			result, err := cache.SearchPhase(tt.tGiven, tt.moonTable, tt.phase)

			if err != nil {
				t.Errorf("SearchPhase() for %s failed with error: %v", tt.name, err)
				return
			}

			if !result.Equal(tt.wantTime) {
				t.Errorf("SearchPhase() for %s gotTime = %v, want %v", tt.name, result, tt.wantTime)
			}

			if result.Location() != tt.wantTime.Location() {
				t.Errorf("SearchPhase() for %s location = %v, want %v",
					tt.name, result.Location(), tt.wantTime.Location())
			}
		})
	}
}

func TestSearchPhase_PhaseOrder(t *testing.T) {
	location := time.UTC

	moonTable := &MoonTable{
		Elems: []*MoonTableElement{
			{
				NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
				FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
				FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
				LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
				t1:           1.0,
				t2:           2.0,
			},
		},
	}

	tests := []struct {
		name     string
		tGiven   time.Time
		phase    EnumPhase
		wantTime time.Time
		desc     string
	}{
		{
			name:     "between new moon and first quarter - find first quarter",
			tGiven:   time.Date(2023, 1, 3, 0, 0, 0, 0, location),
			phase:    FirstQuarter,
			wantTime: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
		},
		{
			name:     "between first quarter and full moon - find full moon",
			tGiven:   time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			phase:    FullMoon,
			wantTime: time.Date(2023, 1, 15, 0, 0, 0, 0, location),
		},
		{
			name:     "between full moon and last quarter - find last quarter",
			tGiven:   time.Date(2023, 1, 18, 0, 0, 0, 0, location),
			phase:    LastQuarter,
			wantTime: time.Date(2023, 1, 22, 0, 0, 0, 0, location),
		},
	}

	cache := &Cache{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cache.SearchPhase(tt.tGiven, moonTable, tt.phase)

			if err != nil {
				t.Errorf("SearchPhase() for %s failed with error: %v", tt.name, err)
				return
			}

			if !result.Equal(tt.wantTime) {
				t.Errorf("SearchPhase() for %s gotTime = %v, want %v. %s",
					tt.name, result, tt.wantTime, tt.desc)
			}
		})
	}
}

func TestSearchPhase_SingleElement(t *testing.T) {
	location := time.UTC

	singleElementTable := &MoonTable{
		Elems: []*MoonTableElement{
			{
				NewMoon:      time.Date(2023, 1, 1, 0, 0, 0, 0, location),
				FirstQuarter: time.Date(2023, 1, 8, 0, 0, 0, 0, location),
				FullMoon:     time.Date(2023, 1, 15, 0, 0, 0, 0, location),
				LastQuarter:  time.Date(2023, 1, 22, 0, 0, 0, 0, location),
				t1:           1.0,
				t2:           2.0,
			},
		},
	}

	tests := []struct {
		name      string
		tGiven    time.Time
		phase     EnumPhase
		wantTime  time.Time
		wantError error
	}{
		{
			name:     "before all phases - return first phase",
			tGiven:   time.Date(2022, 12, 31, 0, 0, 0, 0, location),
			phase:    NewMoon,
			wantTime: time.Date(2023, 1, 1, 0, 0, 0, 0, location),
		},
		{
			name:      "after all phases - should try to create next table",
			tGiven:    time.Date(2023, 12, 31, 0, 0, 0, 0, location),
			phase:     NewMoon,
			wantError: errors.New("not found"),
		},
		{
			name:     "exactly in middle - find next phase",
			tGiven:   time.Date(2023, 1, 10, 0, 0, 0, 0, location),
			phase:    FullMoon,
			wantTime: time.Date(2023, 1, 15, 0, 0, 0, 0, location),
		},
	}

	cache := &Cache{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cache.SearchPhase(tt.tGiven, singleElementTable, tt.phase)

			if tt.wantError != nil {
				if err == nil {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				} else if err.Error() != tt.wantError.Error() {
					t.Errorf("Test: %v, SearchPhase() error = %v, wantErr %v", tt.name, err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Errorf("Test: %v, SearchPhase() unexpected error = %v", tt.name, err)
				return
			}

			if !result.Equal(tt.wantTime) {
				t.Errorf("Test: %v, SearchPhase() gotTime = %v, want %v", tt.name, result, tt.wantTime)
			}
		})
	}
}

func TestSearchPhase_Performance(t *testing.T) {
	location := time.UTC

	var elems []*MoonTableElement
	for i := 0; i < 100; i++ {
		baseTime := time.Date(2023, time.Month(i+1), 1, 0, 0, 0, 0, location)
		elem := &MoonTableElement{
			NewMoon:      baseTime,
			FirstQuarter: baseTime.AddDate(0, 0, 7),
			FullMoon:     baseTime.AddDate(0, 0, 14),
			LastQuarter:  baseTime.AddDate(0, 0, 21),
			t1:           float64(i * 2),
			t2:           float64(i*2 + 1),
		}
		elems = append(elems, elem)
	}

	largeMoonTable := &MoonTable{Elems: elems}
	cache := &Cache{}

	tGiven := time.Date(2023, 6, 15, 0, 0, 0, 0, location)
	result, err := cache.SearchPhase(tGiven, largeMoonTable, FullMoon)

	if err != nil {
		t.Errorf("SearchPhase() with large table failed: %v", err)
		return
	}

	expected := time.Date(2023, 7, 15, 0, 0, 0, 0, location)
	if !result.Equal(expected) {
		t.Errorf("SearchPhase() with large table gotTime = %v, want %v", result, expected)
	}
}
