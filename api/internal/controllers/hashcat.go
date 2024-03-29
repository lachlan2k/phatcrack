package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/hashcathelpers"
	"github.com/lachlan2k/phatcrack/api/internal/resources"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookHashcatEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api.GET("/hashtypes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, apitypes.HashTypesDTO{
			HashTypes: resources.GetHashTypeMap(),
		})
	})

	api.POST("/detect-hashtype", func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.DetectHashTypeRequestDTO](c)
		if err != nil {
			return err
		}

		possibleTypes, err := hashcathelpers.IdentifyHashTypes(req.TestHash, req.HasUsername)
		if err != nil {
			return util.ServerError("Failed to find hash candidates", err)
		}

		return c.JSON(http.StatusOK, apitypes.DetectHashTypeResponseDTO{
			PossibleTypes: possibleTypes,
		})
	})

	api.POST("/verify-hashes", func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.VerifyHashesRequestDTO](c)
		if err != nil {
			return err
		}

		// TODO: Do more investigation as to what happens when hashcat rejects hashes
		normalized, err := hashcathelpers.NormalizeHashes(req.Hashes, req.HashType, req.HasUsernames)

		isValid := len(normalized) > 0 && err != nil
		return c.JSON(http.StatusOK, apitypes.VerifyHashesResponseDTO{
			Valid: isValid,
		})
	})

	api.POST("/normalize-hashes", func(c echo.Context) error {
		req, err := util.BindAndValidate[apitypes.NormalizeHashesRequestDTO](c)
		if err != nil {
			return err
		}

		// TODO: Do more investigation as to what happens when hashcat rejects hashes
		normalized, err := hashcathelpers.NormalizeHashes(req.Hashes, req.HashType, req.HasUsernames)

		isValid := len(normalized) > 0 && err != nil
		if !isValid {
			normalized = []string{}
		}

		return c.JSON(http.StatusOK, apitypes.NormalizeHashesResponseDTO{
			Valid:            isValid,
			NormalizedHashes: normalized,
		})
	})
}
