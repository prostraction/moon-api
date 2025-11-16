package server

import (
	"moon/pkg/moon"
	"moon/pkg/phase"
	pos "moon/pkg/position"
	"moon/pkg/zodiac"
	"time"
)

type MoonStat struct {
	MoonDays     float64
	Illumination float64
	Phase        phase.PhaseResp
	Zodiac       zodiac.Zodiac
	Position     *pos.MoonPosition `json:"MoonPosition,omitempty"`
}

type FullInfo struct {
	MoonDaysBegin   float64
	MoonDaysEnd     float64
	MoonDaysCurrent float64

	IlluminationBeginDay float64
	IlluminationCurrent  float64
	IlluminationEndDay   float64
}

type MoonDay struct {
	Begin time.Time
	End   time.Time
}

type MoonPhaseResponse struct {
	BeginDay     *MoonStat
	CurrentState *MoonStat
	EndDay       *MoonStat

	MoonDaysDetailed *moon.MoonDaysDetailed `json:"MoonDaysDetailed,omitempty"`
	ZodiacDetailed   *zodiac.Zodiacs

	MoonRiseAndSet *pos.DayData `json:"MoonRiseAndSet,omitempty"`

	info *FullInfo
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
	IsValid   bool
}

type JulianTimeResp struct {
	CivilDate          string
	CivilDateTimestamp int64
	JulianDate         float64
}

type NextMoonPhaseElement struct {
	Date      string
	Countdown int64
}

type NextMoonPhase struct {
	New   *NextMoonPhaseElement
	First *NextMoonPhaseElement
	Full  *NextMoonPhaseElement
	Last  *NextMoonPhaseElement
}
