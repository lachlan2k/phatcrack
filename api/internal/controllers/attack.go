package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"go.mongodb.org/mongo-driver/mongo"
)

func handleAttackGetAllForHashlist(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	res := apitypes.AttackMultipleDTO{}

	// TODO: stop being lazy and write a db query for this
	for _, hashlist := range proj.Hashlists {
		if hashlist.ID.Hex() == hashlistId {
			res.Attacks = make([]apitypes.AttackDTO, len(hashlist.Attacks))
			for i, attack := range hashlist.Attacks {
				res.Attacks[i] = attack.ToDTO()
			}

			return c.JSON(http.StatusOK, res)
		}
	}

	return echo.ErrNotFound
}

func handleAttackGet(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	for _, hashlist := range proj.Hashlists {
		if hashlist.ID.Hex() == hashlistId {
			for _, attack := range hashlist.Attacks {
				if attack.ID.Hex() == attackId {
					return c.JSON(http.StatusOK, attack.ToDTO())
				}
			}
		}
	}

	return echo.ErrForbidden
}

func handleAttackCreate(c echo.Context) error {
	_, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	return echo.ErrNotImplemented
}
