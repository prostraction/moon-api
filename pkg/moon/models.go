package moon

import "time"

type MoonDaysInDay struct {
	Begin   time.Duration
	Current time.Duration
	End     time.Duration
}

// used by:
// all moon package
type MoonTable struct {
	Elems []*MoonTableElement
}

type MoonTableElement struct {
	NewMoon      time.Time
	FirstQuarter time.Time
	FullMoon     time.Time
	LastQuarter  time.Time
}

type EnumPhase int

const (
	NewMoon EnumPhase = iota
	FirstQuarter
	FullMoon
	LastQuarter
)

// used by:
// - phase methods
type MoonDaysDetailed struct {
	Count int
	Day   []MoonDay
}

type MoonDay struct {
	Begin         *any `json:"Begin,omitempty"`
	IsBeginExists bool
	End           *any `json:"End,omitempty"`
	IsEndExists   bool
}

// used by:
// - route-phase
// - rounte-moon-table
type NearestPhase struct {
	NewMoon      time.Time
	FirstQuarter time.Time
	FullMoon     time.Time
	LastQuarter  time.Time
}
type NearestPhaseTimestamp struct {
	NewMoon      int64
	FirstQuarter int64
	FullMoon     int64
	LastQuarter  int64
}
type NearestPhaseString struct {
	NewMoon      string
	FirstQuarter string
	FullMoon     string
	LastQuarter  string
}

type SeachMoonDayResp struct {
	From time.Time
	To   time.Time
}

type SeachMoonDayRespTimestamp struct {
	From int64
	To   int64
}

type SeachMoonDayRespString struct {
	From string
	To   string
}
