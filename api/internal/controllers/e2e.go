package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func HookE2EEndpoints(g *echo.Group, apiKey string) {
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-E2E-Key") == apiKey {
				return next(c)
			}

			return echo.ErrUnauthorized
		}
	})

	g.POST("/wipe", func(c echo.Context) error {
		AuditLog(c, nil, "Wiping database for e2e purposes")

		err := db.WipeEverything()
		if err != nil {
			return util.ServerError("Failed to wipe database", err)
		}

		config.Update(func(rc *config.RuntimeConfig) error {
			*rc = config.MakeDefaultConfig()
			return nil
		})

		return c.JSON(http.StatusOK, "ok")
	})
}
