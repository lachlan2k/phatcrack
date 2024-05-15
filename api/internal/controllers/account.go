package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"golang.org/x/crypto/bcrypt"
)

func HookAccountEndpoints(api *echo.Group) {
	api.PUT("/change-password", handleChangePassword)
}

func handleChangePassword(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AccountChangePasswordRequestDTO](c)
	if err != nil {
		return err
	}

	u := auth.UserFromReq(c)
	if u == nil {
		return echo.ErrUnauthorized
	}

	AuditLog(c, nil, "User is updating password")

	dbUser, err := db.GetUserByID(u.ID.String())
	if err != nil {
		return util.ServerError("Failed to update password", err)
	}

	if dbUser.IsPasswordLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Password locked")
	}

	currentPasswordCheckErr := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(req.CurrentPassword))
	if currentPasswordCheckErr == bcrypt.ErrMismatchedHashAndPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect current password")
	}
	if currentPasswordCheckErr != nil {
		return util.ServerError("Failed to update password", err)
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return util.ServerError("Failed to update password", err)
	}

	dbUser.PasswordHash = string(newHash)
	err = db.Save(dbUser)
	if err != nil {
		return util.ServerError("Failed to update password", err)
	}

	AuditLog(c, nil, "User successfully updated their password")

	return c.JSON(http.StatusOK, "ok")
}
