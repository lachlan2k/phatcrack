package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/resources"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookResourceEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api.GET("/hashtypes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, apitypes.HashTypesDTO{
			HashTypes: resources.GetHashTypeMap(),
		})
	})
}
