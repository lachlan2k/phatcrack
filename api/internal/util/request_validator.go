package util

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func BindAndValidate[DTO interface{}](c echo.Context) (DTO, error) {
	var req DTO
	if err := c.Bind(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Bad request").SetInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return req, err
	}
	return req, nil
}

type RequestValidator struct {
	Validator *validator.Validate
}

func (v *RequestValidator) Init() {
	// Register other validators
	v.Validator.RegisterValidation("userrole", func(fl validator.FieldLevel) bool {
		switch fl.Field().String() {
		case "admin", "standard":
			return true
		default:
			return false
		}
	})
}

// Validate Data
func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
