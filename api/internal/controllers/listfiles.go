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

	// TODO: disabling these until I figure out the permissions and workflowuser, err := auth.ClaimsFromReq(c)
	// api.POST("/wordlist/create", handleWordlistCreate)
	// api.POST("/rulefile/create", handleRuleFileCreate)

	api.GET("/wordlist/all", handleGetAllWordlists)
	api.GET("/rulefile/all", handleGetAllRuleFiles)

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
		return util.ServerError("Failed to fetch wordlist", err)
	}

	return c.JSON(http.StatusOK, apitypes.ListsWordlistResponseDTO{
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
		return util.ServerError("Failed to fetch rulefile", err)
	}

	return c.JSON(http.StatusOK, apitypes.ListsRuleFileResponseDTO{
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
		return util.ServerError("Failed to fetch wordlists", err)
	}

	var res apitypes.ListsGetAllWordlistsDTO
	res.Wordlists = make([]apitypes.ListsWordlistResponseDTO, len(lists))
	for i, list := range lists {
		res.Wordlists[i] = apitypes.ListsWordlistResponseDTO{
			ID: list.ID.Hex(),

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
		return util.ServerError("Failed to fetch rulefiles", err)
	}

	var res apitypes.ListsGetAllRuleFilesDTO
	res.RuleFiles = make([]apitypes.ListsRuleFileResponseDTO, len(lists))
	for i, list := range lists {
		res.RuleFiles[i] = apitypes.ListsRuleFileResponseDTO{
			ID:          list.ID.Hex(),
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
	req, err := util.BindAndValidate[apitypes.ListsWordlistCreateDTO](c)
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
		return util.ServerError("Couldn't create wordlist", err)
	}

	return c.NoContent(http.StatusCreated)
}

func handleRuleFileCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.ListsRuleFileCreateDTO](c)
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
		return util.ServerError("Couldn't create rulefile", err)
	}

	return c.NoContent(http.StatusCreated)
}
