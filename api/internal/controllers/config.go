package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/config"
)

func HookConfigEndpoints(api *echo.Group) {
	api.GET("/public", func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.Get().ToPublicDTO())
	})
}
