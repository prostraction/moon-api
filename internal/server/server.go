package server

import (
	"moon/pkg/moon"
	"moon/pkg/position"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	moonCache     *moon.Cache
	positionCache *position.Cache
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Fiber",
		CaseSensitive: false,
		StrictRouting: true,
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,HEAD,OPTIONS",
	}))

	// web vanilla JS because it's cool
	app.Static("/styles", "web/styles/")
	app.Static("/assets", "web/assets/")
	app.Static("/", "web")

	// moon table for year:
	// new moon, first quarter, full moon, last quarter
	// maybe rename moonTable -> phaseTable
	app.Get("/v1/moonTableYear", s.moonTableYearV1)
	app.Get("/v1/moonTableCurrent", s.moonTableCurrentV1)
	app.Get("/api/v1/moonTableYear", s.moonTableYearV1)
	app.Get("/api/v1/moonTableCurrent", s.moonTableCurrentV1)

	// maybe rename moonPhase -> phase later

	// moon phase for day:
	// - begin, current, end of the day:
	//	- moon days, illumintaion, phase, zodiac
	// - moon days by time
	// - zodiacs by time
	// - rise, set and meridian:
	//	- time, direction, azimuth, exists
	app.Get("/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/v1/moonPhaseDate", s.moonPhaseDatetV1)
	app.Get("/api/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/api/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/api/v1/moonPhaseDate", s.moonPhaseDatetV1)

	// per month
	app.Get("/v1/moonPositionMonthly", s.moonPositionMonthly)
	app.Get("/api/v1/moonPositionMonthly", s.moonPositionMonthly)
	//app.Get("/v1/moonRiseSetCalendar")
	//app.Get("/v1/moonPhaseCalendar")
	//app.Get("/v1/moonZodiacCalendar")
	//app.Get("/v1/moonMonthCalendar") -- all combined? think more....

	// eclipses
	//app.Get("/v1/moonEclipseYear")
	//app.Get("/v1/moonEclipseCalendar")
	//app.Get("/v1/sunEclipseYear")
	//app.Get("/v1/sunEclipseCalendar")

	// methods when?
	app.Get("/v1/nextMoonPhase", s.moonNextMoonPhaseV1)
	app.Get("/api/v1/nextMoonPhase", s.moonNextMoonPhaseV1)

	//app.Get("/v1/nextMoonEclipse")
	//app.Get("/v1/nextSunEclipse")
	//app.Get("/v1/nextMoonZodiac")
	//app.Get("/v1/nextMoonSet")
	//app.Get("/v1/nextMoonRise")

	// jtime methods:
	app.Get("/v1/toJulianTimeByDate", s.toJulianTimeByDateV1)
	app.Get("/v1/toJulianTimeByTimestamp", s.toJulianTimeByTimestampV1)
	app.Get("/v1/fromJulianTime", s.fromJulianTimeV1)

	app.Get("/api/v1/toJulianTimeByDate", s.toJulianTimeByDateV1)
	app.Get("/api/v1/toJulianTimeByTimestamp", s.toJulianTimeByTimestampV1)
	app.Get("/api/v1/fromJulianTime", s.fromJulianTimeV1)

	// some kind of faq
	//app.Get("/v1")

	// just print version
	app.Get("/v1/version", s.versionV1)
	app.Get("/api/v1/version", s.versionV1)

	s.moonCache = new(moon.Cache)
	s.positionCache = new(position.Cache)
	return app
}

func (s *Server) versionV1(c *fiber.Ctx) error {
	return c.JSON("1.2.0rc6")
}
