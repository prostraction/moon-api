package server

import (
	jt "moon/pkg/julian-time"
	phase "moon/pkg/phase"
	pos "moon/pkg/position"
	zodiac "moon/pkg/zodiac"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Server) versionV1(c *fiber.Ctx) error {
	return c.JSON("1.1.2")
}

/*    MOON PHASE    */
func (s *Server) moonPhaseCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

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
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, loc)

	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

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

	year := strToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := strToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)
	day := strToInt(c.Query("day", strconv.Itoa(int(tNow.Day()))), int(tNow.Day()), 1, 31)

	err := isValidDate(year, month, day)
	if err != nil {
		c.Status(400)
		return c.SendString("Validation error: " + err.Error())
	}

	hour := strToInt(c.Query("hour", strconv.Itoa(int(tNow.Hour()))), int(tNow.Hour()), 0, 23)
	minute := strToInt(c.Query("minute", strconv.Itoa(int(tNow.Minute()))), int(tNow.Minute()), 0, 59)
	second := strToInt(c.Query("second", strconv.Itoa(int(tNow.Second()))), int(tNow.Second()), 0, 59)

	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	tGiven := time.Date(year, jt.GetMonth(month), day, hour, minute, second, 0, loc)
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

	resp.BeginDay.MoonDays = toFixed(resp.info.MoonDaysBegin, precision)
	resp.CurrentState.MoonDays = toFixed(resp.info.MoonDaysCurrent, precision)
	resp.EndDay.MoonDays = toFixed(resp.info.MoonDaysEnd, precision)

	resp.info.IlluminationCurrent, resp.info.IlluminationBeginDay, resp.info.IlluminationEndDay, resp.CurrentState.Phase, resp.BeginDay.Phase, resp.EndDay.Phase = phase.CurrentMoonPhase(tGiven, lang)

	resp.BeginDay.Illumination = toFixed(resp.info.IlluminationBeginDay*100, precision)
	resp.CurrentState.Illumination = toFixed(resp.info.IlluminationCurrent*100, precision)
	resp.EndDay.Illumination = toFixed(resp.info.IlluminationEndDay*100, precision)

	resp.ZodiacDetailed, resp.BeginDay.Zodiac, resp.CurrentState.Zodiac, resp.EndDay.Zodiac = zodiac.CurrentZodiacs(tGiven, loc, lang, s.moonCache.CreateMoonTable(tGiven))

	if locationCords.IsValid {
		resp.MoonRiseAndSet, err = pos.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision, locationCords.Longitude, locationCords.Latitude)
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
		resp.MoonRiseAndSet, err = pos.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision)
	}

	if err != nil && err.Error() != "no location prodived" {
		log.Error(err.Error())
	}

	return c.JSON(resp)
}

/*    MOON TABLE    */

func (s *Server) moonTableCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableYearV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tNow := time.Now()
	year := strToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)

	tGiven := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableV1(c *fiber.Ctx, tGiven time.Time) error {
	resp := MoonTable{}
	resp.Table = s.moonCache.GenerateMoonTable(tGiven)
	return c.JSON(resp.Table)
}

/* Julian Time methods */
func (s *Server) toJulianTimeByDateV1(c *fiber.Ctx) error {
	tNow := time.Now()
	tNow = tNow.In(time.UTC)

	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

	year := strToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := strToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)
	day := strToInt(c.Query("day", strconv.Itoa(int(tNow.Day()))), int(tNow.Day()), 1, 31)

	err := isValidDate(year, month, day)
	if err != nil {
		c.Status(400)
		return c.SendString("Validation error: " + err.Error())
	}

	hour := strToInt(c.Query("hour", strconv.Itoa(int(tNow.Hour()))), int(tNow.Hour()), 0, 23)
	minute := strToInt(c.Query("minute", strconv.Itoa(int(tNow.Minute()))), int(tNow.Minute()), 0, 59)
	second := strToInt(c.Query("second", strconv.Itoa(int(tNow.Second()))), int(tNow.Second()), 0, 59)

	tGiven := time.Date(year, jt.GetMonth(month), day, hour, minute, second, 0, time.UTC)
	return s.toJulianTimeV1(c, tGiven, precision)
}

func (s *Server) toJulianTimeByTimestampV1(c *fiber.Ctx) error {
	tStr := c.Query("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}

	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, time.Local)
	tGiven = tGiven.In(time.UTC)
	return s.toJulianTimeV1(c, tGiven, precision)
}

func (s *Server) toJulianTimeV1(c *fiber.Ctx, tGiven time.Time, precision int) error {
	resp := JulianTimeResp{}
	resp.CivilDate = tGiven.String()
	resp.CivilDateTimestamp = tGiven.Unix()
	resp.JulianDate = toFixed(jt.ToJulianDate(tGiven), precision)
	return c.JSON(resp)
}

func (s *Server) fromJulianTimeV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)

	jtStr := c.Query("jtime", "none")
	jtime, err := strconv.ParseFloat(jtStr, 64)
	if err != nil {
		c.Status(400)
		c.JSON("missing required parameter: 'jtime' (float)")
		return nil
	}

	t := jt.FromJulianDate(jtime, loc)

	precision := strToInt(c.Query("precision", "2"), 2, 0, 20)

	resp := JulianTimeResp{}
	resp.CivilDate = t.String()
	resp.CivilDateTimestamp = t.Unix()
	resp.JulianDate = toFixed(jtime, precision)

	return c.JSON(resp)
}
