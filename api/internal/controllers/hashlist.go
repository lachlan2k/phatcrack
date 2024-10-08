package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
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
	api.POST("/:hashlist-id/append", handleHashlistAppend)
	api.DELETE("/:hashlist-id", handleHashlistDelete)
	api.GET("/:hashlist-id/attacks", handleAttackGetAllForHashlist)
	api.GET("/:hashlist-id/attacks-with-jobs", handleAttacksAndJobsForHashlist)
}

func handleHashlistGetAllForProj(c echo.Context) error {
	projId := c.Param("proj-id")
	if !util.AreValidUUIDs(projId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
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

	user := auth.UserFromReq(c)
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
		return util.ServerError("Failed to check project", err)
	}
	if !allowed {
		return echo.ErrForbidden
	}

	return c.JSON(http.StatusOK, hashlist.ToDTO(true))
}

func handleHashlistDelete(c echo.Context) error {
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(hashlistId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
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
		return util.ServerError("Failed to check project", err)
	}
	if !allowed {
		return echo.ErrForbidden
	}

	jobsToStop, err := db.GetJobsForHashlist(hashlistId)
	if err != nil {
		return util.ServerError("Failed to get jobs to kill", err)
	}

	AuditLog(c, log.Fields{
		"hashlist_id":   hashlist.ID.String(),
		"hashlist_name": hashlist.Name,
		"project_id":    hashlist.ProjectID.String(),
	}, "User is deleting hashlist")

	err = db.HardDelete(hashlist)
	if err != nil {
		return util.ServerError("Failed to delete project", err)
	}

	for _, job := range jobsToStop {
		fleet.StopJob(job, db.JobStopReasonUserStopped)
	}

	return c.JSON(http.StatusOK, "ok")
}

func handleHashlistCreate(c echo.Context) error {
	user := auth.UserFromReq(c)
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

	log.Infof("Validating hashes for new hashlist (%q) for project %q", req.Name, req.ProjectID)

	// Ensure provided algorithm type is valid and normalize
	normalizedHashes, err := hashcathelpers.NormalizeHashes(req.InputHashes, req.HashType, req.HasUsernames)
	if err != nil {
		log.Warnf("Failed to validated hashes for project %q because %v", req.ProjectID, err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate and normalize hashes. Please ensure your hashes are valid for the given hash type.").SetInternal(err)
	}

	hashes := make([]db.HashlistHash, len(normalizedHashes))
	for i, inputHash := range req.InputHashes {
		if req.HasUsernames {
			username, splitHash, err := hashcathelpers.SplitUsername(inputHash)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate hash, error splitting username on line %d", i).SetInternal(err)
			}

			hashes[i].Username = username
			hashes[i].InputHash = splitHash
		} else {
			hashes[i].InputHash = inputHash
		}

		hashes[i].NormalizedHash = normalizedHashes[i]
	}

	newHashlist, err := db.CreateHashlist(&db.Hashlist{
		ProjectID: uuid.MustParse(req.ProjectID),

		Name:    req.Name,
		Version: 1,

		HasUsernames: req.HasUsernames,
		HashType:     req.HashType,
		Hashes:       hashes,
	})

	if err != nil {
		return util.ServerError("Failed to create hashlist", err)
	}

	AuditLog(c, log.Fields{
		"hashlist_name": newHashlist.Name,
		"hashlist_id":   newHashlist.ID.String(),
		"project_id":    req.ProjectID,
	}, "User created a new hashlist")

	numFromPotfile, err := db.PopulateHashlistFromPotfile(newHashlist.ID.String())
	if err != nil {
		log.WithFields(log.Fields{
			"hashlist_name": newHashlist.Name,
			"hashlist_id":   newHashlist.ID.String(),
			"project_id":    req.ProjectID,
		}).WithError(err).Warn("Failed to populated hashlist from potfile")
	}

	return c.JSON(http.StatusCreated, apitypes.HashlistCreateResponseDTO{
		ID:                      newHashlist.ID.String(),
		NumPopulatedFromPotfile: numFromPotfile,
	})
}

func handleHashlistAppend(c echo.Context) error {
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(hashlistId) {
		return echo.ErrBadRequest
	}

	req, err := util.BindAndValidate[apitypes.HashlistAppendRequestDTO](c)
	if err != nil {
		return err
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	allowed, err := accesscontrol.HasRightsToHashlistID(user, hashlistId)
	if err != nil {
		return util.GenericServerError(err)
	}
	if !allowed {
		return echo.ErrForbidden
	}

	hashlist, err := db.GetHashlist(hashlistId)
	if err != nil {
		return err
	}

	// Ensure provided algorithm type is valid and normalize
	normalizedHashes, err := hashcathelpers.NormalizeHashes(req.Hashes, hashlist.HashType, hashlist.HasUsernames)
	if err != nil {
		log.Warnf("Failed to validated newly appended hashes for hashlist %q because %v", hashlistId, err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate and normalize hashes. Please ensure your hashes are valid for the given hash type.").SetInternal(err)
	}

	hashes := make([]db.HashlistHash, len(normalizedHashes))
	for i, inputHash := range req.Hashes {
		if hashlist.HasUsernames {
			username, splitHash, err := hashcathelpers.SplitUsername(inputHash)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate hash, error splitting username on line %d", i).SetInternal(err)
			}

			hashes[i].Username = username
			hashes[i].InputHash = splitHash
		} else {
			hashes[i].InputHash = inputHash
		}

		hashes[i].NormalizedHash = normalizedHashes[i]
		hashes[i].HashlistID = hashlist.ID
	}

	numNewHashes, err := db.AppendToHashlist(hashes)
	if err != nil {
		return util.GenericServerError(err)
	}

	numFromPotfile, err := db.PopulateHashlistFromPotfile(hashlistId)
	if err != nil {
		return util.GenericServerError(err)
	}

	return c.JSON(http.StatusOK, apitypes.HashlistAppendResponseDTO{
		NumNewHashes:            numNewHashes,
		NumPopulatedFromPotfile: numFromPotfile,
	})
}
