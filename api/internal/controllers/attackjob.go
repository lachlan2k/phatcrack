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
	"go.mongodb.org/mongo-driver/mongo"
)

func handleAttackStart(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.CanGetProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	var foundHashlist *db.ProjectHashlist = nil
	var attackToStart *db.HashlistAttack = nil

	for _, hashlist := range proj.Hashlists {
		if hashlist.ID.Hex() == hashlistId {
			foundHashlist = &hashlist
			for _, attack := range hashlist.Attacks {
				if attack.ID.Hex() == attackId {
					attackToStart = &attack
					break
				}
			}
			break
		}
	}

	if attackToStart == nil {
		return echo.ErrNotFound
	}

	normalizedHashes := make([]string, len(attackToStart.Hashes))
	for i, hash := range attackToStart.Hashes {
		normalizedHashes[i] = hash.NormalizedHash
	}

	// Create a job to start the attack
	newJobId, err := db.CreateJob(db.Job{
		ProjectID:       proj.ID,
		HashlistID:      foundHashlist.ID,
		AttackID:        attackToStart.ID,
		HashlistVersion: foundHashlist.Version,

		HashcatParams: db.HashcatParams{
			AttackMode:        attackToStart.HashcatParams.AttackMode,
			HashType:          attackToStart.HashcatParams.HashType,
			Mask:              attackToStart.HashcatParams.Mask,
			WordlistFilenames: attackToStart.HashcatParams.WordlistFilenames,
			RulesFilenames:    attackToStart.HashcatParams.RulesFilenames,
			AdditionalArgs:    attackToStart.HashcatParams.AdditionalArgs,
			OptimizedKernels:  attackToStart.HashcatParams.OptimizedKernels,
		},

		Hashes:   normalizedHashes,
		HashType: attackToStart.HashType,

		RuntimeData: db.RuntimeData{
			Status: db.JobStatusCreated,
		},
	})

	if err != nil {
		return util.ServerError("Couldn't create attack job", err)
	}

	agentId, err := fleet.ScheduleJob(newJobId)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, apitypes.AttackStartResponseDTO{
			NewJobID: newJobId,
			AgentID:  agentId,
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
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.CanGetProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	jobs, err := db.GetJobsForAttack(projId, hashlistId, attackId)
	if err != nil {
		return util.ServerError("Failed to get jobs for attack", err)
	}

	jobDTOs := make([]apitypes.AttackJobSimpleDTO, len(jobs))
	for i, job := range jobs {
		jobDTOs[i] = apitypes.AttackJobSimpleDTO{
			ID:         job.ID.Hex(),
			ProjectID:  job.ProjectID.Hex(),
			HashlistID: job.HashlistID.Hex(),
			AttackID:   job.AttackID.Hex(),
		}
	}

	return c.JSON(http.StatusOK, apitypes.AttackJobMultipleDTO{
		Jobs: jobDTOs,
	})
}

func handleAttackJobGet(c echo.Context) error {
	projId := c.Param("proj-id")
	hashlistId := c.Param("hashlist-id")
	attackId := c.Param("attack-id")
	jobId := c.Param("job-id")
	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	proj, err := db.GetProjectForUser(projId, user.ID)
	if err == mongo.ErrNoDocuments {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.CanGetProject(&user.UserClaims, proj) {
		return echo.ErrForbidden
	}

	job, err := db.GetJobForAttack(projId, hashlistId, attackId, jobId)
	if err == mongo.ErrNoDocuments {
		return echo.ErrNotFound
	}

	if err != nil {
		return util.ServerError("Failed to get job", err)
	}

	return c.JSON(http.StatusOK, apitypes.AttackJobSimpleDTO{
		ID:         job.ID.Hex(),
		ProjectID:  job.ProjectID.Hex(),
		HashlistID: job.HashlistID.Hex(),
		AttackID:   job.AttackID.Hex(),
	})
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
	jobProjId, err := db.GetJobProjID(jobId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
