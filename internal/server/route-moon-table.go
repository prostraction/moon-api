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
	table := s.moonCache.GenerateMoonTable(tGiven)

	if strings.ToLower(timeFormat) == "timestamp" {
		val := []moon.NearestPhaseTimestamp{}
		for i := range table.Elems {
			val = append(val, moon.NearestPhaseTimestamp{
				NewMoon:      table.Elems[i].NewMoon.Unix(),
				FirstQuarter: table.Elems[i].FirstQuarter.Unix(),
				FullMoon:     table.Elems[i].FullMoon.Unix(),
				LastQuarter:  table.Elems[i].FullMoon.Unix(),
			})
		}
		return c.JSON(val)
	} else if strings.ToLower(timeFormat) != "iso" {
		val := []moon.NearestPhaseString{}
		for i := range table.Elems {
			val = append(val, moon.NearestPhaseString{
				NewMoon:      table.Elems[i].NewMoon.Format(timeFormat),
				FirstQuarter: table.Elems[i].FirstQuarter.Format(timeFormat),
				FullMoon:     table.Elems[i].FullMoon.Format(timeFormat),
				LastQuarter:  table.Elems[i].FullMoon.Format(timeFormat),
			})
		}
		return c.JSON(val)
	}

	return c.JSON(table.Elems)
}
