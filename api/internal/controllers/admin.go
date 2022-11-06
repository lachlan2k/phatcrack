package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookAdminEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.POST("/whoami", func(c echo.Context) error {
		user := c.Get("user")
		return c.JSON(http.StatusOK, user)
	})

	api.POST("/user/create", handleCreateUser)
	api.POST("/agent/create", handleAgentCreate)
}

func handleCreateUser(c echo.Context) error {
	var req apitypes.UserCreateRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	username := db.NormalizeUsername(req.Username)

	userId, err := db.RegisterUser(username, req.Password, req.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, apitypes.UserCreateResponseDTO{
		ID:       userId,
		Username: username,
		Role:     req.Role,
	})
}

func handleAgentCreate(c echo.Context) error {
	var req apitypes.AgentCreateRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	agentId, key, err := db.CreateAgent(req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couldn't create agent").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, apitypes.AgentCreateResponseDTO{
		Name: req.Name,
		ID:   agentId,
		Key:  key,
	})
}
