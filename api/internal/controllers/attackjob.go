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
)

func HookJobEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong jobs")
	})

	api.GET("/:job-id", handleAttackJobGet)
	api.GET("/:job-id/watch", handleAttackJobWatch)
}

func handleAttackJobGet(c echo.Context) error {
	jobId := c.Param("job-id")
	if !util.AreValidUUIDs(jobId) {
		return echo.ErrBadRequest
	}

	user, err := auth.ClaimsFromReq(c)
	if err != nil {
		return err
	}

	projId, err := db.GetJobProjID(jobId)
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
