package phase

import (
	"math"
	il "moon/pkg/illumination"
	ma "moon/pkg/math-helpers"
	"time"
)

type illumFunc func(tGiven time.Time, loc *time.Location) float64

func CurrentMoonPhase(tGiven time.Time, lang string) *PhaseResp {
	pr := new(PhaseResp)
	newT := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), tGiven.Hour(), tGiven.Minute(), tGiven.Second(), 0, time.UTC)

	currentMoonIllumination, currentMoonIlluminationBefore, currentMoonIlluminationAfter := currentMoonPhaseCalc(newT, time.FixedZone("UTC+12", 12*60*60), il.GetCurrentMoonIllumination)
	dayBeginMoonIllumination, dayBeginMoonIlluminationBefore, dayBeginMoonIlluminationAfter := currentMoonPhaseCalc(newT, time.FixedZone("UTC+12", 12*60*60), il.GetDailyMoonIllumination)
	dayEndMoonIllumination, dayEndMoonIlluminationBefore, dayEndMoonIlluminationAfter := currentMoonPhaseCalc(newT.AddDate(0, 0, 1), time.FixedZone("UTC+12", 12*60*60), il.GetDailyMoonIllumination)

	moonPhaseCurrent := GetMoonPhase(currentMoonIlluminationBefore, currentMoonIllumination, currentMoonIlluminationAfter, lang)
	moonPhaseBegin := GetMoonPhase(dayBeginMoonIlluminationBefore, dayBeginMoonIllumination, dayBeginMoonIlluminationAfter, lang)
	moonPhaseEnd := GetMoonPhase(dayEndMoonIlluminationBefore, dayEndMoonIllumination, dayEndMoonIlluminationAfter, lang)

	if dayBeginMoonIllumination <= currentMoonIllumination && currentMoonIllumination <= dayEndMoonIllumination {
		moonPhaseCurrent.IsWaxing = true
		moonPhaseBegin.IsWaxing = true
		moonPhaseEnd.IsWaxing = true
	} else if dayBeginMoonIllumination > currentMoonIllumination && currentMoonIllumination > dayEndMoonIllumination {
		moonPhaseCurrent.IsWaxing = false
		moonPhaseBegin.IsWaxing = false
		moonPhaseEnd.IsWaxing = false
	} else if dayBeginMoonIllumination > currentMoonIllumination {
		moonPhaseCurrent.IsWaxing = false
		moonPhaseBegin.IsWaxing = false
		moonPhaseEnd.IsWaxing = true
	} else if dayBeginMoonIllumination < currentMoonIllumination {
		moonPhaseCurrent.IsWaxing = true
		moonPhaseBegin.IsWaxing = true
		moonPhaseEnd.IsWaxing = false
	}

	pr.BeginDay = &moonPhaseBegin
	pr.Current = &moonPhaseCurrent
	pr.EndDay = &moonPhaseEnd

	pr.Illumination.BeginDay = dayBeginMoonIllumination
	pr.Illumination.Current = currentMoonIllumination
	pr.Illumination.EndDay = dayEndMoonIllumination

	return pr
}

func GetMoonPhase(before, current, after float64, lang string) Phase {
	phaseName, phangeNameLocalized, phaseEmoji := "", "", ""
	switch {
	case current > 0.05 && current < 0.45 && current < after:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 0)
		phaseName, phaseEmoji = getMoonPhases(0)
	case current >= 0.45 && current <= 0.55 && current < after:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 1)
		phaseName, phaseEmoji = getMoonPhases(1)
	case current > 0.55 && current < 0.95 && current > before:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 2)
		phaseName, phaseEmoji = getMoonPhases(2)
	case current >= 0.95:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 3)
		phaseName, phaseEmoji = getMoonPhases(3)
	case current < 0.95 && current > 0.55 && current < before:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 4)
		phaseName, phaseEmoji = getMoonPhases(4)
	case current <= 0.55 && current >= 0.45 && current < before:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 5)
		phaseName, phaseEmoji = getMoonPhases(5)
	case current < 0.45 && current > 0.05 && current < before:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 6)
		phaseName, phaseEmoji = getMoonPhases(6)
	case current <= 0.05:
		phangeNameLocalized = getMoonPhasesLocalized(lang, 7)
		phaseName, phaseEmoji = getMoonPhases(7)
	}
	return Phase{Name: phaseName, NameLocalized: phangeNameLocalized, Emoji: phaseEmoji}
}

func getMoonPhases(position int) (string, string) {
	return phasesEn[position], phasesEmoji[position]
}

