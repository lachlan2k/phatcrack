package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
			for _, userRole := range u.Roles {
				if userRole == RoleMFAEnrolled {
					userIsEnrolled = true
				}

				// Early exit if they're exempt
				if userRole == RoleMFAExempt {
					return next(c)
				}
			}

			userHasCompleted := u.HasCompletedMFA

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

func (a *AuthHandler) Middleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: middleware.AlgorithmHS256,
		SigningKey:    a.Secret,
		TokenLookup:   "cookie:" + TokenCookieName,
		Claims:        &AuthClaims{},
		Skipper:       a.shouldSkip,
	})
}

func (a *AuthHandler) RoleRestrictedMiddleware(allowedRoles []string, disallowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if a.shouldSkip(c) {
				return next(c)
			}

			user, ok := c.Get("user").(*jwt.Token)
			if user == nil || !ok {
				return echo.ErrUnauthorized
			}

			claims, ok := user.Claims.(*AuthClaims)
			if claims == nil || !ok {
				return echo.ErrUnauthorized
			}

			for _, disallowedRole := range disallowedRoles {
				for _, userRole := range claims.Roles {
					if disallowedRole == userRole {
						return echo.ErrUnauthorized
					}
				}
			}

			for _, allowedRole := range allowedRoles {
				for _, userRole := range claims.Roles {
					if allowedRole == userRole {
						return next(c)
					}
				}
			}

			return echo.ErrUnauthorized
		}
	}
}

func (a *AuthHandler) AdminOnlyMiddleware() echo.MiddlewareFunc {
	return a.RoleRestrictedMiddleware([]string{"admin"}, nil)
}
