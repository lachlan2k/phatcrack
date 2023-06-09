package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/hashcathelpers"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookHashlistEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong hashlist")
	})

	api.POST("/create", handleHashlistCreate)
	api.GET("/:hashlist-id", handleHashlistGet)
	api.GET("/:hashlist-id/attacks", handleAttackGetAllForHashlist)
	api.GET("/:hashlist-id/attacks-with-jobs", handleAttacksAndJobsForHashlist)
}

func handleHashlistGetAllForProj(c echo.Context) error {
	projId := c.Param("proj-id")
	if !util.AreValidUUIDs(projId) {
		return echo.ErrBadRequest
	}

	user, _ := auth.UserAndSessFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	allowed, err := accesscontrol.HasRightsToProjectID(user, projId)
	if err != nil {
		return err
	}
	if !allowed {
		return echo.ErrForbidden
	}

	hashlists, err := db.GetAllHashlistsForProject(projId)
	if err != nil {
		return util.ServerError("Failed to get hashlists", err)
	}

	res := apitypes.HashlistResponseMultipleDTO{}
	res.Hashlists = make([]apitypes.HashlistDTO, len(hashlists))

	for i, hashlist := range hashlists {
		res.Hashlists[i] = hashlist.ToDTO(false)
	}

	return c.JSON(http.StatusOK, res)
}

func handleHashlistGet(c echo.Context) error {
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(hashlistId) {
		return echo.ErrBadRequest
	}

	user, _ := auth.UserAndSessFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	hashlist, err := db.GetHashlistWithHashes(hashlistId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to load hashlist", err)
	}

	allowed, err := accesscontrol.HasRightsToProjectID(user, hashlist.ProjectID.String())
	if err != nil {
		return err
	}
	if !allowed {
		return echo.ErrForbidden
	}

	return c.JSON(http.StatusOK, hashlist.ToDTO(true))
}

func handleHashlistCreate(c echo.Context) error {
	user, _ := auth.UserAndSessFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	req, err := util.BindAndValidate[apitypes.HashlistCreateRequestDTO](c)
	if err != nil {
		return err
	}

	// Access control
	allowed, err := accesscontrol.HasRightsToProjectID(user, req.ProjectID)
	if err != nil {
		return err
	}
	if !allowed {
		return echo.ErrForbidden
	}

	// Ensure provided algorithm type is valid and normalize
	normalizedHashes, err := hashcathelpers.NormalizeHashes(req.InputHashes, req.HashType, req.HasUsernames)
	if err != nil {
		c.Logger().Printf("Failed to validated hashes for project %s because %v", req.ProjectID, err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate and normalize hashes. Please ensure your hashes are valid for the given hash type.").SetInternal(err)
	}

	hashes := make([]db.HashlistHash, len(normalizedHashes))
	for i, inputHash := range req.InputHashes {
		hashes[i].InputHash = inputHash
		hashes[i].NormalizedHash = normalizedHashes[i]
	}

	newHashlist, err := db.CreateHashlist(&db.Hashlist{
		ProjectID: uuid.MustParse(req.ProjectID),

		Name:    req.Name,
		Version: 1,

		HashType: req.HashType,
		Hashes:   hashes,
	})

	if err != nil {
		return util.ServerError("Failed to create hashlist", err)
	}

	return c.JSON(http.StatusCreated, apitypes.HashlistCreateResponseDTO{
		ID: newHashlist.ID.String(),
	})
}
