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
)

func HookAttackTemplateEndpoints(api *echo.Group) {
	api.GET("/all", func(c echo.Context) error {
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
	})

	api.POST("/create", func(c echo.Context) error {
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
			Description:     req.Description,
			CreatedByUserID: user.ID,
		})
		if err != nil {
			return util.ServerError("Failed to create new attack template", err)
		}

		return c.JSON(http.StatusOK, newAttackTemplate.ToDTO())
	})

	api.DELETE("/:id", func(c echo.Context) error {
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
			if !(user.HasRole(roles.UserRoleAdmin) || attackTemplate.CreatedByUserID.String() == user.ID.String()) {
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

		if !(user.HasRole(roles.UserRoleAdmin) || attackTemplateSet.CreatedByUserID.String() == user.ID.String()) {
			return echo.ErrForbidden
		}

		err = db.DeleteAttackTemplateSet(id)
		if err != nil {
			return util.ServerError("Failed to delete attack template", err)
		}

		return c.JSON(http.StatusOK, "ok")
	})
}
