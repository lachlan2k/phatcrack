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
		},
	}

	// Slightly annoying, the auth middleware, by default, uses a 400 error when the auth is missing
	// We want a 401
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if err == middleware.ErrJWTMissing {
			c.Error(echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
			return
		}
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}

	api := e.Group("/api/v1")

	// Agent auth is done separately in the controller, so it can go before auth middleware
	controllers.HookAgentEndpoints(api.Group("/agent"))

	api.Use(authHandler.Middleware())

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// If a user has "requires_password_change" etc they need to be able to do that
	// Don't worry, the authHandler.Middleware() is already enforcing auth
	controllers.HookAuthEndpoints(api.Group("/auth"), &authHandler)

	api.Use(authHandler.EnforceMFA())

	api.Use(authHandler.RoleRestrictedMiddleware(
		[]string{auth.RoleAdmin, auth.RoleStandard},
		[]string{auth.RoleRequiresPasswordChange}, // disallowed
	))

	controllers.HookHashcatEndpoints(api.Group("/hashcat"))
	controllers.HookProjectEndpoints(api.Group("/project"))
	controllers.HookListsEndpoints(api.Group("/list"))
	controllers.HookHashlistEndpoints(api.Group("/hashlist"))
	controllers.HookAttackEndpoints(api.Group(("/attack")))
	controllers.HookJobEndpoints(api.Group("/job"))

	adminAPI := api.Group("/admin")
	adminAPI.Use(authHandler.AdminOnlyMiddleware())

	controllers.HookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
