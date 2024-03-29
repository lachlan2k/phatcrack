package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/filerepo"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func HookAgentHandlerEndpoints(api *echo.Group) {
	// NOTE: this is just for agent handling
	// These endpoints are exempt from useful authz/n
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong agent")
	})

	api.GET("/ws", handleAgentWs)
	api.GET("/download-file/:id", handleAgentDownloadFile)
}

func handleAgentDownloadFile(c echo.Context) error {
	fileId := c.Param("id")
	if !util.AreValidUUIDs(fileId) {
		return echo.ErrBadRequest
	}

	authKey := c.Request().Header.Get("Authorization")
	if len(authKey) == 0 {
		return echo.ErrUnauthorized
	}

	agentId, err := db.FindAgentIDByAuthKey(authKey)

	if err != nil || agentId == "" {
		return echo.ErrUnauthorized
	}

	filename, err := filerepo.GetPathToFile(uuid.MustParse(fileId))
	if err != nil {
		return util.ServerError("Failed to get file", err)
	}

	err = c.File(filename)
	if err == echo.ErrNotFound {
		log.WithField("file_id", fileId).WithField("agent_id", agentId).WithError(err).Warn("Agent tried to download a file that doesn't exist")
		return err
	}
	if err != nil {
		log.WithField("file_id", fileId).WithField("agent_id", agentId).WithError(err).Warn("Agent tried to download a file but encountered an error")
	}
	return err
}

func handleAgentWs(c echo.Context) error {
	authKey := c.Request().Header.Get("Authorization")
	if len(authKey) == 0 {
		return echo.ErrUnauthorized
	}

	agentData, err := db.FindAgentByAuthKey(authKey)
	if err != nil {
		return echo.ErrUnauthorized
	}

	ws, err := (&websocket.Upgrader{}).Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return util.ServerError("Couldn't upgrade websocket", err)
	}

	defer ws.Close()

	agent := fleet.RegisterAgentFromWebsocket(ws, agentData.ID.String())

	AuditLog(c, log.Fields{
		"agent_id":   agentData.ID.String(),
		"agent_name": agentData.Name,
	}, "Agent has connected and is being handled")

	err = agent.Handle()
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"agent_id":   agentData.ID.String(),
			"agent_name": agentData.Name,
		}).Warn("Error from agent")
	}

	return nil
}
