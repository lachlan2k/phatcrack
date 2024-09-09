package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	log "github.com/sirupsen/logrus"
)

func HookPotfileEndpoints(api *echo.Group) {
	api.POST("/search", handlePotfileSearch)
}

func handlePotfileSearch(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.PotfileSearchRequestDTO](c)
	if err != nil {
		return err
	}

	AuditLog(c, log.Fields{
		"num_hashes": len(req.Hashes),
	}, "User is performing a hash search")

	results, err := db.SearchPotfile(req.Hashes)
	if err != nil {
		return util.ServerError("Failed to search for hash", err)
	}

	dtoResults := make([]apitypes.PotfileSearchResultDTO, len(results))
	for i, result := range results {
		dtoResults[i] = result.ToDTO()
	}

	return c.JSON(http.StatusOK, apitypes.PotfileSearchResponseDTO{
		Results: dtoResults,
	})
}
