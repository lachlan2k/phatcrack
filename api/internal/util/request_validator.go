package util

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
)

func BindAndValidate[DTO interface{}](c echo.Context) (DTO, error) {
	var req DTO
	if err := c.Bind(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Bad request").SetInternal(err)
	}
	err := c.Validate(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

type RequestValidator struct {
	Validator *validator.Validate
}

func (v *RequestValidator) Init() {
	// Register other validators
	v.Validator.RegisterValidation("userroles", func(fl validator.FieldLevel) bool {
		roles, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}

		for _, role := range roles {
			found := false
			for _, allowedRole := range apitypes.UserSignupRoles {
				if role == allowedRole {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	})
}

// Validate Data
func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
