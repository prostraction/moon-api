package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	jt "moon/pkg/julian-time"
)

/* Julian Time methods */
func (s *Server) toJulianTimeByDateV1(c *fiber.Ctx) error {
	tNow := time.Now()
	tNow = tNow.In(time.UTC)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := StrToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)
	day := StrToInt(c.Query("day", strconv.Itoa(int(tNow.Day()))), int(tNow.Day()), 1, 31)

	timeFormat := c.Query("timeFormat", "ISO")

	err := IsValidDate(year, month, day)
	if err != nil {
		c.Status(400)
		errPrintable := ErrorPrintable{
			Status:  400,
			Message: "Validation error: " + err.Error(),
		}
		return c.JSON(errPrintable)
	}

	hour := StrToInt(c.Query("hour", strconv.Itoa(int(tNow.Hour()))), int(tNow.Hour()), 0, 23)
	minute := StrToInt(c.Query("minute", strconv.Itoa(int(tNow.Minute()))), int(tNow.Minute()), 0, 59)
	second := StrToInt(c.Query("second", strconv.Itoa(int(tNow.Second()))), int(tNow.Second()), 0, 59)

	tGiven := time.Date(year, jt.GetMonth(month), day, hour, minute, second, 0, time.UTC)

	resp := s.toJulianTimeV1(c, tGiven, timeFormat, precision)
	return c.JSON(resp)
}

func (s *Server) toJulianTimeByTimestampV1(c *fiber.Ctx) error {
	tStr := c.Query("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}

	timeFormat := c.Query("timeFormat", "ISO")
	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, time.Local)
	tGiven = tGiven.In(time.UTC)

	resp := s.toJulianTimeV1(c, tGiven, timeFormat, precision)
	return c.JSON(resp)
}

func (s *Server) toJulianTimeV1(c *fiber.Ctx, tGiven time.Time, timeFormat string, precision int) JulianTimeResp {
	resp := JulianTimeResp{}

	var t any
	if strings.ToLower(timeFormat) == "timestamp" {
		t = tGiven.Unix()
	} else if strings.ToLower(timeFormat) != "iso" {
		t = tGiven.Format(timeFormat)
	} else {
		t = tGiven
	}

	resp.CivilDate = &t
	resp.JulianDate = ToFixed(jt.ToJulianDate(tGiven), precision)
	return resp
}

func (s *Server) fromJulianTimeV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)

	jtStr := c.Query("jtime", "none")
	jtime, err := strconv.ParseFloat(jtStr, 64)
	if err != nil {
		c.Status(400)
		errPrintable := ErrorPrintable{
			Status:  400,
			Message: "missing required parameter: 'jtime' (float)",
		}
		return c.JSON(errPrintable)
	}

	t := jt.FromJulianDate(jtime, loc)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)
	timeFormat := c.Query("timeFormat", "ISO")

	var tFormat any
	if strings.ToLower(timeFormat) == "timestamp" {
		tFormat = t.Unix()
	} else if strings.ToLower(timeFormat) != "iso" {
		tFormat = t.Format(timeFormat)
	} else {
		tFormat = t
	}

	resp := JulianTimeResp{}
	resp.CivilDate = &tFormat
	resp.JulianDate = ToFixed(jtime, precision)

	return c.JSON(resp)
}
