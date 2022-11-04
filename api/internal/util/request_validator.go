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
	// register other validators
}

// Validate Data
func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
