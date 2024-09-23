package webserver

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/controllers"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/sirupsen/logrus"
)

func makeSessionHandler() auth.SessionHandler {
	return &auth.InMemorySessionHandler{
		SessionTimeout:     30 * time.Minute,
		SessionMaxLifetime: 4 * time.Hour,
	}
}

func Listen(baseURL string,port string) error {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin != "" && origin != baseURL {
				return echo.NewHTTPError(http.StatusForbidden, "Origin not allowed")
			}
			return next(c)
		}
	})

	validator := util.NewRequestValidator()
	e.Validator = validator

	e.Use(makeLoggerMiddleware())
	e.Use(middleware.Recover())

	sessionHandler := makeSessionHandler()

	api := e.Group("/api/v1")

	// Agent auth is done separately in the controller, so it can go before auth middleware
	controllers.HookAgentHandlerEndpoints(api.Group("/agent-handler"))

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	if strings.TrimSpace(os.Getenv("E2E_TEST_ENABLE_FIXTURES")) != "" && len(strings.TrimSpace(os.Getenv("E2E_TEST_FIXTURE_KEY"))) > 1 {
		logrus.Info("Testing fixtures are enabled -- do not run this in production")
		controllers.HookE2EEndpoints(api.Group("/e2e"), strings.TrimSpace(os.Getenv("E2E_TEST_FIXTURE_KEY")))
	}

	api.Use(auth.CreateHeaderAuthMiddleware())
	api.Use(sessionHandler.CreateMiddleware())

	api.Use(auth.EnforceAuthMiddleware(auth.EnforceAuthArgs{
		BypassPaths: []string{
			"/api/v1/config/public",
			"/api/v1/auth/login/credentials",
			"/api/v1/auth/login/oidc/start",
			"/api/v1/auth/login/oidc/callback",
		},
	}))

	// If a user has "requires_password_change" etc they need to be able to do that
	// Don't worry, the sessionhandler middleware is already enforcing auth
	controllers.HookAuthEndpoints(api.Group("/auth"), sessionHandler)

	// Config endpoints are publicly accessible (as they include things like what auth options should be listed)
	controllers.HookConfigEndpoints(api.Group("/config"))

	api.Use(auth.EnforceMFAMiddleware(sessionHandler))

	api.Use(auth.RoleRestrictedMiddleware(
		[]string{roles.UserRoleAdmin, roles.UserRoleStandard},
		[]string{roles.UserRoleRequiresPasswordChange}, // disallowed
	))


	controllers.HookHashcatEndpoints(api.Group("/hashcat"))
	controllers.HookProjectEndpoints(api.Group("/project"))
	controllers.HookListsEndpoints(api.Group("/listfiles"))
	controllers.HookHashlistEndpoints(api.Group("/hashlist"))
	controllers.HookAttackEndpoints(api.Group(("/attack")))
	controllers.HookAttackTemplateEndpoints(api.Group(("/attack-template")))
	controllers.HookAgentEndpoints(api.Group("/agent"))
	controllers.HookJobEndpoints(api.Group("/job"))
	controllers.HookAccountEndpoints(api.Group("/account"))
	controllers.HookUserEndpoints(api.Group("/user"))
	controllers.HookPotfileEndpoints(api.Group("/potfile"))

	adminAPI := api.Group("/admin")
	adminAPI.Use(auth.AdminOnlyMiddleware(sessionHandler))
	controllers.HookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
