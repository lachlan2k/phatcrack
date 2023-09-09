package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookProjectEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong project")
	})

	api.GET("/all", handleProjectGetAll)
	api.POST("/create", handleProjectCreate)
	api.GET("/:id", handleProjectGet)
	api.DELETE("/:id", handleProjectDelete)

	api.GET("/:proj-id/hashlists", handleHashlistGetAllForProj)
}

func handleProjectCreate(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	req, err := util.BindAndValidate[apitypes.ProjectCreateRequestDTO](c)
	if err != nil {
		return err
	}

	newProj, err := db.CreateProject(&db.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerUserID: user.ID,
	})
	if err != nil {
		return util.ServerError("Failed to create project", err)
	}

	AuditLog(c, log.Fields{
		"project_id":   newProj.ID,
		"project_name": newProj.Name,
	}, "User created a new project")

	return c.JSON(http.StatusCreated, newProj.ToDTO())
}

func handleProjectGet(c echo.Context) error {
	projId := c.Param("id")
	if !util.AreValidUUIDs(projId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	proj, err := db.GetProjectForUser(projId, user.ID.String())
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	// Even though our DB query should've constrained it, sanity check with access control regardless
	if !accesscontrol.HasRightsToProject(user, proj) {
		log.WithField("user_id", user.ID).WithField("project_id", proj.ID.String()).Error("Access control violation: db query returned a project the user should not have access to")
		return echo.ErrForbidden
	}

	return c.JSON(http.StatusOK, proj.ToDTO())
}

func handleProjectDelete(c echo.Context) error {
	projId := c.Param("id")
	if !util.AreValidUUIDs(projId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	proj, err := db.GetProjectForUser(projId, user.ID.String())
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	// Only the owner of a project, or an admin, can delete a project
	if !user.HasRole(auth.RoleAdmin) && proj.OwnerUserID != user.ID {
		AuditLog(c, log.Fields{
			"project_name": proj.Name,
			"project_id":   proj.ID.String(),
		}, "User tried to delete project, but they are not allowed")
		return echo.ErrForbidden
	}

	err = db.HardDelete(proj)
	if err != nil {
		return util.ServerError("Failed to delete project", err)
	}

	AuditLog(c, log.Fields{
		"project_name": proj.Name,
		"project_id":   proj.ID.String(),
	}, "User deleted project")

	return c.JSON(http.StatusOK, "ok")
}

func handleProjectGetAll(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projects, err := db.GetAllProjectsForUser(user.ID.String())
	var res apitypes.ProjectResponseMultipleDTO

	if err == db.ErrNotFound {
		return c.JSON(http.StatusOK, res)
	}
	if err != nil {
		return util.ServerError("Failed to fetch projects", err)
	}

	res.Projects = make([]apitypes.ProjectDTO, 0, len(projects))

	for _, project := range projects {
		// Sanity check access control
		if !accesscontrol.HasRightsToProject(user, &project) {
			log.WithField("user_id", user.ID).WithField("project_id", project.ID.String()).Error("Access control violation: db query returned a project the user should not have access to")
			continue
		}

		res.Projects = append(res.Projects, project.ToDTO())
	}

	return c.JSON(http.StatusOK, res)
}
