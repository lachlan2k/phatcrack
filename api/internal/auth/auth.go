package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lachlan2k/phatcrack/api/internal/db"
)

const TokenCookieName = "auth"
const TokenLifetime = 15 * time.Minute

type AuthHandler struct {
	Secret         []byte
	WhitelistPaths []string
}

type UserClaims struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type AuthClaims struct {
	UserClaims
	jwt.StandardClaims
}

func UserToClaims(user *db.User) *AuthClaims {
	return &AuthClaims{
		UserClaims: UserClaims{
			ID:       user.ID.String(),
			Username: user.Username,
			Roles:    user.Roles,
		},
	}
}

func ClaimsFromReq(c echo.Context) (*AuthClaims, error) {
	u, ok := c.Get("user").(*jwt.Token)
	if u == nil || !ok {
		return nil, fmt.Errorf("couldn't cast token %v", c.Get("user"))
	}
	claims, ok := u.Claims.(*AuthClaims)
	if claims == nil || !ok {
		return nil, fmt.Errorf("couldn't cast token claims %v", u.Claims)
	}
	return claims, nil
}

func (a *AuthHandler) SignJwt(claims *AuthClaims) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(TokenLifetime)

	claims.StandardClaims = jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expires.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.Secret)
	if err != nil {
		return "", now, err
	}

	return tokenString, expires, nil
}

func (a *AuthHandler) SignAndSetJWT(c echo.Context, claims *AuthClaims) error {
	token, expires, err := a.SignJwt(claims)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https",
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func (a *AuthHandler) shouldSkip(c echo.Context) bool {
	path := c.Request().URL.Path
	for _, bypassPath := range a.WhitelistPaths {
		if path == bypassPath {
			return true
		}
	}
	return false
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
