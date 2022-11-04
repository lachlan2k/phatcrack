package webserver

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/fleet"
	"github.com/lachlan2k/phatcrack/api/pkg/apitypes"
	"golang.org/x/net/websocket"
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
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return c.NoContent(http.StatusUnauthorized)
	}

	authKey := authHeader[len("bearer "):]

	agentData, err := db.FindAgentByAuthKey(authKey)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		agent, err := fleet.RegisterAgentFromWebsocket(ws, agentData.ID.String())
		if err != nil {
			// TODO: c.logger?
			log.Printf("failed to register agent: %v", err)
			return
		}

		err = agent.Handle()
		if err != nil {
			c.Logger().Warnf("Error from agent: %v", err)
		}

	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
