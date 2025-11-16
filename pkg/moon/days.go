package moon

import (
	"math"
	"strings"
	"time"

	il "moon/pkg/illumination"
	jt "moon/pkg/julian-time"
	phase "moon/pkg/phase"
)

func (c *Cache) CreateMoonTable(timeGiven time.Time) *MoonTable {
	t := time.Date(timeGiven.Year(), 0, 0, 0, 0, 0, 0, timeGiven.Location())
	if c.moonTable != nil && c.moonTable[t.String()] != nil {
		return c.moonTable[t.String()]
	}
	moonTable := new(MoonTable)

	var l int
	var k1, mtime float64
	var minx int
	var phaset []float64

	phaset = make([]float64, 0)

	// Tabulate new and full moons surrounding the year
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

			elem.t1 = mp
			elem.t2 = lastnew

			lastnew = mp
		}

		elem.t1 = mp
		elem.t2 = lastnew

		firstQuarterTime := il.BinarySearchIllumination(lastnew, mp, timeGiven.Location(), true)
		elem.FirstQuarter = jt.FromJulianDate(firstQuarterTime, timeGiven.Location())

		lastQuarterTime := il.BinarySearchIllumination(mp, mp+10, timeGiven.Location(), false)
		elem.LastQuarter = jt.FromJulianDate(lastQuarterTime, timeGiven.Location())

		elem.NewMoon = jt.FromJulianDate(lastnew, timeGiven.Location())
		elem.FullMoon = jt.FromJulianDate(mp, timeGiven.Location())

		if elem.t1 != elem.t2 {
			moonTable.Elems = append(moonTable.Elems, elem)
		}

		if elem.LastQuarter.Year() > timeGiven.Year() {
			break
		}
	}
	if c.moonTable == nil {
		c.moonTable = make(map[string]*MoonTable)
	}
	if c.moonTable[t.String()] == nil {
		c.moonTable[t.String()] = moonTable
	}
	return moonTable
}

func BeginMoonDayToEarthDay(tGiven time.Time, duration time.Duration, timeFormat string, moonTable []*MoonTableElement) *any {
	var time any = time.Time{}
	if len(moonTable) == 0 {
		return &time
	}
	for i := range moonTable {
		elem := moonTable[i]
		if elem.t1 != elem.t2 {
			if tGiven.After(elem.NewMoon) && tGiven.Before(elem.LastQuarter) {
				t := elem.NewMoon
				t = t.Add(duration)
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
					t := elem.NewMoon
					t = t.Add(duration)

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
	}
	return &time
}

func (c *Cache) FindNearestPhase(tGiven time.Time) NearestPhase {
	var np NearestPhase
	table := c.CreateMoonTable(tGiven)

	if val, err := c.SearchPhase(tGiven, table, NewMoon); err == nil {
		np.NewMoon = val
	}
	if val, err := c.SearchPhase(tGiven, table, FirstQuarter); err == nil {
		np.FirstQuarter = val
	}
	if val, err := c.SearchPhase(tGiven, table, FullMoon); err == nil {
		np.FullMoon = val
	}
	if val, err := c.SearchPhase(tGiven, table, LastQuarter); err == nil {
		np.LastQuarter = val
	}

	return np
}

func GetMoonDays(tGiven time.Time, table []*MoonTableElement) time.Duration {
	var moonDays time.Duration
	for i := range table {
		elem := table[i]

		if elem.t1 != elem.t2 {
			if tGiven.After(elem.NewMoon) /*&& tGiven.Before(elem.TFull)*/ {
				moonDays = tGiven.Sub(elem.NewMoon)
			}
		}
	}
	return moonDays
}
