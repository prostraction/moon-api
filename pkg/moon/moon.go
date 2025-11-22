package moon

import (
	"errors"
	"log"
	pos "moon/pkg/position"
	"time"
)

// to do
func CurrentMoonDays(tGiven time.Time, loc *time.Location, moonTable *MoonTable) MoonDaysInDay {
	var mday MoonDaysInDay
	if moonTable == nil {
		return mday
	}

	if loc == nil {
		loc = time.UTC
	}
	currentDayTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), tGiven.Hour(), tGiven.Minute(), tGiven.Second(), tGiven.Nanosecond(), loc)

	mday, _ = GetMoonDays(currentDayTime, moonTable.Elems)
	return mday
}

func MoonDetailed(tGiven time.Time, loc *time.Location, lang string, timeFormat string, longitude float64, latitude float64) *MoonDaysDetailed {
	if loc == nil {
		loc = time.UTC
	}

	dayYesterday := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()-1, 0, 0, 0, 0, loc)
	dayToday := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayTomorrow := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonDaysDetailed := new(MoonDaysDetailed)
	moonDaysDetailed.Day = make([]MoonDay, 2)
	moonDaysDetailed.Count = 2

	moonRiseYesterday, err1 := pos.GetRisesDay(dayYesterday.Year(), int(dayYesterday.Month()), dayYesterday.Day(), loc, 2, timeFormat, longitude, latitude)
	moonRiseToday, err2 := pos.GetRisesDay(dayToday.Year(), int(dayToday.Month()), dayToday.Day(), loc, 2, timeFormat, longitude, latitude)
	moonRiseTomorrow, err3 := pos.GetRisesDay(dayTomorrow.Year(), int(dayTomorrow.Month()), dayTomorrow.Day(), loc, 2, timeFormat, longitude, latitude)

	if err1 == nil && err2 == nil {
		if moonRiseYesterday.IsMoonRise {
			moonDaysDetailed.Day[0].Begin = moonRiseYesterday.Moonrise.Time
			moonDaysDetailed.Day[0].IsBeginExists = true
		}
		if moonRiseToday.IsMoonRise {
			moonDaysDetailed.Day[0].End = moonRiseToday.Moonrise.Time
			moonDaysDetailed.Day[0].IsEndExists = true
		}
	}
	if err2 == nil && err3 == nil {
		if moonRiseToday.IsMoonRise {
			moonDaysDetailed.Day[1].Begin = moonRiseToday.Moonrise.Time
			moonDaysDetailed.Day[1].IsBeginExists = true
		}
		if moonRiseTomorrow.IsMoonRise {
			moonDaysDetailed.Day[1].End = moonRiseTomorrow.Moonrise.Time
			moonDaysDetailed.Day[1].IsEndExists = true
		}
	}

	if !(moonDaysDetailed.Day[0].IsBeginExists && moonDaysDetailed.Day[0].IsEndExists) {
		moonDaysDetailed.Count = 1
	}

	return moonDaysDetailed
}

func SearchPhase(tGiven time.Time, moonTable *MoonTable, phase EnumPhase) (t time.Time, err error) {
	if moonTable == nil {
		err = errors.New("passed empty moonTable to SearchNewMoon")
		log.Println(err.Error())
		return
	}
	err = errors.New("not found")
	for i := range moonTable.Elems {
		elem := moonTable.Elems[i]
		elemSearch1 := elem.NewMoon
		if tGiven.Before(elemSearch1) {
			elemSearch1 = time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()-1, 0, 0, 0, 0, tGiven.Location())
		}
		elemSearch2 := elem.LastQuarter

		// range current phase:
		if tGiven.After(elemSearch1) && tGiven.Before(elemSearch2) {
			// found in current phase
			switch phase {
			case NewMoon:
				if tGiven.Before(elem.NewMoon) {
					return elem.NewMoon, nil
				}
			case FirstQuarter:
				if tGiven.Before(elem.FirstQuarter) {
					return elem.FirstQuarter, nil
				}
			case FullMoon:
				if tGiven.Before(elem.FullMoon) {
					return elem.FullMoon, nil
				}
			case LastQuarter:
				if tGiven.Before(elem.LastQuarter) {
					return elem.LastQuarter, nil
				}
			}
			// found in next phase
			if i < len(moonTable.Elems)-1 {
				// use values if in table
				switch phase {
				case NewMoon:
					return moonTable.Elems[i+1].NewMoon, nil
				case FirstQuarter:
					return moonTable.Elems[i+1].FirstQuarter, nil
				case FullMoon:
					return moonTable.Elems[i+1].FullMoon, nil
				case LastQuarter:
					return moonTable.Elems[i+1].LastQuarter, nil
				}
			} else {
				// create table for next table
				newT := time.Date(tGiven.Year()+1, 0, 0, 0, 0, 0, 0, tGiven.Location())
				newMoonTable := CreateMoonTable(newT)

				if newMoonTable != nil && newMoonTable.Elems != nil && len(newMoonTable.Elems) > 0 {
					switch phase {
					case NewMoon:
						if tGiven.Before(newMoonTable.Elems[0].NewMoon) {
							return newMoonTable.Elems[0].NewMoon, nil
						}
					case FirstQuarter:
						if tGiven.Before(newMoonTable.Elems[0].FirstQuarter) {
							return newMoonTable.Elems[0].FirstQuarter, nil
						}
					case FullMoon:
						if tGiven.Before(newMoonTable.Elems[0].FullMoon) {
							return newMoonTable.Elems[0].FullMoon, nil
						}
					case LastQuarter:
						if tGiven.Before(newMoonTable.Elems[0].LastQuarter) {
							return newMoonTable.Elems[0].LastQuarter, nil
						}
					}
				}
			}
		}
		// range next phase
		if i < len(moonTable.Elems)-1 {
			// try to find in current table
			elem2 := moonTable.Elems[i+1]
			elemSearch1 = elem.LastQuarter

			switch phase {
			case NewMoon:
				elemSearch2 = elem2.NewMoon
			case FirstQuarter:
				elemSearch2 = elem2.FirstQuarter
			case FullMoon:
				elemSearch2 = elem2.FullMoon
			case LastQuarter:
				elemSearch2 = elem2.LastQuarter
			}

			if tGiven.After(elemSearch1) && tGiven.Before(elemSearch2) {
				return elemSearch2, nil
			}

		} else {
			// try to find in next table
			newT := time.Date(tGiven.Year()+1, 0, 0, 0, 0, 0, 0, tGiven.Location())
			newMoonTable := CreateMoonTable(newT)
			if newMoonTable != nil && newMoonTable.Elems != nil && len(newMoonTable.Elems) > 0 {
				elemSearch1 = elem.LastQuarter

				switch phase {
				case NewMoon:
					elemSearch2 = newMoonTable.Elems[0].NewMoon
				case FirstQuarter:
					elemSearch2 = newMoonTable.Elems[0].FirstQuarter
				case FullMoon:
					elemSearch2 = newMoonTable.Elems[0].FullMoon
				case LastQuarter:
					elemSearch2 = newMoonTable.Elems[0].LastQuarter
				}

				if tGiven.After(elemSearch1) && tGiven.Before(elemSearch2) {
					return elemSearch2, nil
				}
			}
		}
	}
	return
}
