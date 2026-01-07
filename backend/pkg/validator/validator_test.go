package validator

import (
	"testing"
)

type testRegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=128,strongpassword"`
}

type testLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name     string
		req      testRegisterRequest
		wantErrs int
	}{
		{
			name:     "all valid",
			req:      testRegisterRequest{Email: "test@example.com", Username: "johndoe", Password: "Password123"},
			wantErrs: 0,
		},
		{
			name:     "all empty",
			req:      testRegisterRequest{Email: "", Username: "", Password: ""},
			wantErrs: 3,
		},
		{
			name:     "invalid email",
			req:      testRegisterRequest{Email: "invalid", Username: "johndoe", Password: "Password123"},
			wantErrs: 1,
		},
		{
			name:     "username too short",
			req:      testRegisterRequest{Email: "test@example.com", Username: "ab", Password: "Password123"},
			wantErrs: 1,
		},
		{
			name:     "username with special chars",
			req:      testRegisterRequest{Email: "test@example.com", Username: "john@doe", Password: "Password123"},
			wantErrs: 1,
		},
		{
			name:     "password too short",
			req:      testRegisterRequest{Email: "test@example.com", Username: "johndoe", Password: "Pass1"},
			wantErrs: 1,
		},
		{
			name:     "password no uppercase",
			req:      testRegisterRequest{Email: "test@example.com", Username: "johndoe", Password: "password123"},
			wantErrs: 1,
		},
		{
			name:     "password no lowercase",
			req:      testRegisterRequest{Email: "test@example.com", Username: "johndoe", Password: "PASSWORD123"},
			wantErrs: 1,
		},
		{
			name:     "password no digit",
			req:      testRegisterRequest{Email: "test@example.com", Username: "johndoe", Password: "Passworddd"},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Struct(&tt.req)
			if len(errs) != tt.wantErrs {
				t.Errorf("Struct() got %d errors, want %d errors. Errors: %v", len(errs), tt.wantErrs, errs)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name     string
		req      testLoginRequest
		wantErrs int
	}{
		{
			name:     "valid login",
			req:      testLoginRequest{Email: "test@example.com", Password: "anypassword"},
			wantErrs: 0,
		},
		{
			name:     "empty email",
			req:      testLoginRequest{Email: "", Password: "anypassword"},
			wantErrs: 1,
		},
		{
			name:     "invalid email format",
			req:      testLoginRequest{Email: "notanemail", Password: "anypassword"},
			wantErrs: 1,
		},
		{
			name:     "empty password",
			req:      testLoginRequest{Email: "test@example.com", Password: ""},
			wantErrs: 1,
		},
		{
			name:     "all empty",
			req:      testLoginRequest{Email: "", Password: ""},
			wantErrs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Struct(&tt.req)
			if len(errs) != tt.wantErrs {
				t.Errorf("Struct() got %d errors, want %d errors. Errors: %v", len(errs), tt.wantErrs, errs)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"trim spaces", "  hello  ", "hello"},
		{"remove newlines", "hello\nworld", "helloworld"},
		{"remove tabs", "hello\tworld", "helloworld"},
		{"remove control chars", "hello\x00world", "helloworld"},
		{"normal string", "hello world", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidatorSingleton(t *testing.T) {
	v1 := Get()
	v2 := Get()
	if v1 != v2 {
		t.Error("Get() should return the same validator instance")
	}
}
