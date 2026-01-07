package validator

import (
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a slice of validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range v {
		msgs = append(msgs, e.Field+": "+e.Message)
	}
	return strings.Join(msgs, "; ")
}

// Get returns the singleton validator instance
func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// Use JSON tag names for field names in errors
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom password validation
		_ = validate.RegisterValidation("strongpassword", validateStrongPassword)
	})
	return validate
}

// validateStrongPassword checks for at least one uppercase, one lowercase, and one digit
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// Struct validates a struct and returns formatted errors
func Struct(s interface{}) ValidationErrors {
	err := Get().Struct(s)
	if err == nil {
		return nil
	}

	var errors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(err.Field()),
			Message: msgForTag(err),
		})
	}

	return errors
}

// msgForTag returns a human-readable message for a validation tag
func msgForTag(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "invalid email format"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "alphanum":
		return field + " can only contain letters and numbers"
	case "strongpassword":
		return "password must contain at least one uppercase letter, one lowercase letter, and one digit"
	default:
		return fe.Error()
	}
}

// SanitizeString trims whitespace and removes control characters
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
