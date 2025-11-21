package zodiac

import (
	"log"
	jt "moon/pkg/julian-time"
	moon "moon/pkg/moon"
	"time"
)

func CurrentZodiacs(tGiven time.Time, loc *time.Location, lang string, timeFormat string, moonTable []*moon.MoonTableElement) (*Zodiacs, Zodiac, Zodiac, Zodiac) {
	zods := new(Zodiacs)

	zodiacBegin := Zodiac{}
	zodiacCurrent := Zodiac{}
	zodiacEnd := Zodiac{}

	dayBeginTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	beginMoonDays, _, endMoonDays, _ := moon.GetMoonDaysPrecise(dayBeginTime, moonTable)

	zodiacPositionBegin := int((beginMoonDays.Minutes()/jt.Fminute*360.)/30./30.) % 12
	zodiacPositionCurrent := int(((beginMoonDays.Minutes()+float64(tGiven.Hour()*60)+float64(tGiven.Minute()))/jt.Fminute*360.)/30./30.) % 12
	zodiacPositionEnd := int(((beginMoonDays.Minutes()+1440)/jt.Fminute*360.)/30./30.) % 12

	log.Println(zodiacPositionBegin, zodiacPositionEnd)

	if zodiacPositionBegin == zodiacPositionEnd {
		zods.Count = 1
		zodBegin := zodiacPositionBegin * jt.Fminute / 360 * 30. * 30.
		zodEnd := (zodiacPositionEnd + 1) * jt.Fminute / 360 * 30. * 30.
		var tBegin any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodBegin)*time.Minute, timeFormat, moonTable)
		var tEnd any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodEnd)*time.Minute, timeFormat, moonTable)
		zods.Zodiac = make([]ZodiacDetailed, 1)
		zods.Zodiac[0].Begin = &tBegin
		zods.Zodiac[0].End = &tEnd
		zods.Zodiac[0].Name, zods.Zodiac[0].Emoji = getZodiacResp(zodiacPositionBegin)
		zods.Zodiac[0].NameLocalized = getZodiacRespLocalized(zodiacPositionBegin, lang)
	} else {
		zods.Count = 2
		zodBegin1 := (zodiacPositionBegin) * jt.Fminute / 360 * 30. * 30.
		zodEnd1 := (zodiacPositionBegin + 1) * jt.Fminute / 360 * 30. * 30.

		var tBegin1 any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodBegin1)*time.Minute, timeFormat, moonTable)
		var tEnd1 any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodEnd1)*time.Minute, timeFormat, moonTable)
		zods.Zodiac = make([]ZodiacDetailed, 2)
		zods.Zodiac[0].Begin = &tBegin1
		zods.Zodiac[0].End = &tEnd1
		zods.Zodiac[0].Name, zods.Zodiac[0].Emoji = getZodiacResp(zodiacPositionBegin)
		zods.Zodiac[0].NameLocalized = getZodiacRespLocalized(zodiacPositionBegin, lang)

		if int(endMoonDays.Minutes()/jt.Fminute) == 0 {
			endMoonDays += (beginMoonDays + 24*time.Hour)
			zodiacPositionEnd = int((endMoonDays.Minutes()/jt.Fminute*360.)/30./30.) % 12
		}

		zodBegin2 := (zodiacPositionEnd) * jt.Fminute / 360 * 30. * 30.
		zodEnd2 := (zodiacPositionEnd + 1) * jt.Fminute / 360 * 30. * 30.
		var tBegin2 any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodBegin2)*time.Minute, timeFormat, moonTable)
		var tEnd2 any = moon.BeginMoonDayToEarthDay(tGiven, time.Duration(zodEnd2)*time.Minute, timeFormat, moonTable)
		zods.Zodiac[1].Begin = &tBegin2
		zods.Zodiac[1].End = &tEnd2
		zods.Zodiac[1].Name, zods.Zodiac[1].Emoji = getZodiacResp(zodiacPositionEnd)
		zods.Zodiac[1].NameLocalized = getZodiacRespLocalized(zodiacPositionEnd, lang)
	}

	zodiacBegin.Name, zodiacBegin.Emoji = getZodiacResp(zodiacPositionBegin)
	zodiacCurrent.Name, zodiacCurrent.Emoji = getZodiacResp(zodiacPositionCurrent)
	zodiacEnd.Name, zodiacEnd.Emoji = getZodiacResp(zodiacPositionEnd)

	zodiacBegin.NameLocalized = getZodiacRespLocalized(zodiacPositionBegin, lang)
	zodiacCurrent.NameLocalized = getZodiacRespLocalized(zodiacPositionCurrent, lang)
	zodiacEnd.NameLocalized = getZodiacRespLocalized(zodiacPositionEnd, lang)

	return zods, zodiacBegin, zodiacCurrent, zodiacEnd
}

func getZodiacResp(position int) (string, string) {
	if position >= 0 && position < len(signsEn) && position < len(signsEmoji) {
		return signsEn[position], signsEmoji[position]
	}
	return "", ""
}

func getZodiacRespLocalized(position int, lang string) string {
	switch lang {
	case "en":
		if position >= 0 && position < len(signsEn) && position < len(signsEmoji) {
			return signsEn[position]
		}
	case "ru":
		if position >= 0 && position < len(signsRu) && position < len(signsEmoji) {
			return signsRu[position]
		}
	case "es":
		if position >= 0 && position < len(signsEs) && position < len(signsEmoji) {
			return signsEs[position]
		}
	case "de":
		if position >= 0 && position < len(signsDe) && position < len(signsEmoji) {
			return signsDe[position]
		}
	case "fr":
		if position >= 0 && position < len(signsFr) && position < len(signsEmoji) {
			return signsFr[position]
		}
	case "jp":
		if position >= 0 && position < len(signsJp) && position < len(signsEmoji) {
			return signsJp[position]
		}
	default:
		if position >= 0 && position < len(signsEn) && position < len(signsEmoji) {
			return signsEn[position]
		}
	}
	return ""
}
