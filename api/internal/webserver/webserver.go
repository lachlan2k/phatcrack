package webserver

import (
	"crypto/rand"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func Listen(port string) error {
	e := echo.New()

	e.Validator = &util.RequestValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if len(jwtKey) == 0 {
		keyBuf := make([]byte, 32)
		rand.Reader.Read(keyBuf)
		jwtKey = keyBuf[:]
		e.Logger.Printf("Generating jwtKey")
	} else {
		e.Logger.Printf("Using JWT_KEY from environment")
	}

	authHandler := auth.AuthHandler{
		Secret:         jwtKey,
		WhitelistPaths: []string{"/api/v1/agent/handler/ws", "/api/v1/auth/login"},
	}

	e.Use(authHandler.Middleware())

	api := e.Group("/api/v1")

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	hookAuthEndpoints(api.Group("/auth"), &authHandler)
	hookJobEndpoints(api.Group("/job"))
	hookAgentEndpoints(api.Group("/agent"))

	adminAPI := api.Group("/admin")
	adminAPI.Use(authHandler.AdminOnlyMiddleware())
	hookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
