package controllers

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/attacksharder"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"gorm.io/datatypes"
)

func HookAttackEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong attacks")
	})

	api.GET("/:attack-id", handleAttackGet)
	api.GET("/all-initialising", handleAttacksGetInitialising)
	api.PUT("/:attack-id/start", handleAttackStart)
	api.POST("/create", handleAttackCreate)

	api.DELETE("/:attack-id/stop", handleAttackStopAllJobs)
	api.DELETE("/:attack-id", handleDeleteAttack)

	api.GET("/:attack-id/jobs", handleAttackJobGetAll)
}

func handleAttacksGetInitialising(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	attacks, err := db.GetAllAttacksWithProgressStringsForUser(user)
	if err != nil {
		return util.ServerError("Failed to get list of initialising attacks", err)
	}

	attackDTOs := make([]apitypes.AttackIDTreeDTO, len(attacks))
	for i, attack := range attacks {
		attackDTOs[i] = attack.ToDTO()
	}

	return c.JSON(http.StatusOK, apitypes.AttackIDTreeMultipleDTO{
		Attacks: attackDTOs,
	})
}

func handleAttackGetAllForHashlist(c echo.Context) error {
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

	proj, err := db.GetProjectForUser(projId, user)
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

	attackDTOs := make([]apitypes.AttackDTO, len(attacks))
	for i, attack := range attacks {
		attackDTOs[i] = attack.ToDTO()
	}

	return c.JSON(http.StatusOK, apitypes.AttackMultipleDTO{
		Attacks: attackDTOs,
	})
}

func handleDeleteAttack(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	attack, err := db.GetAttack(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch attack", err)
	}

	projId, err := db.GetHashlistProjID(attack.HashlistID.String())
	if err != nil {
		return util.ServerError("Failed to fetch project id", err)
	}

	proj, err := db.GetProjectForUser(projId, user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	AuditLog(c, log.Fields{
		"project_name": proj.Name,
		"project_id":   proj.ID.String(),
		"attack_id":    attackId,
	}, "User deleted attack")

	jobsToStop, err := db.GetJobsForAttack(attackId, false)
	if err != nil {
		return util.ServerError("Failed to get jobs to stop", err)
	}

	err = db.HardDelete(attack)
	if err != nil {
		return util.ServerError("Failed to delete attack", err)
	}

	for _, job := range jobsToStop {
		fleet.StopJob(job, db.JobStopReasonUserStopped)
	}

	return c.JSON(http.StatusOK, "ok")
}

func handleAttackStopAllJobs(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	jobs, err := db.GetJobsForAttack(attackId, false)
	if err != nil {
		return util.ServerError("Failed to get jobs for attack", err)
	}

	for _, job := range jobs {
		fleet.StopJob(job, db.JobStopReasonUserStopped)
	}
	return c.JSON(http.StatusOK, "ok")
}

func handleAttackJobGetAll(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for hashlist", err)
	}

	proj, err := db.GetProjectForUser(projId, user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	jobs, err := db.GetJobsForAttack(attackId, c.QueryParams().Has("includeRuntimeData"))
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

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for attack", err)
	}

	proj, err := db.GetProjectForUser(projId, user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	attack, err := db.GetAttack(attackId)
	if err != nil {
		return util.ServerError("Failed to get attack", err)
	}

	return c.JSON(http.StatusOK, attack.ToDTO())
}

func handleAttackCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AttackCreateRequestDTO](c)
	if err != nil {
		return err
	}

	if !util.AreValidUUIDs(req.HashlistID) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	hashlist, err := db.GetHashlist(req.HashlistID)
	if err != nil {
		return util.ServerError("Failed to fetch hahlist for attack", err)
	}

	proj, err := db.GetProjectForUser(hashlist.ProjectID.String(), user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	rulefiles, err := db.GetAllRulefiles()
	if err != nil {
		return util.ServerError("Failed to get information to validate hashcat params", err)
	}
	wordlists, err := db.GetAllWordlists()
	if err != nil {
		return util.ServerError("Failed to get information to validate hashcat params", err)
	}

	// Check all specified wordlists exactly match the ID of a known wordlist
	for _, suppliedWordlist := range req.HashcatParams.WordlistFilenames {
		found := false
		for _, dbWordlist := range wordlists {
			if dbWordlist.ID.String() == suppliedWordlist {
				found = true
				break
			}
		}

		if !found {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid wordlist supplied: %q", suppliedWordlist))
		}
	}

	// Same for rulefiles
	for _, suppliedRulefile := range req.HashcatParams.RulesFilenames {
		found := false
		for _, dbRulefile := range rulefiles {
			if dbRulefile.ID.String() == suppliedRulefile {
				found = true
				break
			}
		}

		if !found {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid rulefile supplied: %v", suppliedRulefile)
		}
	}

	// Don't allow any additional args
	req.HashcatParams.AdditionalArgs = []string{}

	// Enforce correct hashtype
	req.HashcatParams.HashType = uint(hashlist.HashType)

	hashcatParams := datatypes.JSONType[hashcattypes.HashcatParams]{
		Data: req.HashcatParams,
	}

	attack, err := db.CreateAttack(&db.Attack{
		HashcatParams:  hashcatParams,
		IsDistributed:  req.IsDistributed,
		HashlistID:     uuid.MustParse(req.HashlistID),
		ProgressString: "Created",
	})
	if err != nil {
		return util.ServerError("Failed to create new attack", err)
	}

	AuditLog(c, log.Fields{
		"project_id":   proj.ID.String(),
		"project_name": proj.Name,
		"hashlist_id":  req.HashlistID,
	}, "New attack created")

	return c.JSON(http.StatusCreated, attack.ToDTO())
}

