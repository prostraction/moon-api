package server

import (
	jt "moon/pkg/julian-time"
	"moon/pkg/moon"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*    MOON TABLE    */
func (s *Server) moonTableCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	timeFormat := c.Query("timeFormat", "ISO")

	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	return s.moonTableV1(c, timeFormat, tGiven)
}

func (s *Server) moonTableYearV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	timeFormat := c.Query("timeFormat", "ISO")

	tNow := time.Now()
	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)

	tGiven := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	return s.moonTableV1(c, timeFormat, tGiven)
}

func (s *Server) moonTableV1(c *fiber.Ctx, timeFormat string, tGiven time.Time) error {
	moonTable := moon.CreateMoonTable(tGiven)

	if strings.ToLower(timeFormat) == "timestamp" {
		val := []moon.NearestPhaseTimestamp{}
		for i := range moonTable.Elems {
			val = append(val, moon.NearestPhaseTimestamp{
				NewMoon:      moonTable.Elems[i].NewMoon.Unix(),
				FirstQuarter: moonTable.Elems[i].FirstQuarter.Unix(),
				FullMoon:     moonTable.Elems[i].FullMoon.Unix(),
				LastQuarter:  moonTable.Elems[i].FullMoon.Unix(),
			})
		}
		return c.JSON(val)
	} else if strings.ToLower(timeFormat) != "iso" {
		val := []moon.NearestPhaseString{}
		for i := range moonTable.Elems {
			val = append(val, moon.NearestPhaseString{
				NewMoon:      moonTable.Elems[i].NewMoon.Format(timeFormat),
				FirstQuarter: moonTable.Elems[i].FirstQuarter.Format(timeFormat),
				FullMoon:     moonTable.Elems[i].FullMoon.Format(timeFormat),
				LastQuarter:  moonTable.Elems[i].FullMoon.Format(timeFormat),
			})
		}
		return c.JSON(val)
	}

	return c.JSON(moonTable.Elems)
}
