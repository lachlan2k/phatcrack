package controllers

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"golang.org/x/crypto/bcrypt"
)

func HookAuthEndpoints(api *echo.Group, sessHandler auth.SessionHandler) {
	// Note: these endpoints are MFA-exempt, so tread carefully before you add anything else
	// If adding a generic endpoint to update password, etc. maybe that should go elsewhere
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.PUT("/refresh", handleRefresh(sessHandler))
	api.POST("/login", handleLogin(sessHandler))

	api.POST("/logout", func(c echo.Context) error {
		sessHandler.Destroy(c)
		AuditLog(c, nil, "User has logged out")

		return c.JSON(http.StatusOK, "Goodbye")
	})

	api.GET("/whoami", func(c echo.Context) error {
		user, sessData := auth.UserAndSessFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		return c.JSON(http.StatusOK, apitypes.AuthWhoamiResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
			},
			IsAwaitingMFA:          user.HasRole(roles.RoleMFAEnrolled) && !sessData.HasCompletedMFA,
			RequiresPasswordChange: user.HasRole(roles.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  config.Get().IsMFARequired && !user.HasRole(roles.RoleMFAEnrolled) && !user.HasRole(roles.RoleMFAExempt),
		})
	})

	// Reminder, we're MFA exempt here, so this is only for users changing their first-set temporary password
	api.POST("/change-temporary-password", func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.AuthChangePasswordRequestDTO](c)
		if err != nil {
			return err
		}

		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		if !user.HasRole(roles.RoleRequiresPasswordChange) {
			return echo.NewHTTPError(http.StatusBadRequest, "This endpoint is only available to users who are required to change their password")
		}

		if len(user.PasswordHash) > 0 {
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Old password was incorrect")
			}
		}

		if req.NewPassword == req.OldPassword {
			return echo.NewHTTPError(http.StatusBadRequest, "New password must be different to old password")
		}

		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return util.ServerError("Failed to update password", err)
		}

		user.PasswordHash = string(newPasswordHash)

		for i, role := range user.Roles {
			if role == roles.RoleRequiresPasswordChange {
				user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
				break
			}
		}

		AuditLog(c, nil, "User has changed their password")

		err = db.Save(user)
		if err != nil {
			return util.ServerError("Failed to update password", err)
		}

		return c.JSON(http.StatusOK, "Ok")
	})

	api.GET("/mfa/is-awaiting-challenge", func(c echo.Context) error {
		user, sessData := auth.UserAndSessFromReq(c)
		if user == nil {
			return echo.ErrUnauthorized
		}

		if sessData.HasCompletedMFA {
			return c.JSON(http.StatusOK, false)
		}

		if user.HasRole(roles.RoleMFAEnrolled) || config.Get().IsMFARequired {
			return c.JSON(http.StatusOK, true)
		}

		return c.JSON(http.StatusOK, false)
	})

	api.POST("/mfa/start-enrollment", func(c echo.Context) error {
		method := c.QueryParam("method")
		// Only supports webauthn right now
		if method != auth.MFATypeWebAuthn {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported MFA type")
		}

		AuditLog(c, log.Fields{
			"mfa_method": method,
		}, "User is attempting to enroll MFA")

		resp, userPresentableError, internalErr := auth.MFAWebAuthnBeginRegister(c, sessHandler)
		if userPresentableError != nil || internalErr != nil {
			AuditLog(c, log.Fields{
				"user_presentable_error": userPresentableError,
				"internal_err":           internalErr,
			}, "User failed to start MFA enrollment")
		}
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA enrollment", internalErr)
		}
		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

		AuditLog(c, nil, "User started MFA enrollment successfully")

		return c.JSONBlob(http.StatusOK, resp)
	})

	api.POST("/mfa/finish-enrollment", func(c echo.Context) error {
		method := c.QueryParam("method")
		// Only supports webauthn right now
		if method != auth.MFATypeWebAuthn {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported MFA type")
		}

		AuditLog(c, nil, "User is attempting to finish MFA enrollment")

		userPresentableError, internalErr := auth.MFAWebAuthnFinishRegister(c, sessHandler)
		if userPresentableError != nil || internalErr != nil {
			AuditLog(c, log.Fields{
				"user_presentable_error": userPresentableError,
				"internal_err":           internalErr,
			}, "User failed to finish MFA enrollment")
			sessHandler.Destroy(c)
		}
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA enrollment", internalErr)
		}
		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

		AuditLog(c, nil, "User successfully finished MFA enrollment")

		return c.JSON(http.StatusCreated, "Ok")
	})

	api.POST("/mfa/start-challenge", func(c echo.Context) error {
		method := c.QueryParam("method")

		// Only supports webauthn right now
		if method != auth.MFATypeWebAuthn {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported MFA type")
		}

		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrUnauthorized
		}

		if user.MFAType != method {
			return echo.NewHTTPError(http.StatusBadRequest, "Incorrect MFA type")
		}

		AuditLog(c, nil, "User has begun MFA challenge")

		resp, userPresentableError, internalErr := auth.MFAWebAuthnBeginLogin(c, sessHandler)
		if userPresentableError != nil || internalErr != nil {
			AuditLog(c, log.Fields{
				"user_presentable_error": userPresentableError,
				"internal_err":           internalErr,
			}, "User failed to begin MFA challenge")
			sessHandler.Destroy(c)
		}
		if internalErr != nil {
			return util.ServerError("Something went wrong with starting MFA challenge", internalErr)
		}
		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

		return c.JSONBlob(http.StatusOK, resp)
	})

	api.POST("/mfa/finish-challenge", func(c echo.Context) error {
		method := c.QueryParam("method")

		// Only supports webauthn right now
		if method != auth.MFATypeWebAuthn {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported MFA type")
		}

		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrUnauthorized
		}

		if user.MFAType != method {
			return echo.NewHTTPError(http.StatusBadRequest, "Incorrect MFA type")
		}

		AuditLog(c, nil, "User is completing MFA Challenge")

		userPresentableError, internalErr := auth.MFAWebAuthnFinishLogin(c, sessHandler)
		if userPresentableError != nil || internalErr != nil {
			AuditLog(c, log.Fields{
				"user_presentable_error": userPresentableError,
				"internal_err":           internalErr,
			}, "User failed to complete MFA challenge")
			sessHandler.Destroy(c)
		}
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA completion", internalErr)
		}
		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

		err := sessHandler.UpdateSessionData(c, func(sd *auth.SessionData) error {
			sd.HasCompletedMFA = true
			return nil
		})
		sessHandler.Rotate(c)
		if err != nil {
			return util.ServerError("Something went wrong with MFA completion", err)
		}

		AuditLog(c, nil, "User has authenticated")
		return c.JSON(http.StatusOK, "Ok")
	})
}

func handleRefresh(sessHandler auth.SessionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := sessHandler.Refresh(c)
		if err != nil {
			return err
		}

		user, sessData := auth.UserAndSessFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		return c.JSON(http.StatusOK, apitypes.AuthWhoamiResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
			},
			IsAwaitingMFA:          user.HasRole(roles.RoleMFAEnrolled) && !sessData.HasCompletedMFA,
			RequiresPasswordChange: user.HasRole(roles.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  config.Get().IsMFARequired && !user.HasRole(roles.RoleMFAEnrolled) && !user.HasRole(roles.RoleMFAExempt),
		})
	}
}

func handleLogin(sessHandler auth.SessionHandler) echo.HandlerFunc {
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
			RequiresMFAEnrollment:  config.Get().IsMFARequired && !user.HasRole(roles.RoleMFAEnrolled) && !user.HasRole(roles.RoleMFAExempt),
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
		}, logMessage)

		return c.JSON(http.StatusOK, response)
	}
}
