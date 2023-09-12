package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func HookUserEndpoints(api *echo.Group) {
	api.GET("/all", func(c echo.Context) error {
		users, err := db.GetAllUsers()
		if err != nil {
			return util.ServerError("Failed to fetch users", err)
		}

		userDTOs := []apitypes.UserMinimalDTO{}
		for _, user := range users {
			userDTOs = append(userDTOs, user.ToMinimalDTO())
		}

		return c.JSON(http.StatusOK, apitypes.UsersGetAllResponseDTO{
			Users: userDTOs,
		})
	})
}
