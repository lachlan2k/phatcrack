package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/accesscontrol"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
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
	api.PUT("/:attack-id/start", handleAttackStart)
	api.POST("/create", handleAttackCreate)

	api.GET("/:attack-id/jobs", handleAttackJobGetAll)
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

	attackDTOs := make([]apitypes.AttackDTO, len(attacks))
	for i, attack := range attacks {
		attackDTOs[i] = attack.ToDTO()
	}

	return c.JSON(http.StatusOK, apitypes.AttackMultipleDTO{
		Attacks: attackDTOs,
	})
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

	// TODO: this could be collapsed into less queries too
	projId, err := db.GetAttackProjID(attackId)
	if err != nil {
		return util.ServerError("Failed to fetch project id for attack", err)
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

	projId, err := db.GetHashlistProjID(req.HashlistID)
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
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid wordlist supplied: %v", suppliedWordlist)
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

	hashcatParams := datatypes.JSONType[hashcattypes.HashcatParams]{
		Data: req.HashcatParams,
	}

	attack, err := db.CreateAttack(&db.Attack{
		HashcatParams: hashcatParams,
		HashlistID:    uuid.MustParse(req.HashlistID),
	})
	if err != nil {
		return util.ServerError("Failed to create new attack", err)
	}

	return c.JSON(http.StatusCreated, attack.ToDTO())
}

func handleAttackStart(c echo.Context) error {
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

	attack, err := db.GetAttack(attackId)
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Something went wrong getting attack to start", err)
	}

	hashlist, err := db.GetHashlistWithHashes(attack.HashlistID.String())
	if err == db.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.ServerError("Something went wrong getting hashlist to start attack on", err)
	}

	targetHashes := make([]string, len(hashlist.Hashes))
	for _, hash := range hashlist.Hashes {
		if !hash.IsCracked {
			targetHashes = append(targetHashes, hash.NormalizedHash)
		}
	}

	// TODO: this will eventually be replaced by some special sauce responsible for sharding and creating jobs
	// Could possibly go into the fleet package, or another package called "sharding", idk

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

	// TODO: now that we've refactored the relationship between attacks and jobs...
	// ...we should probably just delete the job if it fails, right?
	// TODO: ...done?
	if err != nil {
		// If the deletion fails, there's not much for us to do really
		db.GetInstance().Delete(newJob)
	}

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
