package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/pkg/apitypes"
)

func hookJobEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong job")
	})

	api.POST("/create", handleJobCreate)
}

func handleJobCreate(c echo.Context) error {
	var req apitypes.JobStartRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	newJobId, err := db.CreateJob(db.Job{
		// TODO from DTO helpers?
		HashcatParams: db.HashcatParams(req.HashcatParams),
		Hashes:        req.Hashes,
		Name:          req.Name,
		Description:   req.Description,
	})

	if err != nil {
		return err
	}

	if req.StartImmediately {
		fleet.ScheduleJob(newJobId)
	}

	return nil
}
