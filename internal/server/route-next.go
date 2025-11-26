package server

import (
	jt "moon/pkg/julian-time"
	"moon/pkg/moon"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Server) moonNextMoonPhaseV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Local)
	tGiven = tGiven.In(loc)

	moonTable := moon.CreateMoonTable(tGiven)
	np := moon.FindNearestPhase(tGiven, moonTable)

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

func (s *Server) moonNextMoonDayV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := jt.SetTimezoneLocFromString(utc)
	day := c.Query("day", "not-set")
	if day == "not-set" {
		e := ErrorPrintable{Status: 400, Message: "day is required for this method."}
		return c.JSON(e)
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		e := ErrorPrintable{Status: 400, Message: "day is required to be int."}
		return c.JSON(e)
	}
	if dayInt < 0 || dayInt > 30 {
		e := ErrorPrintable{Status: 400, Message: "day is required to be [0, 30]. No other days allowed!"}
		return c.JSON(e)
	}

	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Local)
	tGiven = tGiven.In(loc)

	moonTable := moon.CreateMoonTable(tGiven)
	resp, err := moon.SearchMoonDay(tGiven, moonTable, dayInt)
	if err != nil {
		log.Error(err)
		e := ErrorPrintable{Status: 500, Message: "bad things happen: " + err.Error()}
		return c.JSON(e)
	}

	format := c.Query("timeFormat", "ISO")
	if strings.ToLower(format) == "timestamp" {
		resp := moon.SeachMoonDayRespTimestamp{
			From: resp.From.Unix(),
			To:   resp.To.Unix(),
		}
		return c.JSON(resp)
	} else if strings.ToLower(format) == "duration" {
		resp := moon.SeachMoonDayRespTimestamp{
			From: resp.From.Unix() - tGiven.Unix(),
			To:   resp.To.Unix() - tGiven.Unix(),
		}
		return c.JSON(resp)
	} else if strings.ToLower(format) == "iso" {
		return c.JSON(resp)
	} else {
		resp := moon.SeachMoonDayRespString{
			From: resp.From.Format(format),
			To:   resp.To.Format(format),
		}
		return c.JSON(resp)
	}
}
