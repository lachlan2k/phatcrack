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
		if !isCredentialAuthAllowed() {
			return echo.NewHTTPError(http.StatusBadRequest, "Credential auth is not enabled")
		}

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

		if user.IsPasswordLocked() {
			AuditLog(c, log.Fields{
				"attempted_username": user.Username,
			}, "User attempted to log in with credentials, but password is locked")
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
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
				ID:               user.ID.String(),
				Username:         user.Username,
				Roles:            user.Roles,
				IsPasswordLocked: user.IsPasswordLocked(),
			},
			IsAwaitingMFA:          user.HasRole(roles.UserRoleMFAEnrolled),
			RequiresPasswordChange: user.HasRole(roles.UserRoleRequiresPasswordChange),
			RequiresMFAEnrollment:  config.Get().Auth.General.IsMFARequired && !user.HasRole(roles.UserRoleMFAEnrolled) && !user.HasRole(roles.UserRoleMFAExempt),
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
		}, "%s", logMessage)

		return c.JSON(http.StatusOK, response)
	}
}
