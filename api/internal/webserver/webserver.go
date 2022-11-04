package webserver

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func Listen(port string) error {
	e := echo.New()

	e.Validator = &util.RequestValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	hookAgentEndpoints(api.Group("/agent"))
	hookJobEndpoints(api.Group("/job"))

	return e.Start(":" + port)
}
