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

			u, err := ClaimsFromReq(c)
			if err != nil {
				return echo.ErrUnauthorized
			}

			userIsEnrolled := false
			userHasCompleted := false
			for _, userRole := range u.Roles {
				if userRole == RoleMFAEnrolled {
					userIsEnrolled = true
				}

				if userRole == RoleMFACompleted {
					userHasCompleted = true
				}

				// Early exit if they're exempt
				if userRole == RoleMFAExempt {
					return next(c)
				}
			}

			if config.Get().IsMFARequired {
				if !userIsEnrolled {
					return echo.NewHTTPError(http.StatusForbidden, "MFA not yet enrolled")
				}
			}

			// Even if we don't globally enforce MFA, we need to enforce it if the user has chosen to configure it themselves
			if userIsEnrolled {
				if !userHasCompleted {
					return echo.NewHTTPError(http.StatusForbidden, "MFA not yet completed")
				}
			}

			return next(c)
		}
	}
}
