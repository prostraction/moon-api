package server

import (
	"moon/pkg/moon"
	"moon/pkg/phase"
	pos "moon/pkg/position"
	"moon/pkg/zodiac"
	"time"
)

type ErrorPrintable struct {
	Status  int
	Message string
}

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
	BeginDay     *MoonStat `json:"BeginDay,omitempty"`
	CurrentState *MoonStat `json:"CurrentState,omitempty"`
	EndDay       *MoonStat `json:"EndDay,omitempty"`

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
