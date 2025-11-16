package moon

import "time"

type EnumPhase int

const (
	NewMoon EnumPhase = iota
	FirstQuarter
	FullMoon
	LastQuarter
)

type MoonTable struct {
	Elems []*MoonTableElement
}

type MoonTableElement struct {
	NewMoon      time.Time
	FirstQuarter time.Time
	FullMoon     time.Time
	LastQuarter  time.Time
	t1           float64
	t2           float64
}

type Cache struct {
	moonTable map[string]*MoonTable
}

type MoonDay struct {
	Begin         *any `json:"Begin,omitempty"`
	IsBeginExists bool
	End           *any `json:"End,omitempty"`
	IsEndExists   bool
}

type MoonDaysDetailed struct {
	Count int
	Day   []MoonDay
}

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
