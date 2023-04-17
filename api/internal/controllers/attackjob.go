package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

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

func handleAttackJobGetAll(c echo.Context) error {
	projId := c.Param("proj-id")
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(projId, attackId) {
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

func handleAttackJobGet(c echo.Context) error {
	projId := c.Param("proj-id")
	jobId := c.Param("job-id")
	if !util.AreValidUUIDs(projId, jobId) {
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

	job, err := db.GetJob(jobId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}

	if err != nil {
		return util.ServerError("Failed to get job", err)
	}

	return c.JSON(http.StatusOK, job.ToDTO())
}

func handleAttackJobWatch(c echo.Context) error {
	origin := c.Request().Header.Get("origin")
	originU, err := url.Parse(origin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid origin header").SetInternal(err)
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	if len(originU.Host) == 0 || c.Request().Header.Get("host") != originU.Host {
		return echo.NewHTTPError(http.StatusBadRequest, "Cross-origin request are not allowed")
	}

	jobId := c.Param("job-id")
	if !util.AreValidUUIDs(jobId) {
		return echo.ErrBadRequest
	}

	jobProjId, err := db.GetJobProjID(jobId)
	if err != nil {
		if err == db.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Job couldn't be found")
		} else {
			return util.ServerError("Error fetching job", err)
		}
	}

	// Access control
	allowed, err := accesscontrol.HasRightsToProjectID(&user.UserClaims, jobProjId)
	if err != nil {
		return util.ServerError("Error fetching job", err)
	}

	if !allowed {
		return echo.ErrUnauthorized
	}

	ws, err := (&websocket.Upgrader{}).Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return util.ServerError("Couldn't upgrade websocket", err)
	}
	defer ws.Close()
	notifs := fleet.Observe(jobId)
	defer fleet.RemoveObserver(notifs, jobId)

	for {
		notif := <-notifs
		err := ws.WriteJSON(notif)
		if err != nil {
			return err
		}
	}
}
