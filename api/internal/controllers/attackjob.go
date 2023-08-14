package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookJobEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong jobs")
	})

	api.GET("/:job-id", handleAttackJobGet)
}

func handleAttackJobGet(c echo.Context) error {
	jobId := c.Param("job-id")
	if !util.AreValidUUIDs(jobId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projId, err := db.GetJobProjID(jobId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user.ID.String())
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	job, err := db.GetJob(jobId, c.QueryParams().Has("includeRuntimeData"))

	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}

	if err != nil {
		return util.ServerError("Failed to get job", err)
	}

	return c.JSON(http.StatusOK, job.ToDTO())
}

func handleAttacksAndJobsForHashlist(c echo.Context) error {
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(hashlistId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	// TODO: This is all a bit gross and could ideally be collapsed into a shorter number of queries?
	projId, err := db.GetHashlistProjID(hashlistId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user.ID.String())
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	attacks, err := db.GetAllAttacksForHashlist(hashlistId)
	if err != nil {
		return util.ServerError("Failed to get attacks for hashlist", err)
	}

	attackDTOs := make([]apitypes.AttackWithJobsDTO, len(attacks))
	for i, attack := range attacks {
		attackDTOs[i].AttackDTO = attack.ToDTO()

		jobs, err := db.GetJobsForAttack(attack.ID.String(), c.QueryParams().Has("includeRuntimeData"))
		if err != nil {
			return util.ServerError("Failed to get job for an attack", err)
		}

		jobDTOs := make([]apitypes.JobDTO, len(jobs))
		for j, job := range jobs {
			jobDTOs[j] = job.ToDTO()
		}

		attackDTOs[i].Jobs = jobDTOs
	}

	return c.JSON(http.StatusOK, apitypes.AttackWithJobsMultipleDTO{
		Attacks: attackDTOs,
	})
}
