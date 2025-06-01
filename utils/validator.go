package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()

	Validator.RegisterValidation("password", validatePassword)
}

func validatePassword(filedLevel validator.FieldLevel) bool {
	password := filedLevel.Field().String()

	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-0]`).MatchString(password)

	return hasNumber && hasLower && hasUpper && hasSpecial
}
