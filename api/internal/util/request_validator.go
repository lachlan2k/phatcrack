package util

import (
	"net/http"
	"regexp"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/resources"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
)

func BindAndValidate[DTO interface{}](c echo.Context) (DTO, error) {
	var req DTO
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return req, echo.NewHTTPError(http.StatusUnsupportedMediaType, "Only application/json content type is supported")
	}
	if err := c.Bind(&req); err != nil {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}
	err := c.Validate(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

type RequestValidator struct {
	Validator  *validator.Validate
	translator ut.Translator
}

func NewRequestValidator() *RequestValidator {
	v := &RequestValidator{}

	eng := english.New()
	uni := ut.New(eng, eng)
	translator, _ := uni.GetTranslator("en")
	v.Validator = validator.New()
	v.translator = translator
	_ = en.RegisterDefaultTranslations(v.Validator, translator)

	// Register other validators
	v.Validator.RegisterValidation("userroles", func(fl validator.FieldLevel) bool {
		providedRoles, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}

		return roles.AreRolesAssignable(providedRoles)
	})

	standardNameRegex := regexp.MustCompile(`^[\w \-\.']*$`)

	v.Validator.RegisterValidation("standardname", func(fl validator.FieldLevel) bool {
		name, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return standardNameRegex.MatchString(name)
	})

	usernameRegex := regexp.MustCompile(`^[\w\.@\-_]*$`)

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

	return v
}

func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, v.makeValidationError(err))
	}
	return nil
}

func (v *RequestValidator) makeValidationError(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	res := ""

	for i := range errs {
		res += errs[i].Translate(v.translator) + "."

		if i < len(errs)-1 {
			res += "\n"
		}
	}

	return res
}
