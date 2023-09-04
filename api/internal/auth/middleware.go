package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/config"
)

func EnforceMFAMiddleware(s SessionHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, sess := UserAndSessFromReq(c)
			if user == nil {
				return echo.ErrUnauthorized
			}

			userIsEnrolled := false
			for _, userRole := range user.Roles {
				if userRole == RoleMFAEnrolled {
					userIsEnrolled = true
				}

				// Early exit if they're exempt
				if userRole == RoleMFAExempt {
					return next(c)
				}
			}

			userHasCompleted := sess.HasCompletedMFA

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

func EnforceAuthMiddleware(bypassPaths []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			for _, bypassPath := range bypassPaths {
				if path == bypassPath {
					return next(c)
				}
			}

			if !AuthIsValid(c) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Login required")
			}

			return next(c)
		}
	}
}

func RoleRestrictedMiddleware(h SessionHandler, allowedRoles []string, disallowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := UserAndSessFromReq(c)

			if user == nil {
				return echo.ErrUnauthorized
			}

			for _, disallowedRole := range disallowedRoles {
				for _, userRole := range user.Roles {
					if disallowedRole == userRole {
						return echo.ErrUnauthorized
					}
				}
			}

			for _, allowedRole := range allowedRoles {
				for _, userRole := range user.Roles {
					if allowedRole == userRole {
						return next(c)
					}
				}
			}

			return echo.ErrUnauthorized
		}
	}
}

func AdminOnlyMiddleware(h SessionHandler) echo.MiddlewareFunc {
	return RoleRestrictedMiddleware(h, []string{"admin"}, nil)
}
