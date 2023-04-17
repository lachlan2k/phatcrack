package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookListsEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong lists")
	})

	api.GET("/wordlist/all", handleGetAllWordlists)
	api.GET("/rulefile/all", handleGetAllRuleFiles)

	api.GET("/wordlist/:id", handleGetWordlist)
	api.GET("/rulefile/:id", handlGetRuleFile)
}

func handleGetWordlist(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	list, err := db.GetWordlist(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch wordlist", err)
	}

	return c.JSON(http.StatusOK, list.ToDTO())
}

func handlGetRuleFile(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	list, err := db.GetRuleFile(id)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Failed to fetch rulefile", err)
	}

	return c.JSON(http.StatusOK, list.ToDTO())
}

func handleGetAllWordlists(c echo.Context) error {
	lists, err := db.GetAllWordlists()
	if err != nil {
		return util.ServerError("Failed to fetch wordlists", err)
	}

	var res apitypes.GetAllWordlistsDTO
	res.Wordlists = make([]apitypes.WordlistDTO, len(lists))
	for i, list := range lists {
		res.Wordlists[i] = list.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}

func handleGetAllRuleFiles(c echo.Context) error {
	lists, err := db.GetAllRuleFiles()
	if err != nil {
		return util.ServerError("Failed to fetch rulefiles", err)
	}

	var res apitypes.GetAllRuleFilesDTO
	res.RuleFiles = make([]apitypes.RuleFileDTO, len(lists))
	for i, list := range lists {
		res.RuleFiles[i] = list.ToDTO()
	}

	return c.JSON(http.StatusOK, res)
}
