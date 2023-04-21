package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookAttackEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong attacks")
	})

	api.GET("/:attack-id", handleAttackGet)
	api.PUT("/:attack-id/start", handleAttackStart)
	api.POST("/create", handleAttackCreate)

	api.GET("/:attack-id/jobs", handleAttackJobGetAll)
}

func handleAttackGetAllForHashlist(c echo.Context) error {
	hashlistId := c.Param("hashlist-id")
	if !util.AreValidUUIDs(hashlistId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projId, err := db.GetHashlistProjID(hashlistId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	// res := apitypes.AttackMultipleDTO{}
	return echo.ErrNotImplemented
}

func handleAttackJobGetAll(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	jobs, err := db.GetJobsForAttack(attackId, projId)
	if err != nil {
		return util.ServerError("Failed to get jobs for attack", err)
	}

	jobDTOs := make([]apitypes.JobSimpleDTO, len(jobs))
	for i, job := range jobs {
		jobDTOs[i] = job.ToSimpleDTO()
	}

	return c.JSON(http.StatusOK, apitypes.JobMultipleDTO{
		Jobs: jobDTOs,
	})
}

func handleAttackGet(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	return echo.ErrNotImplemented
}

func handleAttackCreate(c echo.Context) error {
	_, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	return echo.ErrNotImplemented
}

func handleAttackStart(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(projId, hashlistId, attackId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	hashlist, err := db.GetHashlist(hashlistId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Something went wrong getting hashlist to start attack on", err)
	}

	if hashlist.ProjectID != proj.ID {
		return echo.ErrForbidden
	}

	attack, err := db.GetAttack(attackId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Something went wrong getting attack to start", err)
	}

	if attack.HashlistID != hashlist.ID {
		return echo.ErrForbidden
	}

	targetHashes := make([]string, len(hashlist.Hashes))
	for i, hash := range hashlist.Hashes {
		targetHashes[i] = hash.NormalizedHash
	}

	newJob, err := db.CreateJob(&db.Job{
		HashlistVersion: hashlist.Version,
		AttackID:        &attack.ID,
		HashcatParams:   attack.HashcatParams,
		TargetHashes:    targetHashes,
		HashType:        hashlist.HashType,
	})

	if err != nil {
		return util.ServerError("Couldn't create attack job", err)
	}

	_, err = fleet.ScheduleJob(newJob.ID.String())

	switch err {
	case nil:
		return c.JSON(http.StatusOK, apitypes.AttackStartResponseDTO{
			JobIDs: []string{newJob.ID.String()},
		})

	case fleet.ErrJobDoesntExist:
		return echo.NewHTTPError(http.StatusNotFound, "Job doesn't exist")

	case fleet.ErrJobAlreadyScheduled:
		return echo.NewHTTPError(http.StatusBadRequest, "Job already scheduled")

	case fleet.ErrNoAgentsOnline:
		return echo.NewHTTPError(http.StatusServiceUnavailable, "No agents are online")

	default:
		return util.ServerError("Unexpected error occured when scheduling job", err)
	}
}
