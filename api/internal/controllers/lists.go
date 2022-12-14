package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"go.mongodb.org/mongo-driver/mongo"
)

func HookListsEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong lists")
	})

	api.POST("/wordlist/create", handleWordlistCreate)
	api.POST("/rulefile/create", handleRuleFileCreate)

	api.GET("/wordlists", handleGetAllWordlists)
	api.GET("/rulefiles", handleGetAllRuleFiles)

	api.GET("/wordlist/:id", handleGetWordlist)
	api.GET("/rulefile/:id", handlGetRuleFile)
}

func handleGetWordlist(c echo.Context) error {
	id := c.Param("id")
	list, err := db.GetWordlist(id)
	if err == mongo.ErrNoDocuments {
		return echo.ErrNotFound
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch wordlist").SetInternal(err)
	}

	return c.JSON(http.StatusOK, apitypes.WordlistResponseDTO{
		Name:        list.Name,
		Description: list.Description,
		Filename:    list.Filename,
		Size:        list.Size,
		Lines:       list.Lines,
	})
}

func handlGetRuleFile(c echo.Context) error {
	id := c.Param("id")
	list, err := db.GetRuleFile(id)
	if err == mongo.ErrNoDocuments {
		return echo.ErrNotFound
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch wordlist").SetInternal(err)
	}

	return c.JSON(http.StatusOK, apitypes.RuleFileResponseDTO{
		Name:        list.Name,
		Description: list.Description,
		Filename:    list.Filename,
		Size:        list.Size,
		Lines:       list.Lines,
	})
}

func handleGetAllWordlists(c echo.Context) error {
	lists, err := db.GetAllWordlists()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load wordlists").SetInternal(err)
	}

	var res apitypes.GetAllWordlistsDTO
	res.Wordlists = make([]apitypes.WordlistResponseDTO, len(lists))
	for i, list := range lists {
		res.Wordlists[i] = apitypes.WordlistResponseDTO{
			Name:        list.Name,
			Description: list.Description,
			Filename:    list.Filename,
			Size:        list.Size,
			Lines:       list.Lines,
		}
	}

	return c.JSON(http.StatusOK, res)
}

func handleGetAllRuleFiles(c echo.Context) error {
	lists, err := db.GetAllRuleFiles()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load wordlists").SetInternal(err)
	}

	var res apitypes.GetAllRuleFilesDTO
	res.RuleFiles = make([]apitypes.RuleFileResponseDTO, len(lists))
	for i, list := range lists {
		res.RuleFiles[i] = apitypes.RuleFileResponseDTO{
			Name:        list.Name,
			Description: list.Description,
			Filename:    list.Filename,
			Size:        list.Size,
			Lines:       list.Lines,
		}
	}

	return c.JSON(http.StatusOK, res)
}

func handleWordlistCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.WordlistCreateDTO](c)
	if err != nil {
		return err
	}

	err = db.AddWordlist(db.Wordlist{
		Name:        req.Name,
		Description: req.Description,
		Filename:    req.Filename,
		Size:        req.Size,
		Lines:       req.Lines,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create wordlist").SetInternal(err)
	}

	return c.NoContent(http.StatusCreated)
}

func handleRuleFileCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.RuleFileCreateDTO](c)
	if err != nil {
		return err
	}

	err = db.AddRuleFile(db.RuleFile{
		Name:        req.Name,
		Description: req.Description,
		Filename:    req.Filename,
		Size:        req.Size,
		Lines:       req.Lines,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add rulefile").SetInternal(err)
	}

	return c.NoContent(http.StatusCreated)
}