func getMoonPhasesLocalized(lang string, position int) string {
	switch lang {
	case "en":
		return phasesEn[position]
	case "ru":
		return phasesRu[position]
	case "es":
		return phasesEs[position]
	case "de":
		return phasesDe[position]
	case "fr":
		return phasesFr[position]
	case "jp":
		return phasesJp[position]
	}
	return phasesEn[position]
}

func currentMoonPhaseCalc(tGiven time.Time, loc *time.Location, calcF illumFunc) (float64, float64, float64) {
	moonIllumination := calcF(tGiven, loc)
	moonIlluminationBefore := calcF(tGiven.AddDate(0, 0, -1), loc)
	moonIlluminationAfter := calcF(tGiven.AddDate(0, 0, 1), loc)

	// in rare UTC-12 case they are equal
	if moonIllumination == moonIlluminationBefore {
		moonIlluminationBefore = calcF(tGiven.AddDate(0, 0, -2), loc)
	}
	// just in case
	if moonIllumination == moonIlluminationAfter {
		moonIlluminationAfter = calcF(tGiven.AddDate(0, 0, 2), loc)
	}

	return moonIllumination, moonIlluminationBefore, moonIlluminationAfter
}

func Truephase(k, phase float64) float64 {
	var t, t2, t3, pt, m, mprime, f float64
	// to do
	SynMonth := 29.53058868 // Synodic month (mean time from new to next new Moon)

	k += phase           // Add phase to new moon time
	t = k / 1236.85      // Time in Julian centuries from 1900 January 0.5
	t2 = t * t           // Square for frequent use
	t3 = t2 * t          // Cube for frequent use
	pt = 2415020.75933 + // Mean time of phase
		SynMonth*k +
		0.0001178*t2 -
		0.000000155*t3 +
		0.00033*ma.Dsin(166.56+132.87*t-0.009173*t2)

	m = 359.2242 + // Sun's mean anomaly
		29.10535608*k -
		0.0000333*t2 -
		0.00000347*t3
	mprime = 306.0253 + // Moon's mean anomaly
		385.81691806*k +
		0.0107306*t2 +
		0.00001236*t3
	f = 21.2964 + // Moon's argument of latitude
		390.67050646*k -
		0.0016528*t2 -
		0.00000239*t3

	if (phase < 0.01) || (math.Abs(phase-0.5) < 0.01) {
		// Corrections for New and Full Moon
		pt += (0.1734-0.000393*t)*ma.Dsin(m) +
			0.0021*ma.Dsin(2*m) -
			0.4068*ma.Dsin(mprime) +
			0.0161*ma.Dsin(2*mprime) -
			0.0004*ma.Dsin(3*mprime) +
			0.0104*ma.Dsin(2*f) -
			0.0051*ma.Dsin(m+mprime) -
			0.0074*ma.Dsin(m-mprime) +
			0.0004*ma.Dsin(2*f+m) -
			0.0004*ma.Dsin(2*f-m) -
			0.0006*ma.Dsin(2*f+mprime) +
			0.0010*ma.Dsin(2*f-mprime) +
			0.0005*ma.Dsin(m+2*mprime)
	} else if (math.Abs(phase-0.25) < 0.01) || (math.Abs(phase-0.75) < 0.01) {
		pt += (0.1721-0.0004*t)*ma.Dsin(m) +
			0.0021*ma.Dsin(2*m) -
			0.6280*ma.Dsin(mprime) +
			0.0089*ma.Dsin(2*mprime) -
			0.0004*ma.Dsin(3*mprime) +
			0.0079*ma.Dsin(2*f) -
			0.0119*ma.Dsin(m+mprime) -
			0.0047*ma.Dsin(m-mprime) +
			0.0003*ma.Dsin(2*f+m) -
			0.0004*ma.Dsin(2*f-m) -
			0.0006*ma.Dsin(2*f+mprime) +
			0.0021*ma.Dsin(2*f-mprime) +
			0.0003*ma.Dsin(m+2*mprime) +
			0.0004*ma.Dsin(m-2*mprime) -
			0.0003*ma.Dsin(2*m+mprime)

		if phase < 0.5 {
			// First quarter correction
			pt += 0.0028 - 0.0004*ma.Dcos(m) + 0.0003*ma.Dcos(mprime)
		} else {
			// Last quarter correction
			pt += -0.0028 + 0.0004*ma.Dcos(m) - 0.0003*ma.Dcos(mprime)
		}
	}
	return pt
}
