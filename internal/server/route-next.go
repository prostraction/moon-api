package server

import (
	jt "moon/pkg/julian-time"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) moonNextMoonPhaseV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	tGiven = tGiven.In(loc)

	np := s.moonCache.FindNearestPhase(tGiven)

	return c.JSON(np)
}
