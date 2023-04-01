package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HookAuthEndpoints(api *echo.Group, authHandler *auth.AuthHandler) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.PUT("/refresh", handleRefresh(authHandler))
	api.POST("/login", handleLogin(authHandler))

	api.GET("/whoami", func(c echo.Context) error {
		claims, err := auth.ClaimsFromReq(c)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, apitypes.AuthWhoamiResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       claims.ID,
				Username: claims.Username,
				Role:     claims.Role,
			},
		})
	})
}

func handleRefresh(authHandler *auth.AuthHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.ClaimsFromReq(c)
		if err != nil {
			return err
		}

		user, err := db.GetUserByID(claims.ID)
		if err != nil {
			return util.ServerError("Failed to refresh user data", err)
		}

		newClaims := auth.UserToClaims(user)

		err = authHandler.SignAndSetJWT(c, newClaims)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, apitypes.AuthWhoamiResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.Hex(),
				Username: user.Username,
				Role:     user.Role,
			},
		})
	}
}

func handleLogin(authHandler *auth.AuthHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.AuthLoginRequestDTO](c)
		if err != nil {
			return err
		}

		// TODO: was there a better way of doing this?
		minTime := time.After(250 * time.Millisecond)
		defer func() { <-minTime }()

		username := db.NormalizeUsername(req.Username)

		user, err := db.GetUserByUsername(username)
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		} else if err != nil {
			return util.ServerError("Database error", err)
		}

		hashingTest := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if hashingTest != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}

		claims := auth.UserToClaims(user)
		authHandler.SignAndSetJWT(c, claims)

		return c.JSON(http.StatusOK, apitypes.AuthLoginResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.Hex(),
				Username: user.Username,
				Role:     user.Role,
			},
		})
	}
}
