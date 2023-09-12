package auth

import (
	"regexp"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	log "github.com/sirupsen/logrus"
)

func CreateHeaderAuthMiddleware() echo.MiddlewareFunc {
	headerRe := regexp.MustCompile(`(?i)\s*bearer\s+(.+)\s*`)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Another auth validator has already validated this request
			if AuthIsValid(c) {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			match := headerRe.FindStringSubmatch(authHeader)
			if match == nil || len(match) < 2 {
				return next(c)
			}

			token := match[1]
			user, err := db.GetServiceAccountByAPIKey(token)
			if err == db.ErrNotFound {
				return echo.ErrUnauthorized
			}
			if err != nil {
				return util.ServerError("Failed to check API key", err)
			}
			if !slices.Contains(user.Roles, roles.RoleServiceAccount) {
				log.WithField("user_dto", user.ToDTO()).Warn("Request sucessfully authorized by bearer token, but account isn't a service account")
				return echo.ErrUnauthorized
			}

			// Create dummy session entry
			c.Set(sessionContextKey, SessionData{
				UserID:              user.ID.String(),
				HasCompletedMFA:     false,
				WebAuthnSession:     nil,
				PendingWebAuthnUser: nil,
			})

			c.Set(sessionUserContextKey, user)
			c.Set(sessionAuthIsValidContextKey, true)

			return next(c)
		}
	}
}
