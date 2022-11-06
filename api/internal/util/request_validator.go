package util

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

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
