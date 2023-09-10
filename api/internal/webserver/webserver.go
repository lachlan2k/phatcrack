package webserver

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/controllers"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func makeSessionHandler() auth.SessionHandler {
	return &auth.InMemorySessionHandler{
		SessionTimeout:     30 * time.Minute,
		SessionMaxLifetime: 4 * time.Hour,
	}
}

func Listen(port string) error {
	e := echo.New()

	validator := &util.RequestValidator{Validator: validator.New()}
	validator.Init()
	e.Validator = validator

	e.Use(makeLoggerMiddleware())
	e.Use(middleware.Recover())

	sessionHandler := makeSessionHandler()

	api := e.Group("/api/v1")

	// Agent auth is done separately in the controller, so it can go before auth middleware
	controllers.HookAgentHandlerEndpoints(api.Group("/agent-handler"))

	api.Use(auth.CreateHeaderAuthMiddleware())
	api.Use(sessionHandler.CreateMiddleware())

	api.Use(auth.EnforceAuthMiddleware([]string{
		"/api/v1/auth/login",
	}))

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// If a user has "requires_password_change" etc they need to be able to do that
	// Don't worry, the sessionhandler middleware is already enforcing auth
	controllers.HookAuthEndpoints(api.Group("/auth"), sessionHandler)

	api.Use(auth.EnforceMFAMiddleware(sessionHandler))

	api.Use(auth.RoleRestrictedMiddleware(
		[]string{auth.RoleAdmin, auth.RoleStandard},
		[]string{auth.RoleRequiresPasswordChange}, // disallowed
	))

	controllers.HookHashcatEndpoints(api.Group("/hashcat"))
	controllers.HookProjectEndpoints(api.Group("/project"))
	controllers.HookListsEndpoints(api.Group("/list"))
	controllers.HookHashlistEndpoints(api.Group("/hashlist"))
	controllers.HookAttackEndpoints(api.Group(("/attack")))
	controllers.HookAgentEndpoints(api.Group("/agent"))
	controllers.HookJobEndpoints(api.Group("/job"))
	controllers.HookAccountEndpoints(api.Group("/account"))
	controllers.HookConfigEndpoints(api.Group("/config"))

	adminAPI := api.Group("/admin")
	adminAPI.Use(auth.AdminOnlyMiddleware(sessionHandler))
	controllers.HookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
