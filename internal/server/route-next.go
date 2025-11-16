package server

import (
	jt "moon/pkg/julian-time"
	"moon/pkg/moon"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) moonNextMoonPhaseV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	tGiven = tGiven.In(loc)

	np := s.moonCache.FindNearestPhase(tGiven)

	format := c.Query("timeFormat", "ISO")
	if strings.ToLower(format) == "timestamp" {
		npt := moon.NearestPhaseTimestamp{
			NewMoon:      np.NewMoon.Unix(),
			FirstQuarter: np.FirstQuarter.Unix(),
			FullMoon:     np.FullMoon.Unix(),
			LastQuarter:  np.LastQuarter.Unix(),
		}
		return c.JSON(npt)
	} else if strings.ToLower(format) == "duration" {
		npd := moon.NearestPhaseTimestamp{
			NewMoon:      np.NewMoon.Unix() - tGiven.Unix(),
			FirstQuarter: np.FirstQuarter.Unix() - tGiven.Unix(),
			FullMoon:     np.FullMoon.Unix() - tGiven.Unix(),
			LastQuarter:  np.LastQuarter.Unix() - tGiven.Unix(),
		}
		return c.JSON(npd)
	} else if strings.ToLower(format) == "iso" {
		return c.JSON(np)
	} else {
		npc := moon.NearestPhaseString{
			NewMoon:      np.NewMoon.Format(format),
			FirstQuarter: np.FirstQuarter.Format(format),
			FullMoon:     np.FullMoon.Format(format),
			LastQuarter:  np.LastQuarter.Format(format),
		}
		return c.JSON(npc)
	}

}
