package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func handleCredentialLogin(sessHandler auth.SessionHandler) echo.HandlerFunc {
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
			AuditLog(c, log.Fields{
				"attempted_username": user.Username,
			}, "Incorrect password submitted")
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		if hashingTest != nil {
			return util.ServerError("Internal error", hashingTest)
		}

		sessHandler.Start(c, auth.SessionData{
			UserID:          user.ID.String(),
			HasCompletedMFA: false,
		})

		response := apitypes.AuthLoginResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
			},
			IsAwaitingMFA:          user.HasRole(roles.RoleMFAEnrolled),
			RequiresPasswordChange: user.HasRole(roles.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  config.Get().Auth.IsMFARequired && !user.HasRole(roles.RoleMFAEnrolled) && !user.HasRole(roles.RoleMFAExempt),
		}

		logMessage := "Session started"
		if response.IsAwaitingMFA {
			logMessage += ", user needs to complete MFA"
		}
		if response.RequiresPasswordChange {
			logMessage += ", user needs to change password"
		}
		if response.RequiresMFAEnrollment {
			logMessage += ", user needs to enroll MFA"
		}

		AuditLog(c, log.Fields{
			"authenticated_username": user.Username,
			"auth_type":              "credentials",
		}, logMessage)

		return c.JSON(http.StatusOK, response)
	}
}