func handleAttackStart(c echo.Context) error {
	attackId := c.Param("attack-id")
	if !util.AreValidUUIDs(attackId) {
		return echo.ErrBadRequest
	}

	if config.Get().IsMaintenanceMode {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Phatcrack is in maintenance mode. Attacks cannot be scheduled.")
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for attack", err)
	}

	proj, err := db.GetProjectForUser(projId, user)
	if err == db.ErrNotFound {
		return echo.ErrForbidden
	}
	if err != nil {
		return util.ServerError("Failed to fetch project", err)
	}

	if !accesscontrol.HasRightsToProject(user, proj) {
		return echo.ErrForbidden
	}

	attack, err := db.GetAttack(attackId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Something went wrong getting attack to start", err)
	}

	db.SetAttackProgressString(attackId, "Processing (this can take a while)..")

	jobMultiplier := config.Get().SplitJobsPerAgent
	if jobMultiplier <= 0 {
		jobMultiplier = 1
	}
	numJobs := fleet.NumSchedulableAgents() * jobMultiplier

	errChan := make(chan error, 1)
	successChan := make(chan apitypes.AttackStartResponseDTO, 1)

	AuditLog(c, log.Fields{
		"attack_id":    attack.ID,
		"project_id":   projId,
		"project_name": proj.Name,
		"hashlist_id":  attack.HashlistID,
	}, "User has started attack")

	go func() {
		finalProgressString := "" // blank == success
		defer func() {
			// putting it in a func allows it to capture progress string properly
			db.SetAttackProgressString(attackId, finalProgressString)
		}()

		newJobs, _, err := attacksharder.MakeJobs(attack, numJobs)

		handleErr := func(err error) {
			errId := uuid.NewString()
			finalProgressString = "Error #" + errId
			log.WithFields(log.Fields{
				"attack_id":    attack.ID,
				"project_id":   projId,
				"project_name": proj.Name,
				"hashlist_id":  attack.HashlistID,
				"error_id":     errId,
			}).WithError(err).Warn("Failed to start attack")
			errChan <- err
		}

		if err == db.ErrNotFound {
			errId := uuid.NewString()
			finalProgressString = "Error #" + errId
			log.WithFields(log.Fields{
				"attack_id":    attack.ID,
				"project_id":   projId,
				"project_name": proj.Name,
				"hashlist_id":  attack.HashlistID,
				"error_id":     errId,
			}).WithError(err).Error("Failed to start attack because it was not found")

			errChan <- echo.ErrNotFound
			return
		}

		if err != nil {
			handleErr(util.ServerError("Something went wrong creating attack job", err))
			return
		}

		db.SetAttackProgressString(attackId, "Scheduling jobs...")

		jobIDs := []string{}
		for _, job := range newJobs {
			jobIDs = append(jobIDs, job.ID.String())
		}
		_, err = fleet.ScheduleJobs(jobIDs)

		if err != nil {
			for _, newJob := range newJobs {
				// If the deletion fails, there's not much for us to do really
				db.HardDelete(newJob)
			}
		}

		switch err {
		case nil:
			successChan <- apitypes.AttackStartResponseDTO{
				JobIDs:          jobIDs,
				StillProcessing: false,
			}

		case fleet.ErrJobDoesntExist:
			handleErr(echo.NewHTTPError(http.StatusNotFound, "Job doesn't exist"))

		case fleet.ErrJobAlreadyScheduled:
			handleErr(echo.NewHTTPError(http.StatusBadRequest, "Job already scheduled"))

		case fleet.ErrNoAgentsOnline:
			handleErr(echo.NewHTTPError(http.StatusServiceUnavailable, "No agents are online"))

		default:
			handleErr(util.ServerError("Unexpected error occured when scheduling job", err))
		}

		close(errChan)
		close(successChan)
	}()

	select {
	case e := <-errChan:
		return e

	case res := <-successChan:
		return c.JSON(http.StatusOK, res)

	case <-time.After(3 * time.Second):
		return c.JSON(http.StatusAccepted, apitypes.AttackStartResponseDTO{
			JobIDs:          []string{},
			StillProcessing: true,
		})
	}
}
