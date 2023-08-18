package webserver

import (
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/controllers"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	log "github.com/sirupsen/logrus"
)

func makeSessionHandler() auth.SessionHandler {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	if len(jwtKey) > 0 && 1 == 0 {
		return &auth.JWTSessionHandler{
			Secret: jwtKey,
			WhitelistPaths: []string{
				"/api/v1/auth/login",
			},
			SessionLifetime: 10 * time.Minute,
		}
	}

	return &auth.InMemorySessionHandler{
		WhitelistPaths: []string{
			"/api/v1/auth/login",
		},
		SessionLifetime: 10 * time.Minute,
	}
}

func Listen(port string) error {
	e := echo.New()

	validator := &util.RequestValidator{Validator: validator.New()}
	validator.Init()
	e.Validator = validator

	// e.Use(middleware.LoggerWithConfig(
	// 	middleware.LoggerConfig{
	// 		Skipper: func(c echo.Context) bool {
	// 			return false
	// 			if c.Response().Status == 200 {
	// 				return true
	// 			}
	// 			return false
	// 		},
	// 	},
	// ))
	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{

			LogLatency:       true,
			LogRemoteIP:      true,
			LogMethod:        true,
			LogURI:           true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogContentLength: true,
			LogResponseSize:  true,

			LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
				fields := log.Fields{
					"latency_ms":     values.Latency.Milliseconds(),
					"remote_ip":      values.RemoteIP,
					"method":         values.Method,
					"URI":            values.URI,
					"user_agent":     values.UserAgent,
					"status":         values.Status,
					"content_length": values.ContentLength,
					"response_size":  values.ResponseSize,
				}

				if values.Error != nil {
					log.WithError(values.Error).WithFields(fields).Error("request error")
					return nil
				}

				if values.Status >= 500 && values.Status <= 599 {
					log.WithFields(fields).Error("request error")
					return nil
				}

				if values.Status == 400 || values.Status == 403 {
					log.WithFields(fields).Warn("bad request")
					return nil
				}

				log.WithFields(fields).Info("request")
				return nil
			},
		},
	))
	e.Use(middleware.Recover())

	// Slightly annoying, the auth middleware, by default, uses a 400 error when the auth is missing
	// We want a 401
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if err == middleware.ErrJWTMissing {
			c.Error(echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
			return
		}
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}

	sessionHandler := makeSessionHandler()

	api := e.Group("/api/v1")

	// Agent auth is done separately in the controller, so it can go before auth middleware
	controllers.HookAgentHandlerEndpoints(api.Group("/agent-handler"))

	api.Use(sessionHandler.CreateMiddleware())

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// If a user has "requires_password_change" etc they need to be able to do that
	// Don't worry, the sessionhandler middleware is already enforcing auth
	controllers.HookAuthEndpoints(api.Group("/auth"), sessionHandler)

	api.Use(auth.EnforceMFAMiddleware(sessionHandler))

	api.Use(auth.RoleRestrictedMiddleware(
		sessionHandler,
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

	adminAPI := api.Group("/admin")
	adminAPI.Use(auth.AdminOnlyMiddleware(sessionHandler))

	controllers.HookAdminEndpoints(adminAPI)

	return e.Start(":" + port)
}
