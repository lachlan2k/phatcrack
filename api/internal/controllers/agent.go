package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookAgentEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong agent")
	})

	api.GET("/all", func(c echo.Context) error {
		agents, err := db.GetAllAgents()
		if err != nil {
			return util.ServerError("Couldn't fetch agents", err)
		}

		agentDTOs := make([]apitypes.AgentDTO, len(agents))
		for i, a := range agents {
			agentDTOs[i] = a.ToDTO()
		}

		return c.JSON(http.StatusOK, apitypes.AgentGetAllResponseDTO{
			Agents: agentDTOs,
		})
	})
}
