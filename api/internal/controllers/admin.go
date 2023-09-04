package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/api/internal/version"
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

	api.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, version.Version())
	})

	api.PUT("/config", func(c echo.Context) error {
		user := auth.UserFromReq(c)
		if user == nil {
			return echo.ErrForbidden
		}

		req, err := util.BindAndValidate[apitypes.AdminConfigRequestDTO](c)
		if err != nil {
			return err
		}

		err = config.Update(func(newConf *config.RuntimeConfig) error {
			newConf.IsMFARequired = req.IsMFARequired
			newConf.AutomaticallySyncListfiles = req.AutomaticallySyncListfiles
			newConf.RequirePasswordChangeOnFirstLogin = req.RequirePasswordChangeOnFirstLogin
			newConf.SplitJobsPerAgent = req.SplitJobsPerAgent
			newConf.MaximumUploadedFileSize = req.MaximumUploadedFileSize
			newConf.MaximumUploadedFileLineScanSize = req.MaximumUploadedFileLineScanSize
			return nil
		})

		if err != nil {
			return util.ServerError("Failed to update config", err)
		}

		AuditLog(c, nil, "Admin updated configuration to %v", req)

		conf := config.Get()
		return c.JSON(http.StatusOK, apitypes.AdminConfigResponseDTO{
			IsSetupComplete:                   conf.IsSetupComplete,
			IsMFARequired:                     conf.IsMFARequired,
			AutomaticallySyncListfiles:        conf.AutomaticallySyncListfiles,
			RequirePasswordChangeOnFirstLogin: conf.RequirePasswordChangeOnFirstLogin,
			SplitJobsPerAgent:                 conf.SplitJobsPerAgent,
			MaximumUploadedFileSize:           conf.MaximumUploadedFileSize,
			MaximumUploadedFileLineScanSize:   conf.MaximumUploadedFileLineScanSize,
		})
	})

	api.GET("/config", func(c echo.Context) error {
		conf := config.Get()
		return c.JSON(http.StatusOK, apitypes.AdminConfigResponseDTO{
			IsSetupComplete:                   conf.IsSetupComplete,
			IsMFARequired:                     conf.IsMFARequired,
			AutomaticallySyncListfiles:        conf.AutomaticallySyncListfiles,
			RequirePasswordChangeOnFirstLogin: conf.RequirePasswordChangeOnFirstLogin,
			SplitJobsPerAgent:                 conf.SplitJobsPerAgent,
			MaximumUploadedFileSize:           conf.MaximumUploadedFileSize,
			MaximumUploadedFileLineScanSize:   conf.MaximumUploadedFileLineScanSize,
		})
	})

	api.GET("/user/all", func(c echo.Context) error {
		users, err := db.GetAllUsers()
		if err != nil {
			return util.ServerError("Failed to get users", err)
		}

		userDTOs := make([]apitypes.UserDTO, len(users))
		for i := range users {
			userDTOs[i] = users[i].ToDTO()
		}

		return c.JSON(http.StatusOK, apitypes.AdminGetAllUsersResponseDTO{
			Users: userDTOs,
		})
	})

	api.POST("/user/create", handleCreateUser)
	api.POST("/user/create-service-account", handleCreateServiceAccount)
	api.POST("/agent/create", handleAgentCreate)

	api.DELETE("/user/:id", handleDeleteUser)
	api.DELETE("/agent/:id", handleDeleteAgent)
}

func handleCreateUser(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminUserCreateRequestDTO](c)
	if err != nil {
		return err
	}

	rolesOk := auth.AreRolesAllowedOnRegistration(req.Roles)
	if !rolesOk {
		return echo.NewHTTPError(http.StatusBadRequest, "One or more provided roles are not allowed on registration")
	}

	if config.Get().RequirePasswordChangeOnFirstLogin {
		req.Roles = append(req.Roles, auth.RoleRequiresPasswordChange)
	}

	// TODO: pro-active handling of duplicate username
	// could also check to see what happens when the constraint fails
	newUser, err := db.RegisterUser(req.Username, req.Password, req.Roles)
	if err != nil {
		return util.ServerError("Couldn't create user", err)
	}

	AuditLog(c, log.Fields{
		"new_user": newUser.ToDTO(),
	}, "New user created")

	return c.JSON(http.StatusCreated, apitypes.AdminUserCreateResponseDTO{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Roles:    newUser.Roles,
	})
}

func handleCreateServiceAccount(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminServiceAccountCreateRequestDTO](c)
	if err != nil {
		return err
	}

	rolesOk := auth.AreRolesAllowedOnRegistration(req.Roles)
	if !rolesOk {
		return echo.NewHTTPError(http.StatusBadRequest, "One or more provided roles are not allowed on registration")
	}

	apiKey, _, err := util.GenAPIKeyAndHash()
	if err != nil {
		return util.ServerError("Couldn't create service account", err)
	}

	allRoles := append(req.Roles, auth.RoleMFAExempt, auth.RoleServiceAccount)

	newUser, err := db.RegisterServiceAccount(req.Username, apiKey, allRoles)
	if err != nil {
		return util.ServerError("Couldn't create service account", err)
	}

	AuditLog(c, log.Fields{
		"new_user": newUser.ToDTO(),
	}, "New service account created")

	return c.JSON(http.StatusCreated, apitypes.AdminServiceAccountCreateResponseDTO{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Roles:    newUser.Roles,
		APIKey:   apiKey,
	})
}

func handleAgentCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminAgentCreateRequestDTO](c)
	if err != nil {
		return err
	}

	newAgent, key, err := db.CreateAgent(req.Name)
	if err != nil {
		return util.ServerError("Failed to create agent", err)
	}

	AuditLog(c, log.Fields{
		"new_agent": newAgent.ToDTO(),
	}, "New agent created")

	return c.JSON(http.StatusCreated, apitypes.AdminAgentCreateResponseDTO{
		Name: req.Name,
		ID:   newAgent.ID.String(),
		Key:  key,
	})
}

func handleDeleteUser(c echo.Context) error {
	id := c.Param("id")

	user, err := db.GetUserByID(id)
	if err == db.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "User does not exist")
	}
	if err != nil {
		return util.ServerError("Failed to retrieve user before deletion", err)
	}

	AuditLog(c, log.Fields{
		"user_to_delete": user.ToDTO(),
	}, "Admin is deleting user")

	err = db.Delete(user)
	if err != nil {
		return util.ServerError("Failed to delete user", err)
	}

	return c.JSON(http.StatusOK, "ok")
}

func handleDeleteAgent(c echo.Context) error {
	id := c.Param("id")

	agent, err := db.GetAgent(id)
	if err == db.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Agent does not exist")
	}
	if err != nil {
		return util.ServerError("Failed to retrieve agent before deletion", err)
	}

	AuditLog(c, log.Fields{
		"agent_to_delete": agent.ToDTO(),
	}, "Admin is deleting agent")

	err = db.Delete(agent)
	if err != nil {
		return util.ServerError("Failed to delete agent", err)
	}

	return c.JSON(http.StatusOK, "ok")
}
