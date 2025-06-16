package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookListsEndpoints(api *echo.Group) {
	api.GET("/all", handleGetAllListfiles)

	api.POST("/upload", handleListfileUpload)

	api.GET("/:id", handleGetListfile)
	api.DELETE("/:id", handleListfileDelete)
}

func handleListfileDelete(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	AuditLog(c, log.Fields{
		"listfile_id": id,
	}, "User is deleting listfile")

	listfile, err := db.GetListfile(id)
	if err != nil {
		return util.ServerError("Failed to get listfile prior to deletion", err)
	}

	isAllowed := user.HasRole(roles.UserRoleAdmin) || listfile.CreatedByUserID == user.ID
	if !isAllowed && listfile.AttachedProjectID != nil {
		projId := listfile.AttachedProjectID.String()
		ok, err := accesscontrol.HasRightsToProjectID(user, projId)
		if err != nil {
			log.WithError(err).WithField("project_id", projId).WithField("user_id", user.ID.String()).Warn("Failed to check project access control for listfile")
		} else if ok {
			isAllowed = true
		}
	}

	if !isAllowed {
		return echo.ErrForbidden
	}

	err = db.MarkListfileForDeletion(id)
	if err != nil {
		return util.ServerError("Failed to mark listfile for deletion", err)
	}

	return c.JSON(http.StatusOK, "ok")
}

func handleGetListfile(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	listfile, err := db.GetListfile(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch wordlist", err)
	}

	if listfile.AttachedProjectID != nil {
		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		ok, err := accesscontrol.HasRightsToProjectID(user, listfile.AttachedProjectID.String())
		if err != nil {
			log.WithError(err).
				WithField("project_id", listfile.AttachedProjectID.String()).
				WithField("user_id", user.ID.String()).
				Warn("Failed to lookup user access to listifle")

			return echo.ErrNotFound
		}

		if !ok {
			return echo.ErrNotFound
		}
	}

	return c.JSON(http.StatusOK, listfile.ToDTO())
}

func handleGetAllListfiles(c echo.Context) error {
	listfiles, err := db.GetAllPublicListfiles()
	if err != nil {
		return util.ServerError("Failed to fetch listfiles", err)
	}

	var res apitypes.GetAllListfilesDTO
	res.Listfiles = make([]apitypes.ListfileDTO, len(listfiles))
	for i, lf := range listfiles {
		res.Listfiles[i] = lf.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}
