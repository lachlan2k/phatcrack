package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
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

	api.GET("/logout", func(c echo.Context) error {
		sessHandler.Destroy(c)
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
			IsAwaitingMFA:          user.HasRole(auth.RoleMFAEnrolled) && !sessData.HasCompletedMFA,
			RequiresPasswordChange: user.HasRole(auth.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  !user.HasRole(auth.RoleMFAEnrolled) && config.Get().IsMFARequired,
		})
	})

	// Reminder, we're MFA exempt here, so don't put a general password change here
	api.POST("/change-temporary-password", func(c echo.Context) error {
		// TODO: when implementing this, ensure the user has the RoleRequiresPasswordChange role
		// Because this is a setup endpoint, its MFA-exempt
		// For implementing general password changing, we should use a different endpoint
		return echo.ErrNotImplemented
	})

	api.GET("/mfa/is-awaiting-challenge", func(c echo.Context) error {
		user, sessData := auth.UserAndSessFromReq(c)
		if user == nil {
			return echo.ErrUnauthorized
		}

		if sessData.HasCompletedMFA {
			return c.JSON(http.StatusOK, false)
		}

		if user.HasRole(auth.RoleMFAEnrolled) || config.Get().IsMFARequired {
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

		resp, userPresentableError, internalErr := auth.MFAWebAuthnBeginRegister(c, sessHandler)
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA enrollment", internalErr)
		}

		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

		return c.JSONBlob(http.StatusOK, resp)
	})

	api.POST("/mfa/finish-enrollment", func(c echo.Context) error {
		method := c.QueryParam("method")
		// Only supports webauthn right now
		if method != auth.MFATypeWebAuthn {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported MFA type")
		}

		userPresentableError, internalErr := auth.MFAWebAuthnFinishRegister(c, sessHandler)
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA enrollment", internalErr)
		}

		if userPresentableError != nil {
			return echo.NewHTTPError(http.StatusBadRequest, userPresentableError.Error()).SetInternal(userPresentableError)
		}

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

		resp, userPresentableError, internalErr := auth.MFAWebAuthnBeginLogin(c, sessHandler)
		if internalErr != nil {
			return util.ServerError("Something went wrong with MFA registration", internalErr)
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

		userPresentableError, internalErr := auth.MFAWebAuthnFinishLogin(c, sessHandler)
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
		if err != nil {
			return util.ServerError("Something went wrong with MFA completion", err)
		}

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
			IsAwaitingMFA:          user.HasRole(auth.RoleMFAEnrolled) && !sessData.HasCompletedMFA,
			RequiresPasswordChange: user.HasRole(auth.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  !user.HasRole(auth.RoleMFAEnrolled) && config.Get().IsMFARequired,
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
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		if hashingTest != nil {
			return util.ServerError("Internal error", hashingTest)
		}

		sessHandler.Start(c, auth.SessionData{
			UserID:          user.ID.String(),
			HasCompletedMFA: false,
		})

		return c.JSON(http.StatusOK, apitypes.AuthLoginResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    user.Roles,
			},
			IsAwaitingMFA:          user.HasRole(auth.RoleMFAEnrolled),
			RequiresPasswordChange: user.HasRole(auth.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  !user.HasRole(auth.RoleMFAEnrolled) && config.Get().IsMFARequired,
		})
	}
}
