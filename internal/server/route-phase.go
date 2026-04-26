package server

import (
	jt "moon/pkg/julian-time"
	"moon/pkg/moon"
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
	utc := c.Query("utc", "UTC+0")
	loc, err := jt.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Trace(err)
	}
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	tGiven = tGiven.In(loc)
	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	timeFormat := c.Query("timeFormat", "ISO")

	return s.moonPhaseV1(c, tGiven, precision, locationCords, timeFormat)
}

func (s *Server) moonPhaseTimestampV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC+0")
	loc, err := jt.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Trace(err)
	}

	tStr := c.Query("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}

	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, time.Local).In(loc)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	timeFormat := c.Query("timeFormat", "ISO")

	return s.moonPhaseV1(c, tGiven, precision, locationCords, timeFormat)
}

func (s *Server) moonPhaseDatetV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC+0")
	loc, err := jt.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Trace(err)
	}

	tNow := time.Now()

	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := StrToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)
	day := StrToInt(c.Query("day", strconv.Itoa(tNow.Day())), tNow.Day(), 1, 31)

	err = IsValidDate(year, month, day)
	if err != nil {
		c.Status(400)
		errPrintable := ErrorPrintable{
			Status:  400,
			Message: "Validation error: " + err.Error(),
		}
		return c.JSON(errPrintable)
	}

	hour := StrToInt(c.Query("hour", strconv.Itoa(tNow.Hour())), tNow.Hour(), 0, 23)
	minute := StrToInt(c.Query("minute", strconv.Itoa(tNow.Minute())), tNow.Minute(), 0, 59)
	second := StrToInt(c.Query("second", strconv.Itoa(tNow.Second())), tNow.Second(), 0, 59)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	timeFormat := c.Query("timeFormat", "ISO")

	tGiven := time.Date(year, jt.GetMonth(month), day, hour, minute, second, 0, time.Local)
	tGiven = tGiven.In(loc)
	return s.moonPhaseV1(c, tGiven, precision, locationCords, timeFormat)
}

func (s *Server) moonPhaseV1(c *fiber.Ctx, tGiven time.Time, precision int, locationCords Coordinates, timeFormat string) error {
	var err error

	lang := c.Query("lang", "en")
	utc := c.Query("utc", "UTC+0")

	loc, err := jt.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Trace(err)
	}

	resp := MoonPhaseResponse{}
	resp.BeginDay = new(MoonStat)
	resp.CurrentState = new(MoonStat)
	resp.EndDay = new(MoonStat)

	moonTable := moon.CreateMoonTable(tGiven)

	// moon days
	day := moon.CurrentMoonDays(tGiven, loc, moonTable)

	resp.BeginDay.MoonDays = ToFixed(day.Begin.Minutes()/jt.Fminute, precision)
	resp.CurrentState.MoonDays = ToFixed(day.Current.Minutes()/jt.Fminute, precision)
	resp.EndDay.MoonDays = ToFixed(day.End.Minutes()/jt.Fminute, precision)

	// phase && illum
	phase := phase.CurrentMoonPhase(tGiven, lang)

	resp.BeginDay.Illumination = ToFixed(phase.Illumination.BeginDay*100, precision)
	resp.CurrentState.Illumination = ToFixed(phase.Illumination.Current*100, precision)
	resp.EndDay.Illumination = ToFixed(phase.Illumination.EndDay*100, precision)

	resp.BeginDay.Phase = phase.BeginDay
	resp.CurrentState.Phase = phase.Current
	resp.EndDay.Phase = phase.EndDay

	// zodiac TO DO refactor
	resp.ZodiacDetailed, resp.BeginDay.Zodiac, resp.CurrentState.Zodiac, resp.EndDay.Zodiac = zodiac.CurrentZodiacs(tGiven, loc, lang, timeFormat, moonTable.Elems)

	if locationCords.IsValid {
		resp.MoonRiseAndSet, err = s.positionCache.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision, timeFormat, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}
		resp.MoonDaysDetailed = moon.MoonDetailed(tGiven, loc, lang, timeFormat, locationCords.Longitude, locationCords.Latitude)

		newT := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, tGiven.Location())
		resp.BeginDay.Position, err = pos.GetMoonPosition(newT, newT.Location(), precision, timeFormat, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}

		resp.CurrentState.Position, err = pos.GetMoonPosition(tGiven, tGiven.Location(), precision, timeFormat, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}

		resp.EndDay.Position, err = pos.GetMoonPosition(newT.AddDate(0, 0, 1), newT.AddDate(0, 0, 1).Location(), precision, timeFormat, locationCords.Longitude, locationCords.Latitude)
		if err != nil {
			log.Error(err)
		}
	} else {
		resp.MoonRiseAndSet, err = s.positionCache.GetRisesDay(tGiven.Year(), int(tGiven.Month()), tGiven.Day(), tGiven.Location(), precision, timeFormat)
	}

	if err != nil && err.Error() != "no location prodived" {
		log.Error(err.Error())
	}

	return c.JSON(resp)
}
