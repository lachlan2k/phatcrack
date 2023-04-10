package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func handleAttackGetAllForHashlist(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(projId, hashlistId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := dbnew.GetProjectForUser(projId, user.ID)
	if err == dbnew.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	// res := apitypes.AttackMultipleDTO{}
	return echo.ErrNotImplemented
}

func handleAttackGet(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(projId, hashlistId, attackId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := dbnew.GetProjectForUser(projId, user.ID)
	if err == dbnew.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	return echo.ErrNotImplemented
}

func handleAttackCreate(c echo.Context) error {
	_, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	return echo.ErrNotImplemented
}
