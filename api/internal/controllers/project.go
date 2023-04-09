package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookProjectEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong lists")
	})

	api.GET("", handleProjectGetAll)
	api.POST("/create", handleProjectCreate)
	api.GET("/:id", handleProjectGet)

	api.GET("/:proj-id/hashlist", handleHashlistGetAllForProj)
	api.GET("/:proj-id/hashlist/:hashlist-id", handleHashlistGet)
	api.POST("/:proj-id/hashlist/create", handleHashlistCreate)

	api.GET("/:proj-id/hashlist/:hashlist-id/attack", handleAttackGetAllForHashlist)
	api.GET("/:proj-id/hashlist/:hashlist-id/attack/:attack-id", handleAttackGet)
	api.PUT("/:proj-id/hashlist/:hashlist-id/attack/:attack-id/start", handleAttackStart)
	api.POST("/:proj-id/hashlist/:hashlist-id/attack/create", handleAttackCreate)

	api.GET("/:proj-id/hashlist/:hashlist-id/attack/:attack-id/job", handleAttackJobGetAll)
	api.GET("/:proj-id/hashlist/:hashlist-id/attack/:attack-id/job/:job-id", handleAttackJobGet)
	api.GET("/:proj-id/hashlist/:hashlist-id/attack/:attack-id/job/:job-id", handleAttackJobWatch)

}

func handleProjectCreate(c echo.Context) error {
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	req, err := util.BindAndValidate[apitypes.ProjectCreateDTO](c)
	if err != nil {
		return err
	}

	newProj, err := dbnew.CreateProject(&dbnew.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerUserID: uuid.MustParse(user.ID),
	})
	if err != nil {
		return util.ServerError("Failed to create project", err)
	}

	return c.JSON(http.StatusCreated, newProj.ToDTO())
}

func handleProjectGet(c echo.Context) error {
	projId := c.Param("id")

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := dbnew.GetProjectForUser(projId, user.ID)
	if err == dbnew.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	// Even though our DB query should've constrained it, sanity check with access control regardless
	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		c.Logger().Printf("Something went wrong with getting project %s for user %s, the query returned it, but the user should not have access", proj.ID.String(), user.ID)
		return echo.ErrForbidden
	}

	return c.JSON(http.StatusOK, proj.ToDTO())
}

func handleProjectGetAll(c echo.Context) error {
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projects, err := dbnew.GetAllProjectsForUser(user.ID)
	var res apitypes.ProjectResponseMultipleDTO

	if err == dbnew.ErrNotFound {
		return c.JSON(http.StatusOK, res)
	}
	if err != nil {
		return util.ServerError("Failed to fetch projects", err)
	}

	res.Projects = make([]apitypes.ProjectDTO, 0, len(projects))

	for _, project := range projects {
		// Sanity check access control
		if !accesscontrol.HasRightsToProject(&user.UserClaims, &project) {
			c.Logger().Printf("Something went wrong with getting all projects for user %s, the database query returned project %s, which the user should NOT have access to", user.ID, project.ID.String())
			continue
		}

		res.Projects = append(res.Projects, project.ToDTO())
	}

	return c.JSON(http.StatusOK, res)
}
