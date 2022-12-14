package webserver

import (
	"crypto/rand"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/controllers"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func Listen(port string) error {
	e := echo.New()

	e.Validator = &util.RequestValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if os.Getenv("CORS_ALLOWED") != "" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{os.Getenv("CORS_ALLOWED")},
			// AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			// AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowCredentials: true,
		}))
	}

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
		Secret: jwtKey,
		WhitelistPaths: []string{
			"/api/v1/agent/handle/ws",
			"/api/v1/auth/login",
			"/api/v1/auth/refresh",
			"/api/v1/auth/whoami",
		},
	}

	api := e.Group("/api/v1")

	api.Use(authHandler.Middleware())

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	controllers.HookAuthEndpoints(api.Group("/auth"), &authHandler)
	controllers.HookJobEndpoints(api.Group("/job"))
	controllers.HookResourceEndpoints(api.Group("/resources"))
	controllers.HookAgentEndpoints(api.Group("/agent"))

	adminAPI := api.Group("/admin")
	adminAPI.Use(authHandler.AdminOnlyMiddleware())

	controllers.HookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
