package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookAdminEndpoints(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong auth")
	})

	api.GET("/whoami", func(c echo.Context) error {
		user := c.Get("user")
		return c.JSON(http.StatusOK, user)
	})

	api.POST("/user/create", handleCreateUser)
	api.POST("/agent/create", handleAgentCreate)
}

func handleCreateUser(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminUserCreateRequestDTO](c)
	if err != nil {
		return err
	}

	username := db.NormalizeUsername(req.Username)

	userId, err := db.RegisterUser(username, req.Password, req.Role)
	if err != nil {
		return util.ServerError("Failed to create user", err)
	}

	return c.JSON(http.StatusCreated, apitypes.AdminUserCreateResponseDTO{
		ID:       userId,
		Username: username,
		Role:     req.Role,
	})
}

func handleAgentCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminAgentCreateRequestDTO](c)
	if err != nil {
		return err
	}

	agentId, key, err := db.CreateAgent(req.Name)
	if err != nil {
		return util.ServerError("Failed to create agent", err)
	}

	return c.JSON(http.StatusCreated, apitypes.AdminAgentCreateResponseDTO{
		Name: req.Name,
		ID:   agentId,
		Key:  key,
	})
}
