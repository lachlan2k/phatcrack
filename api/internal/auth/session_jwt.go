package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type JWTSessionHandler struct {
	CookieName      string
	SessionLifetime time.Duration
	Secret          []byte
	WhitelistPaths  []string
}

type JWTClaims struct {
	SessionData
	jwt.StandardClaims
}

func (j *JWTSessionHandler) CreateMiddleware() echo.MiddlewareFunc {
	if j.CookieName == "" {
		j.CookieName = "auth"
	}

	jwtMw := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: middleware.AlgorithmHS256,
		SigningKey:    j.Secret,
		TokenLookup:   "cookie:" + j.CookieName,
		Claims:        &JWTClaims{},
		Skipper:       j.shouldSkip,
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return jwtMw(func(c echo.Context) error {
			u, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return next(c)
			}
			claims, ok := u.Claims.(*JWTClaims)
			if !ok {
				return next(c)
			}

			c.Set(sessionContextKey, claims.SessionData)
			return next(c)
		})
	}
}

func (j *JWTSessionHandler) Start(c echo.Context, sess SessionData) error {
	now := time.Now()
	expires := now.Add(j.SessionLifetime)

	claims := JWTClaims{
		SessionData: sess,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		return fmt.Errorf("failed to start jwt session: %v", err)
	}

	c.SetCookie(&http.Cookie{
		Name:     j.CookieName,
		Value:    tokenString,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https",
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func (j *JWTSessionHandler) Destroy(c echo.Context) error {
	// [Critical] Lack of Server-Side Logout
	c.SetCookie(&http.Cookie{
		Name:     j.CookieName,
		Value:    "",
		Expires:  time.UnixMilli(0),
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https",
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func (j *JWTSessionHandler) Refresh(c echo.Context) error {
	return nil
}

func (j *JWTSessionHandler) Rotate(c echo.Context) error {
	return nil
}

func (j *JWTSessionHandler) UpdateSessionData(c echo.Context, updateFunc func(*SessionData) error) error {
	s := SessionDataFromReq(c)
	if s == nil {
		return errors.New("session was nil")
	}

	newS := *s

	updateFunc(&newS)
	return j.Start(c, newS)
}

func (j *JWTSessionHandler) LogoutAllSessionsForUser(id string) error {
	// [Critical] Lack of Server-Side Logout
	return nil
}

func (a *JWTSessionHandler) shouldSkip(c echo.Context) bool {
	path := c.Request().URL.Path
	for _, bypassPath := range a.WhitelistPaths {
		if strings.HasPrefix(path, bypassPath) {
			return true
		}
	}
	return false
}
