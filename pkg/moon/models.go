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
	Begin         *time.Time `json:"Begin,omitempty"`
	IsBeginExists bool
	End           *time.Time `json:"End,omitempty"`
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
