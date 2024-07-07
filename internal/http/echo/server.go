package server

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	HTTPServerPort string `env:"SERVER_PORT,default=8123"`
	HideBanner     bool   `env:"HIDE_BANNER, default=false"`
}

type Server interface {
	Start(address string) error
	Shutdown(ctx context.Context) error
}

var _ Server = (*echo.Echo)(nil)

func New(c *ServerConfig, h *Handler) Server {
	e := echo.New()
	e.HideBanner = c.HideBanner
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
		Skipper: func(c echo.Context) bool {
			return strings.EqualFold(c.Request().URL.Path, "/")
		},
	}))
	e.GET("/", ping)
	e.POST("/messages", h.message)

	return e
}
