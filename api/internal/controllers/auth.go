package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HookAuthEndpoints(api *echo.Group, authHandler *auth.AuthHandler) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.PUT("/refresh", handleRefresh(authHandler))
	api.POST("/login", handleLogin)

	api.POST("/whoami", func(c echo.Context) error {
		user := c.Get("user")
		return c.JSON(http.StatusOK, user)
	})
}

func handleRefresh(authHandler *auth.AuthHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.ErrUnauthorized
		}
		claims, ok := u.Claims.(*auth.AuthClaims)
		if !ok || claims.ID == "" {
			return echo.ErrUnauthorized
		}

		user, err := db.LookupUserByID(claims.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to refresh user data").SetInternal(err)
		}

		newClaims := auth.UserToClaims(user)

		err = authHandler.SignAndSetJWT(c, newClaims)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

func handleLogin(c echo.Context) error {
	var req apitypes.LoginRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	// TODO: was there a better way of doing this?
	minTime := time.After(time.Second)
	defer func() { <-minTime }()

	username := db.NormalizeUsername(req.Username)

	user, err := db.LookupUserByUsername(username)
	if err == mongo.ErrNoDocuments {
		return echo.ErrUnauthorized
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error").SetInternal(err)
	}

	hashingTest := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if hashingTest != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, apitypes.LoginResponseDTO{
		Username: username,
		Role:     user.Role,
	})
}
