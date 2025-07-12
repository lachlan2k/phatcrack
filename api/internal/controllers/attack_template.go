package controllers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	"gorm.io/datatypes"
)

func HookAttackTemplateEndpoints(api *echo.Group) {
	api.GET("/all", handleGetAllAttackTemplates)
	api.POST("/create", handleCreateAttackTemplate)
	api.POST("/create-set", handleCreateAttackTemplateSet)
	api.PUT("/:id", handleUpdateAttackTemplate)
	api.DELETE("/:id", handleDeleteAttackTemplate)
}

func handleGetAllAttackTemplates(c echo.Context) error {
	ats, err := db.GetAllAttackTemplates()
	if err != nil {
		return util.GenericServerError(err)
	}
	atsets, err := db.GetAllAttackTemplateSets()
	if err != nil {
		return util.GenericServerError(err)
	}

	dtos := make([]apitypes.AttackTemplateDTO, 0)

	for _, at := range ats {
		dtos = append(dtos, at.ToDTO())
	}

	for _, atset := range atsets {
		dtos = append(dtos, atset.ToDTO())
	}

	return c.JSON(http.StatusOK, apitypes.AttackTemplateGetAllResponseDTO{
		AttackTemplates: dtos,
	})
}

func handleCreateAttackTemplate(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	req, err := util.BindAndValidate[apitypes.AttackTemplateCreateRequestDTO](c)
	if err != nil {
		return err
	}

	newAttackTemplate, err := db.CreateAttackTemplate(&db.AttackTemplate{
		Name:            req.Name,
		HashcatParams:   datatypes.NewJSONType(req.HashcatParams),
		CreatedByUserID: user.ID,
	})
	if err != nil {
		return util.ServerError("Failed to create new attack template", err)
	}

	return c.JSON(http.StatusOK, newAttackTemplate.ToDTO())
}

func handleCreateAttackTemplateSet(c echo.Context) error {
	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	req, err := util.BindAndValidate[apitypes.AttackTemplateCreateSetRequestDTO](c)
	if err != nil {
		return err
	}

	newAttackTemplateSet, err := db.CreateAttackTemplateSet(&db.AttackTemplateSet{
		Name:              req.Name,
		AttackTemplateIDs: req.AttackTemplateIDs,
		CreatedByUserID:   user.ID,
	})
	if err != nil {
		return util.ServerError("Failed to create new attack template set", err)
	}

	return c.JSON(http.StatusOK, newAttackTemplateSet.ToDTO())
}

func handleUpdateAttackTemplate(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	req, err := util.BindAndValidate[apitypes.AttackTemplateUpdateRequestDTO](c)
	if err != nil {
		return err
	}

	switch req.Type {
	case apitypes.AttackTemplateType:
		{
			attackTemplate, err := db.GetAttackTemplate(id)
			if errors.Is(err, db.ErrNotFound) {
				return echo.ErrNotFound
			}
			if err != nil {
				return util.GenericServerError(err)
			}
			allowed := user.HasRole(roles.UserRoleAdmin) || attackTemplate.CreatedByUserID.String() == user.ID.String()
			if !allowed {
				return echo.ErrForbidden
			}
			if req.HashcatParams == nil {
				return echo.ErrBadRequest
			}

			attackTemplate.Name = req.Name
			attackTemplate.HashcatParams = datatypes.NewJSONType(*req.HashcatParams)

			err = db.Save(attackTemplate)
			if err != nil {
				return util.ServerError("Failed to save attack template", err)
			}

			return c.JSON(http.StatusOK, attackTemplate.ToDTO())
		}

	case apitypes.AttackTemplateSetType:
		{
			attackTemplateSet, err := db.GetAttackTemplateSet(id)
			if errors.Is(err, db.ErrNotFound) {
				return echo.ErrNotFound
			}
			if err != nil {
				return util.GenericServerError(err)
			}
			allowed := user.HasRole(roles.UserRoleAdmin) || attackTemplateSet.CreatedByUserID.String() == user.ID.String()
			if !allowed {
				return echo.ErrForbidden
			}
			if req.AttackTemplateIDs == nil {
				return echo.ErrBadRequest
			}

			attackTemplateSet.Name = req.Name
			attackTemplateSet.AttackTemplateIDs = req.AttackTemplateIDs

			err = db.Save(attackTemplateSet)
			if err != nil {
				return util.ServerError("Failed to save attack template set", err)
			}

			return c.JSON(http.StatusOK, attackTemplateSet.ToDTO())
		}

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown type %q", req.Type)
	}
}

func handleDeleteAttackTemplate(c echo.Context) error {
	id := c.Param("id")
	if !util.AreValidUUIDs(id) {
		return echo.ErrBadRequest
	}

	user := auth.UserFromReq(c)
	if user == nil {
		return echo.ErrForbidden
	}

	attackTemplate, err := db.GetAttackTemplate(id)
	if err == nil {
		allowed := user.HasRole(roles.UserRoleAdmin) || attackTemplate.CreatedByUserID.String() == user.ID.String()
		if !allowed {
			return echo.ErrForbidden
		}

		err := db.DeleteAttackTemplate(id)
		if err != nil {
			return util.ServerError("Failed to delete attack template", err)
		}

		return c.JSON(http.StatusOK, "ok")
	}

	attackTemplateSet, err := db.GetAttackTemplateSet(id)
	if errors.Is(err, db.ErrNotFound) {
		return echo.ErrNotFound
	}
	if err != nil {
		return util.GenericServerError(err)
	}

	allowed := user.HasRole(roles.UserRoleAdmin) || attackTemplateSet.CreatedByUserID.String() == user.ID.String()
	if !allowed {
		return echo.ErrForbidden
	}

	err = db.DeleteAttackTemplateSet(id)
	if err != nil {
		return util.ServerError("Failed to delete attack template", err)
	}

	return c.JSON(http.StatusOK, "ok")
}
