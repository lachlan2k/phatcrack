package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/hashcathelpers"
	"github.com/lachlan2k/phatcrack/api/internal/resources"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookResourceEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api.GET("/hashtypes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, apitypes.HashTypesDTO{
			HashTypes: resources.GetHashTypeMap(),
		})
	})

	api.POST("/detect_hashtype", func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.DetectHashTypeRequestDTO](c)
		if err != nil {
			return err
		}

		possibleTypes, err := hashcathelpers.IdentifyHashTypes(req.TestHash)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find hash candidates").SetInternal(err)
		}

		return c.JSON(http.StatusOK, apitypes.DetectHashTypeResponseDTO{
			PossibleTypes: possibleTypes,
		})
	})
}
