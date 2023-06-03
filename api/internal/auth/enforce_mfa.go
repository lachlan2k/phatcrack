package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/config"
)

func (a *AuthHandler) EnforceMFA() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if a.shouldSkip(c) {
				return next(c)
			}

			if !config.Get().IsMFARequired {
				return next(c)
			}

			u, err := ClaimsFromReq(c)
			if err != nil {
				return echo.ErrUnauthorized
			}

			for _, userRole := range u.Roles {
				if userRole == RoleMFAEnrolled || userRole == RoleMFAExempt {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "MFA not yet configured")
		}
	}
}
