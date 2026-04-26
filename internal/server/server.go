package server

import (
	"moon/pkg/position"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	//moonCache     *moon.Cache
	positionCache *position.Cache
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:                      true,
		ServerHeader:                 "",
		CaseSensitive:                false,
		StrictRouting:                false,
		ReadTimeout:                  60 * time.Second,
		WriteTimeout:                 60 * time.Second,
		DisableKeepalive:             true,
		DisableStartupMessage:        true,
		DisablePreParseMultipartForm: true,
	})
	s.RegisterRoutes(app)
	return app
}

// RegisterRoutes attaches CORS, static handlers and all API routes to app.
// Split out from NewRouter so tests can construct a Prefork-less app.
func (s *Server) RegisterRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,HEAD,OPTIONS",
	}))

	// web vanilla JS because it's cool
	app.Static("/styles", "web/styles/")
	app.Static("/assets", "web/assets/")
	app.Static("/", "web")

	app.Get("/v1/moonTableYear", s.moonTableYearV1)
	app.Get("/v1/moonTableCurrent", s.moonTableCurrentV1)
	app.Get("/api/v1/moonTableYear", s.moonTableYearV1)
	app.Get("/api/v1/moonTableCurrent", s.moonTableCurrentV1)

	app.Get("/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/v1/moonPhaseDate", s.moonPhaseDatetV1)
	app.Get("/api/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/api/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/api/v1/moonPhaseDate", s.moonPhaseDatetV1)

	app.Get("/v1/moonPositionMonthly", s.moonPositionMonthly)
	app.Get("/api/v1/moonPositionMonthly", s.moonPositionMonthly)

	app.Get("/v1/nextMoonPhase", s.moonNextMoonPhaseV1)
	app.Get("/v1/nextMoonDay", s.moonNextMoonDayV1)
	app.Get("/api/v1/nextMoonPhase", s.moonNextMoonPhaseV1)
	app.Get("/api/v1/nextMoonDay", s.moonNextMoonDayV1)

	app.Get("/v1/toJulianTimeByDate", s.toJulianTimeByDateV1)
	app.Get("/v1/toJulianTimeByTimestamp", s.toJulianTimeByTimestampV1)
	app.Get("/v1/fromJulianTime", s.fromJulianTimeV1)

	app.Get("/api/v1/toJulianTimeByDate", s.toJulianTimeByDateV1)
	app.Get("/api/v1/toJulianTimeByTimestamp", s.toJulianTimeByTimestampV1)
	app.Get("/api/v1/fromJulianTime", s.fromJulianTimeV1)

	app.Get("/v1/version", s.versionV1)
	app.Get("/api/v1/version", s.versionV1)

	if s.positionCache == nil {
		s.positionCache = new(position.Cache)
	}
}

func (s *Server) versionV1(c *fiber.Ctx) error {
	return c.JSON("1.2.2")
}
