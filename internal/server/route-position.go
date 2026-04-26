package server

import (
	jt "moon/pkg/julian-time"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) moonPositionMonthly(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)

	tNow := time.Now()
	year := StrToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 1, 9999)
	month := StrToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 1, 12)

	precision := StrToInt(c.Query("precision", "2"), 2, 0, 20)

	latStr := c.Query("latitude", "no-value")
	lonStr := c.Query("longitude", "no-value")
	locationCords := parseCoords(latStr, lonStr)

	timeFormat := c.Query("timeFormat", "ISO")

	if !locationCords.IsValid {
		e := ErrorPrintable{}
		e.Status = 400
		e.Message = "latitude and longitude are required for this method."
		return c.JSON(e)
	}

	if resp, err := s.positionCache.GetRisesMonthly(year, month, loc, precision, timeFormat, locationCords.Longitude, locationCords.Latitude); err == nil {
		return c.JSON(resp)
	} else {
		e := ErrorPrintable{}
		if strings.Contains(err.Error(), "400 Bad Request") {
			e.Status = 400
			c.Status(400)
		} else {
			e.Status = 500
			c.Status(500)
		}
		e.Message = err.Error()
		return c.JSON(e)
	}
}
