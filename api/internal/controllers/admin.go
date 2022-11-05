package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HookAdminEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.POST("/whoami", func(c echo.Context) error {
		user := c.Get("user")
		return c.JSON(http.StatusOK, user)
	})
}
