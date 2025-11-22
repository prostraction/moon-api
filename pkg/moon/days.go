package moon

import (
	"math"
	"strings"
	"time"

	il "moon/pkg/illumination"
	jt "moon/pkg/julian-time"
	phase "moon/pkg/phase"

	"github.com/gofiber/fiber/v2/log"
)

// to do
func CreateMoonTable(timeGiven time.Time) *MoonTable {
	moonTable := new(MoonTable)

	var l int
	var k1, mtime float64
	var minx int
	var phaset []float64

	phaset = make([]float64, 0)

	// tabulate new and full moons surrounding the year
	k1 = math.Floor((float64(timeGiven.Year()) - 1900) * 12.3685)
	minx = 0
	isNext := true
	for l = 0; ; l++ {
		mtime = phase.Truephase(k1, float64(l&1)*0.5)
		datey, _, _ := jt.Jyear(mtime)
		if datey >= timeGiven.Year() {
			minx++
		}
		phaseSign := mtime
		if (l & 1) == 0 {
			phaseSign = -mtime
		}
		phaset = append(phaset, phaseSign)
		if !isNext {
			break
		}
		if datey > timeGiven.Year() {
			minx++
			isNext = false
		}
		if (l & 1) != 0 {
			k1 += 1
		}
	}

	var lastnew float64
	for l = range minx {
		elem := &MoonTableElement{}
		mp := phaset[l]
		if mp < 0 {
			mp = -mp
			lastnew = mp
		}

		firstQuarterTime := il.BinarySearchIllumination(lastnew, mp, timeGiven.Location(), true)
		elem.FirstQuarter = jt.FromJulianDate(firstQuarterTime, timeGiven.Location())

		lastQuarterTime := il.BinarySearchIllumination(mp, mp+10, timeGiven.Location(), false)
		elem.LastQuarter = jt.FromJulianDate(lastQuarterTime, timeGiven.Location())

		elem.NewMoon = jt.FromJulianDate(lastnew, timeGiven.Location())
		elem.FullMoon = jt.FromJulianDate(mp, timeGiven.Location())

		if mp != lastnew {
			moonTable.Elems = append(moonTable.Elems, elem)
		}

		if elem.LastQuarter.Year() > timeGiven.Year() {
			break
		}
	}

	return moonTable
}

// to do
func BeginMoonDayToEarthDay(tGiven time.Time, duration time.Duration, timeFormat string, moonTable []*MoonTableElement) *any {
	var tt any = time.Time{}
	if len(moonTable) == 0 {
		return &tt
	}
	for i := range moonTable {
		elem := moonTable[i]
		if tGiven.After(elem.NewMoon) && tGiven.Before(elem.NewMoon.Add(time.Hour*24*32)) {
			t := elem.NewMoon.Add(duration)
			if strings.ToLower(timeFormat) == "timestamp" {
				var tRet any = t.Unix()
				return &tRet
			}

			if strings.ToLower(timeFormat) != "iso" {
				var tRet any = t.Format(timeFormat)
				return &tRet
			}

			var tRet any = t
			return &tRet
		}
		if i < len(moonTable)-1 {
			elem2 := moonTable[i+1]
			if tGiven.After(elem.LastQuarter) && tGiven.Before(elem2.NewMoon) {
				t := elem.NewMoon.Add(duration)

				if strings.ToLower(timeFormat) == "timestamp" {
					var tRet any = t.Unix()
					return &tRet
				}

				if strings.ToLower(timeFormat) != "iso" {
					var tRet any = t.Format(timeFormat)
					return &tRet
				}

				var tRet any = t
				return &tRet
			}
		}
	}
	return &tt
}

// to do
func FindNearestPhase(tGiven time.Time, moonTable *MoonTable) NearestPhase {
	if moonTable == nil {
		return NearestPhase{}
	}

	var np NearestPhase

	if val, err := SearchPhase(tGiven, moonTable, NewMoon); err == nil {
		np.NewMoon = val
	}
	if val, err := SearchPhase(tGiven, moonTable, FirstQuarter); err == nil {
		np.FirstQuarter = val
	}
	if val, err := SearchPhase(tGiven, moonTable, FullMoon); err == nil {
		np.FullMoon = val
	}
	if val, err := SearchPhase(tGiven, moonTable, LastQuarter); err == nil {
		np.LastQuarter = val
	}

	return np
}

func GetMoonDays(tGiven time.Time, table []*MoonTableElement) (MoonDaysInDay, time.Duration) {
	tBegin := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, tGiven.Location())
	for i := range table {
		elem := table[i]
		if tBegin.After(elem.NewMoon) && tBegin.Before(elem.NewMoon.Add(time.Hour*24*32)) {
			var elem2 *MoonTableElement
			if i < len(table)-1 {
				elem2 = table[i+1] // new next moon
			} else {
				log.Error("if i < len(table)-1 { // fix it, it is required")
				elem2 = new(MoonTableElement) // new next moon approx (fix?) to do
				elem2.NewMoon = elem.NewMoon.AddDate(0, 0, 29)
				elem2.NewMoon = elem.NewMoon.Add(time.Hour * 12)
			}
			moonMonDays := elem2.NewMoon.Sub(elem.NewMoon) // moon month

			eartbeg := tBegin.Add(-tGiven.Sub(elem.NewMoon))
			eartend := time.Date(eartbeg.Year(), eartbeg.Month()+1, eartbeg.Day(), eartbeg.Hour(), eartbeg.Minute(), eartbeg.Second(), 0, tBegin.Location())
			eartMon := eartend.Unix() - eartbeg.Unix() // earth month

			beginDay := elem.NewMoon
			currentDay := elem.NewMoon
			day := time.Hour * time.Duration(int64(moonMonDays.Seconds()/float64(eartMon)*24.))

			for tBegin.Sub(beginDay) > day {
				beginDay = beginDay.Add(day)
				currentDay = currentDay.Add(day)
			}

			currentDay = currentDay.Add(time.Hour * time.Duration(int64(moonMonDays.Seconds()/float64(eartMon)*float64(tGiven.Hour()))))
			currentDay = currentDay.Add(time.Minute * time.Duration(int64(moonMonDays.Seconds()/float64(eartMon)*float64(tGiven.Minute()))))

			endDay := beginDay.Add(day)
			return MoonDaysInDay{
				Begin:   beginDay.Sub(elem.NewMoon) % moonMonDays,
				Current: currentDay.Sub(elem.NewMoon) % moonMonDays,
				End:     endDay.Sub(elem.NewMoon) % moonMonDays,
			}, moonMonDays
		}
	}
	return MoonDaysInDay{time.Duration(0), time.Duration(0), time.Duration(0)}, time.Duration(0)
}
