package server

import (
	jt "moon/pkg/julian-time"
	"strconv"
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
	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)

	tGiven := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableV1(c *fiber.Ctx, tGiven time.Time) error {
	resp := s.moonCache.GenerateMoonTable(tGiven)
	return c.JSON(resp.Elems)
}
