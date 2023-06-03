package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"golang.org/x/crypto/bcrypt"
)

func HookAuthEndpoints(api *echo.Group, authHandler *auth.AuthHandler) {
	// Note: these endpoints are MFA-exempt, so tread carefully before you add anything else
	// If adding a generic endpoint to update password, etc. maybe that should go elsewhere
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
				Roles:    claims.Roles,
			},
		})
	})

	// Reminder, we're MFA exempt here, so don't put a general password change here
	api.POST("/change-temporary-password", func(c echo.Context) error {
		// TODO: when implementing this, ensure the user has the RoleRequiresPasswordChange role
		// Because this is a setup endpoint, its MFA-exempt
		// For implementing general password changing, we should use a different endpoint
		return echo.ErrNotImplemented
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
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
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

		minTime := time.After(250 * time.Millisecond)
		defer func() { <-minTime }()

		user, err := db.GetUserByUsername(req.Username)
		if err == db.ErrNotFound {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		} else if err != nil {
			return util.ServerError("Internal error", err)
		}

		hashingTest := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if hashingTest == bcrypt.ErrMismatchedHashAndPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		if hashingTest != nil {
			return util.ServerError("Internal error", hashingTest)
		}

		claims := auth.UserToClaims(user)
		authHandler.SignAndSetJWT(c, claims)

		return c.JSON(http.StatusOK, apitypes.AuthLoginResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
			},
		})
	}
}
