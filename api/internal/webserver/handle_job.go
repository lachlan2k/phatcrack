package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/api/pkg/apitypes"
)

func hookJobEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong job")
	})

	api.POST("/create", handleJobCreate)
	api.POST("/:id/start", handleJobStart)
}

func handleJobStart(c echo.Context) error {
	id := c.Param("id")

	agentId, err := fleet.ScheduleJob(id)

	switch err {
	case nil:
		return c.JSON(http.StatusCreated, apitypes.JobStartResponseDTO{
			AgentID: agentId,
		})

	case fleet.ErrJobDoesntExist:
		return echo.NewHTTPError(http.StatusNotFound, "Job doesn't exist")

	case fleet.ErrJobAlreadyScheduled:
		return echo.NewHTTPError(http.StatusBadRequest, "Job already scheduled")

	case fleet.ErrNoAgentsOnline:
		return echo.NewHTTPError(http.StatusServiceUnavailable, "No agents are online")

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Unexpected error occured when scheduling job").SetInternal(err)
	}
}

func handleJobCreate(c echo.Context) error {
	var req apitypes.JobCreateRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	newJobId, err := db.CreateJob(db.Job{
		HashcatParams: db.HashcatParams{
			AttackMode:        req.HashcatParams.AttackMode,
			HashType:          req.HashcatParams.HashType,
			Mask:              req.HashcatParams.Mask,
			WordlistFilenames: req.HashcatParams.WordlistFilenames,
			RulesFilenames:    req.HashcatParams.RulesFilenames,
			AdditionalArgs:    req.HashcatParams.AdditionalArgs,
			OptimizedKernels:  req.HashcatParams.OptimizedKernels,
		},
		RuntimeData: db.RuntimeData{
			StartRequestTime: util.MongoNow(),
			Status:           db.JobStatusCreated,
		},
		Hashes:      req.Hashes,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return err
	}

	if req.StartImmediately {
		fleet.ScheduleJob(newJobId)
	}

	return nil
}
