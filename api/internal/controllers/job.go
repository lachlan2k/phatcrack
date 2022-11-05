package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/api/pkg/apitypes"
	"go.mongodb.org/mongo-driver/mongo"
)

func HookJobEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong job")
	})

	api.POST("/create", handleJobCreate)
	api.POST("/:id/start", handleJobStart)
	api.GET("/:id/watch", handleJobWatch)
}

func handleJobWatch(c echo.Context) error {
	origin := c.Request().Header.Get("origin")
	originU, err := url.Parse(origin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid origin header").SetInternal(err)
	}

	if len(originU.Host) == 0 || c.Request().Header.Get("host") != originU.Host {
		return echo.NewHTTPError(http.StatusBadRequest, "Cross-origin request are not allowed")
	}

	job, err := db.GetJob(c.Param("id"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, "Job couldn't be found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching job")
		}
	}

	ws, err := (&websocket.Upgrader{}).Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couldn't upgrade websocket").SetInternal(err)
	}
	defer ws.Close()
	notifs := fleet.Observe(job.ID.String())
	defer fleet.RemoveObserver(notifs, job.ID.String())

	for {
		notif := <-notifs
		err := ws.WriteJSON(notif)
		if err != nil {
			return err
		}
	}
}

func handleJobStart(c echo.Context) error {
	id := c.Param("id")

	agentId, err := fleet.ScheduleJob(id)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, apitypes.JobStartResponseDTO{
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

	return c.JSON(http.StatusCreated, apitypes.JobCreateResponseDTO{
		ID: newJobId,
	})
}
