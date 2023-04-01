package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/internal/util"
)

func HookAgentEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong agent")
	})

	api.GET("/handle/ws", handleAgentWs)
}

func handleAgentWs(c echo.Context) error {
	authKey := c.Request().Header.Get("X-Agent-Key")
	if len(authKey) == 0 {
		return echo.ErrBadRequest
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

	agent, err := fleet.RegisterAgentFromWebsocket(ws, agentData.ID.Hex())
	if err != nil {
		c.Logger().Printf("Failed to register agent: %v", err)
		return nil
	}

	err = agent.Handle()
	if err != nil {
		c.Logger().Printf("Error from agent: %v", err)
	}

	return nil
}
