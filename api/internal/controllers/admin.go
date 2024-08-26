package controllers

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"crypto/rand"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
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
			if req.Agent != nil {
				a := *req.Agent

				newConf.Agent.AutomaticallySyncListfiles = a.AutomaticallySyncListfiles
				newConf.Agent.SplitJobsPerAgent = a.SplitJobsPerAgent
			}

			if req.Auth != nil {
				a := *req.Auth

				if a.General != nil {
					newConf.Auth.General.EnabledMethods = a.General.EnabledMethods
					newConf.Auth.General.IsMFARequired = a.General.IsMFARequired
					newConf.Auth.General.RequirePasswordChangeOnFirstLogin = a.General.RequirePasswordChangeOnFirstLogin
				}

				if a.OIDC != nil {
					newConf.Auth.OIDC.ClientID = a.OIDC.ClientID
					if a.OIDC.ClientSecret != "redacted" {
						newConf.Auth.OIDC.ClientSecret = a.OIDC.ClientSecret
					}
					newConf.Auth.OIDC.IssuerURL = a.OIDC.IssuerURL
					newConf.Auth.OIDC.RedirectURL = a.OIDC.RedirectURL
					newConf.Auth.OIDC.AutomaticUserCreation = a.OIDC.AutomaticUserCreation
					newConf.Auth.OIDC.UsernameClaim = a.OIDC.UsernameClaim
					newConf.Auth.OIDC.Prompt = a.OIDC.Prompt
					newConf.Auth.OIDC.RolesClaim = a.OIDC.RolesClaim
					newConf.Auth.OIDC.RequiredRole = a.OIDC.RequiredRole
					newConf.Auth.OIDC.AdditionalScopes = a.OIDC.AdditionalScopes
				}
			}

			if req.General != nil {
				g := *req.General

				newConf.General.IsMaintenanceMode = g.IsMaintenanceMode
				newConf.General.MaximumUploadedFileSize = g.MaximumUploadedFileSize
				newConf.General.MaximumUploadedFileLineScanSize = g.MaximumUploadedFileLineScanSize
			}

			return nil
		})

		if err != nil {
			return util.ServerError("Failed to update config", err)
		}

		AuditLog(c, nil, "Admin updated configurationv")

		return c.JSON(http.StatusOK, config.Get().ToAdminDTO())
	})

	api.GET("/config", func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.Get().ToAdminDTO())
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
	api.PUT("/agent/:id/set-maintenance-mode", func(c echo.Context) error {
		id := c.Param("id")
		req, err := util.BindAndValidate[apitypes.AdminAgentSetMaintanceRequestDTO](c)
		if err != nil {
			return err
		}

		err = db.UpdateAgentMaintenanceMode(id, req.IsMaintenanceMode)
		if err != nil {
			return util.ServerError("Failed to set agent's maintenance mode", err)
		}

		return c.JSON(http.StatusOK, "ok")
	})

	api.POST("/agent-registration-key/create", handleAgentRegistrationKeyCreate)

	api.PUT("/user/:id", handleUpdateUser)

	api.DELETE("/user/:id", handleDeleteUser)
	api.DELETE("/agent/:id", handleDeleteAgent)
}

func handleCreateUser(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminUserCreateRequestDTO](c)
	if err != nil {
		return err
	}

	password := req.Password
	generatedPassword := ""

	if req.GenPassword {
		genBuff := make([]byte, 16)
		_, err := rand.Read(genBuff)
		if err != nil {
			return util.ServerError("Couldn't create user", err)
		}

		generatedPassword = hex.EncodeToString(genBuff)
		password = generatedPassword
	} else {
		if pwordOk, pwordFb := util.ValidatePasswordStrength(password); !pwordOk {
			return echo.NewHTTPError(http.StatusBadRequest, "Password did not meet strength requirements: "+pwordFb)
		}
	}

	rolesOk := roles.AreRolesAssignable(req.Roles)
	if !rolesOk {
		return echo.NewHTTPError(http.StatusBadRequest, "One or more provided roles are not allowed on registration")
	}

	if config.Get().Auth.General.RequirePasswordChangeOnFirstLogin {
		req.Roles = append(req.Roles, roles.UserRoleRequiresPasswordChange)
	}

	newUser, err := db.RegisterUserWithCredentials(req.Username, password, req.Roles)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusConflict, "A user with that username already exists")
		}

		return util.ServerError("Couldn't create user", err)
	}

	AuditLog(c, log.Fields{
		"new_user": newUser.ToDTO(),
	}, "New user created")

	return c.JSON(http.StatusCreated, apitypes.AdminUserCreateResponseDTO{
		ID:                newUser.ID.String(),
		Username:          newUser.Username,
		Roles:             newUser.Roles,
		GeneratedPassword: generatedPassword,
	})
}

func handleCreateServiceAccount(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminServiceAccountCreateRequestDTO](c)
	if err != nil {
		return err
	}

	rolesOk := roles.AreRolesAssignable(req.Roles)
	if !rolesOk {
		return echo.NewHTTPError(http.StatusBadRequest, "One or more provided roles are not allowed on registration")
	}

	apiKey, _, err := util.GenAPIKeyAndHash()
	if err != nil {
		return util.ServerError("Couldn't create service account", err)
	}

	allRoles := append(req.Roles, roles.UserRoleMFAExempt, roles.UserRoleServiceAccount)

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

	newAgent, key, err := db.CreateAgent(req.Name, req.Ephemeral)
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

func handleAgentRegistrationKeyCreate(c echo.Context) error {
	req, err := util.BindAndValidate[apitypes.AdminAgentRegistrationKeyCreateRequestDTO](c)
	if err != nil {
		return err
	}

	newRegKey, key, err := db.CreateAgentRegistrationKey(req.Name, req.Ephemeral)
	if err != nil {
		return util.ServerError("Failed to create agent registration key", err)
	}

	AuditLog(c, log.Fields{
		"key_id":   newRegKey.ID,
		"key_name": newRegKey.Name,
	}, "New agent registration key created")

	return c.JSON(http.StatusCreated, apitypes.AdminAgentRegistrationKeyCreateResponseDTO{
		Ephemeral: newRegKey.Ephemeral,
		Name:      newRegKey.Name,
		ID:        strconv.Itoa(int(newRegKey.ID)),
		Key:       key,
	})
}

func handleUpdateUser(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	req, err := util.BindAndValidate[apitypes.AdminUserUpdateRequestDTO](c)
	if err != nil {
		return err
	}

	user, err := db.GetUserByID(id)
	if err == db.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "User does not exist")
	}
	if err != nil {
		return util.ServerError("Failed to retrieve user to update", err)
	}

	AuditLog(c, log.Fields{
		"user_id":      id,
		"new_username": req.Username,
		"new_roles":    req.Roles,
	}, "Admin is updating user")

	if req.Username != user.Username {
		_, err := db.GetUserByUsername(req.Username)
		if err == nil {
			// err is nil, i.e. we found a match
			return echo.NewHTTPError(http.StatusBadRequest, "Username already taken")
		}
		if err == db.ErrNotFound {
			// pass
		} else {
			return util.GenericServerError(err)
		}
	}

	// Already validated
	user.Username = req.Username
	user.Roles = req.Roles
	err = db.Save(user)
	if err != nil {
		return util.ServerError("Failed to save user", err)
	}

	return c.JSON(http.StatusOK, user.ToDTO())
}

func handleDeleteUser(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

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

	err = db.HardDelete(user)
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

	err = db.HardDelete(agent)
	if err != nil {
		return util.ServerError("Failed to delete agent", err)
	}

	return c.JSON(http.StatusOK, "ok")
}
