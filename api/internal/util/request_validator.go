package util

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/resources"
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

	standardNameRegex := regexp.MustCompile(`^[\w \-\.']+$`)

	v.Validator.RegisterValidation("standardname", func(fl validator.FieldLevel) bool {
		name, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return standardNameRegex.MatchString(name)
	})

	usernameRegex := regexp.MustCompile(`^[\w\._]+$`)

	v.Validator.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return usernameRegex.MatchString(username)
	})

	v.Validator.RegisterValidation("hashtype", func(fl validator.FieldLevel) bool {
		t, ok := fl.Field().Interface().(int)
		if !ok {
			return false
		}

		_, ok = resources.GetHashTypeMap()[t]
		return ok
	})
}

// Validate Data
func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
