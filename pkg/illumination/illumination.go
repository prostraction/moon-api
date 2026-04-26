package illimination

import (
	"math"
	jt "moon/pkg/julian-time"
	m "moon/pkg/math-helpers"
	"time"
)

const toRad = math.Pi / 180.

func GetDailyMoonIllumination(tGiven time.Time, loc *time.Location) float64 {
	dailyTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, time.UTC)
	h, m, err := jt.GetTimeFromLocation(loc)
	h = -h
	m = -m
	if err == nil {
		dailyTime = dailyTime.Add(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
	}
	return getIlluminatedFractionOfMoon(jt.ToJulianDate(dailyTime))
}

func GetCurrentMoonIllumination(tGiven time.Time, loc *time.Location) float64 {
	tGiven = time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), tGiven.Hour(), tGiven.Minute(), tGiven.Second(), 0, time.UTC)
	h, m, err := jt.GetTimeFromLocation(loc)
	h = -h
	m = -m
	if err == nil {
		tGiven = tGiven.Add(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
	}
	return getIlluminatedFractionOfMoon(jt.ToJulianDate(tGiven))
}

func getIlluminatedFractionOfMoon(jd float64) float64 {
	T := (jd - 2451545.) / 36525.

	D := m.Constrain(297.8501921+445267.1114034*T-0.0018819*T*T+1./545868.0*T*T*T-1./113065000.0*T*T*T*T) * toRad
	M := m.Constrain(357.5291092+35999.0502909*T-0.0001536*T*T+1./24490000.0*T*T*T) * toRad
	Mp := m.Constrain(134.9633964+477198.8675055*T+0.0087414*T*T+1./69699.0*T*T*T-1./14712000.0*T*T*T*T) * toRad

	i := m.Constrain(180.-D*180./math.Pi-6.289*math.Sin(Mp)+2.1*math.Sin(M)-1.274*math.Sin(2.*D-Mp)-0.658*math.Sin(2.*D)-0.214*math.Sin(2.*Mp)-0.11*math.Sin(D)) * toRad

	return (1. + math.Cos(i)) / 2.
}

func BinarySearchIllumination(jdTimeBegin, jdTimeEnd float64, loc *time.Location, direction bool) (jdTime float64) {
	if loc == nil {
		loc = time.UTC
	}

	it := 0
	low := jdTimeBegin
	high := jdTimeEnd

	mid := low + (high-low)/2.0
	illum := GetCurrentMoonIllumination(jt.FromJulianDate(mid, loc), loc)
	if !direction && illum < 0.5 {
		direction = true
	}

	for it < 50 {
		mid = low + (high-low)/2.0
		illum = GetCurrentMoonIllumination(jt.FromJulianDate(mid, loc), loc)
		if math.Abs(illum-0.5) < 0.0001 {
			return mid
		}
		if direction {
			if illum < 0.5 {
				low = mid
			} else {
				high = mid
			}
		} else {
			if illum > 0.5 {
				low = mid
			} else {
				high = mid
			}
		}
		if high-low < 1e-10 {
			return mid
		}
		it++
	}
	return low + (high-low)/2.0
}
