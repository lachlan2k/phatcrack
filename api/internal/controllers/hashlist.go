package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/hashcathelpers"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func handleHashlistGetAllForProj(c echo.Context) error {
	return echo.ErrNotImplemented
}

func handleHashlistGet(c echo.Context) error {
	return echo.ErrNotImplemented
}

func handleHashlistCreate(c echo.Context) error {
	projId := c.Param("proj-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	req, err := util.BindAndValidate[apitypes.HashlistCreateRequestDTO](c)
	if err != nil {
		return err
	}

	// Access control
	allowed, err := accesscontrol.HasRightsToProjectID(&user.UserClaims, projId)
	if err != nil {
		return err
	}
	if !allowed {
		return echo.ErrForbidden
	}

	// Ensure provided algorithm type is valid and normalize
	normalizedHashes, err := hashcathelpers.NormalizeHashes(req.InputHashes, req.HashType, req.HasUsernames)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate and normalize hashes. Please ensure your hashes are valid for the given hash type.").SetInternal(err)
	}

	hashes := make([]db.HashlistHash, 0, len(normalizedHashes))
	for i, inputHash := range req.InputHashes {
		hashes[i].InputHash = inputHash
		hashes[i].NormalizedHash = normalizedHashes[i]
	}

	newHashlistId, err := db.AddHashlistToProject(projId, db.ProjectHashlist{
		Name:     req.Name,
		HashType: req.HashType,
		Hashes:   hashes,
		Version:  1,
	})

	if err != nil {
		return util.ServerError("Failed to create hashlist", err)
	}

	return c.JSON(http.StatusCreated, apitypes.HashlistCreateResponseDTO{
		ID: newHashlistId,
	})
}
