package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	passwordMinLength = 6
	passwordMaxLength = 16
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
)

var (
	lengthRegexp    = regexp.MustCompile(fmt.Sprintf(`^.{%d,%d}$`, passwordMinLength, passwordMaxLength))
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%^&*]{%d,}`, passwordMinSymbol))
)

type CustomValidator struct {
	v             *validator.Validate
	passwordError error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}
	return cv
}

func (cv *CustomValidator) Validate(i interface{}) (err error) {
	err = cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]
		err = cv.NewValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return
}

func (cv *CustomValidator) NewValidationError(field string, value interface{}, tag string, param string) (err error) {
	switch tag {
	case "required":
		err = fmt.Errorf("field %s is required", field)
	case "email":
		err = fmt.Errorf("field %s must be a valid email address", field)
	case "password":
		err = cv.passwordError
	case "min":
		err = fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		err = fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		err = fmt.Errorf("field %s is invalid", field)
	}
	return
}

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {

	if fl.Field().Kind() != reflect.String {
		cv.passwordError = fmt.Errorf("field %s must be a string ", fl.FieldName())
	}

	fieldValue := fl.Field().String()
	var ok bool
	if ok = lengthRegexp.MatchString(fieldValue); !ok {
		cv.passwordError = fmt.Errorf("Field %s length must be between %d and %d characters", fl.FieldName(), passwordMinLength, passwordMaxLength)
	} else if ok = lowerCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwordError = fmt.Errorf("Field %s must contain at least %d lowercase letter", fl.FieldName(), passwordMinLower)
	} else if ok = upperCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwordError = fmt.Errorf("Field %s must contain at least %d uppercase letter", fl.FieldName(), passwordMinUpper)
	} else if ok = digitRegexp.MatchString(fieldValue); !ok {
		cv.passwordError = fmt.Errorf("Field %s must contain at least %d digit", fl.FieldName(), passwordMinDigit)
	} else if ok = symbolRegexp.MatchString(fieldValue); !ok {
		cv.passwordError = fmt.Errorf("Field %s must contain at least %d special symbols", fl.FieldName(), passwordMinSymbol)
	}
	if !ok {
		return false
	}
	return true
}
