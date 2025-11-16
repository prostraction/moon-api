package server

import (
	jt "moon/pkg/julian-time"
	"moon/pkg/phase"
	pos "moon/pkg/position"
	"moon/pkg/zodiac"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

/*    MOON PHASE    */
func (s *Server) moonPhaseCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	tGiven = tGiven.In(loc)
	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	return s.moonPhaseV1(c, tGiven, precision, locationCords)
}

func (s *Server) moonPhaseTimestampV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tStr := c.Query("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}
	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, time.Local)
	tGiven = tGiven.In(loc)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	return s.moonPhaseV1(c, tGiven, precision, locationCords)
}

func (s *Server) moonPhaseDatetV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tNow := time.Now()

	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := StrToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)
	day := StrToInt(c.Query("day", strconv.Itoa(int(tNow.Day()))), int(tNow.Day()), 1, 31)

	err := IsValidDate(year, month, day)
	if err != nil {
		c.Status(400)
		return c.SendString("Validation error: " + err.Error())
	}

	hour := StrToInt(c.Query("hour", strconv.Itoa(int(tNow.Hour()))), int(tNow.Hour()), 0, 23)
	minute := StrToInt(c.Query("minute", strconv.Itoa(int(tNow.Minute()))), int(tNow.Minute()), 0, 59)
	second := StrToInt(c.Query("second", strconv.Itoa(int(tNow.Second()))), int(tNow.Second()), 0, 59)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	tGiven := time.Date(year, jt.GetMonth(month), day, hour, minute, second, 0, time.Local)
	tGiven = tGiven.In(loc)
	return s.moonPhaseV1(c, tGiven, precision, locationCords)
}

func (s *Server) moonPhaseV1(c *fiber.Ctx, tGiven time.Time, precision int, locationCords Coordinates) error {
	lang := c.Query("lang", "en")
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	resp := MoonPhaseResponse{}
	resp.EndDay = new(MoonStat)
	resp.CurrentState = new(MoonStat)
	resp.BeginDay = new(MoonStat)
	resp.info = new(FullInfo)

	var err error

	var beginDuration, currentDuration, endDuration time.Duration
	beginDuration, currentDuration, endDuration = s.moonCache.CurrentMoonDays(tGiven, loc)

	resp.info.MoonDaysBegin = beginDuration.Minutes() / jt.Fminute
	resp.info.MoonDaysCurrent = currentDuration.Minutes() / jt.Fminute
	resp.info.MoonDaysEnd = endDuration.Minutes() / jt.Fminute

	resp.BeginDay.MoonDays = ToFixed(resp.info.MoonDaysBegin, precision)
	resp.CurrentState.MoonDays = ToFixed(resp.info.MoonDaysCurrent, precision)
	resp.EndDay.MoonDays = ToFixed(resp.info.MoonDaysEnd, precision)

	resp.info.IlluminationCurrent, resp.info.IlluminationBeginDay, resp.info.IlluminationEndDay, resp.CurrentState.Phase, resp.BeginDay.Phase, resp.EndDay.Phase = phase.CurrentMoonPhase(tGiven, lang)

	resp.BeginDay.Illumination = ToFixed(resp.info.IlluminationBeginDay*100, precision)
	resp.CurrentState.Illumination = ToFixed(resp.info.IlluminationCurrent*100, precision)
	resp.EndDay.Illumination = ToFixed(resp.info.IlluminationEndDay*100, precision)

	resp.ZodiacDetailed, resp.BeginDay.Zodiac, resp.CurrentState.Zodiac, resp.EndDay.Zodiac = zodiac.CurrentZodiacs(tGiven, loc, lang, s.moonCache.CreateMoonTable(tGiven).Elems)

	if locationCords.IsValid {
		resp.MoonRiseAndSet, err = s.positionCache.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}
		resp.MoonDaysDetailed = s.moonCache.MoonDetailed(tGiven, loc, lang, locationCords.Longitude, locationCords.Latitude)

		newT := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, tGiven.Location())
		resp.BeginDay.Position, err = pos.GetMoonPosition(newT, newT.Location(), precision, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}

		resp.CurrentState.Position, err = pos.GetMoonPosition(tGiven, tGiven.Location(), precision, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}

		resp.EndDay.Position, err = pos.GetMoonPosition(newT.AddDate(0, 0, 1), newT.AddDate(0, 0, 1).Location(), precision, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}
	} else {
		resp.MoonRiseAndSet, err = s.positionCache.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision)
	}

	if err != nil && err.Error() != "no location prodived" {
		log.Error(err.Error())
	}

	return c.JSON(resp)
}
