package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"go.mongodb.org/mongo-driver/mongo"
)

func HookProjectEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong lists")
	})

	api.GET("", handleProjectGetAll)
	api.POST("/create", handleProjectCreate)
	api.GET("/:id", handleProjectGet)
	// api.POST("/:id/hashlist", handleProjectGet)
	// api.GET("/:id/hashlist/:list-id", handleProjectGet)
	// api.GET("/:id/hashlist", handleProjectGet)
}

func handleProjectCreate(c echo.Context) error {
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error").SetInternal(err)
	}

	req, err := util.BindAndValidate[apitypes.ProjectCreateDTO](c)
	if err != nil {
		return err
	}

	newProjectId, err := db.CreateProject(db.Project{
		Name:        req.Name,
		Description: req.Description,
	}, user.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create project").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, apitypes.ProjectCreateResponseDTO{
		ID:          newProjectId,
		Name:        req.Name,
		Description: req.Description,
	})
}

func handleProjectGet(c echo.Context) error {
	id := c.Param("id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(id, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrNotFound
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch project").SetInternal(err)
	}

	return c.JSON(http.StatusOK, apitypes.ProjectsFullDetailsDTO{
		ID:          proj.ID.Hex(),
		TimeCreated: proj.ID.Timestamp().UnixMilli(),
		Name:        proj.Name,
		Description: proj.Description,
	})
}

func handleProjectGetAll(c echo.Context) error {
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projects, err := db.GetAllProjectForUser(user.ID)
	var res apitypes.ProjectResponseMultipleDTO

	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusOK, res)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch projects").SetInternal(err)
	}

	res.Projects = make([]apitypes.ProjectSimpleDetailsDTO, len(projects))

	for i, project := range projects {
		res.Projects[i] = apitypes.ProjectSimpleDetailsDTO{
			ID:          project.ID.Hex(),
			TimeCreated: project.ID.Timestamp().UnixMilli(),
			Name:        project.Name,
			Description: project.Description,
		}
	}

	return c.JSON(http.StatusOK, res)
}
