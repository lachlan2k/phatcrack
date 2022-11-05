package webserver

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/pkg/apitypes"
)

func hookAgentEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong agent")
	})

	api.POST("/create", handleAgentCreate)
	api.GET("/handle/ws", handleAgentWs)
}

func handleAgentCreate(c echo.Context) error {
	var req apitypes.AgentCreateRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	agentId, key, err := db.CreateAgent(req.Name)
	if err != nil {
		// TODO: err
		return c.String(http.StatusInternalServerError, "couldn't make agent")
	}

	return c.JSON(http.StatusCreated, apitypes.AgentCreateResponseDTO{
		Name: req.Name,
		ID:   agentId,
		Key:  key,
	})
}

func handleAgentWs(c echo.Context) error {
	authKey := c.Request().Header.Get("X-Agent-Key")
	if len(authKey) == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}

	agentData, err := db.FindAgentByAuthKey(authKey)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	ws, err := (&websocket.Upgrader{}).Upgrade(c.Response(), c.Request(), nil)

	defer ws.Close()

	agent, err := fleet.RegisterAgentFromWebsocket(ws, agentData.ID.String())
	if err != nil {
		c.Logger().Warnf("Failed to register agent: %v", err)
		return nil
	}

	err = agent.Handle()
	if err != nil {
		c.Logger().Warnf("Error from agent: %v", err)
	}

	return nil
}
