package user

import (
	"net/mail"
	"strings"
)

func validateRegister(input *UserInput) map[string]string {
	msg := make(map[string]string)
	if strings.TrimSpace(input.Name) == "" {
		msg["name"] = "Name is required"
	}
	if len(input.Name) > 100 {
		msg["name"] = "Name must not exceed 100 characters"
	}
	if strings.TrimSpace(input.Email) == "" {
		msg["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(input.Email); err != nil {
		msg["email"] = "Invalid email format"
	}
	if len(input.Email) > 100 {
		msg["email"] = "Email must not exceed 100 characters"
	}

	if strings.TrimSpace(input.Password) == "" {
		msg["password"] = "Password is required"
	}
	if len(input.Password) > 100 {
		msg["password"] = "Password must not exceed 100 characters"
	}
	if strings.TrimSpace(input.ConfirmPassword) == "" {
		msg["confirm_password"] = "Confirm password is required"
	}
	if len(input.ConfirmPassword) > 100 {
		msg["confirm_password"] = "Confirm password must not exceed 100 characters"
	}
	if input.Password != input.ConfirmPassword {
		msg["password"] = "Password does not match"
	}
	return msg
}

func validateLogin(input *UserInput) map[string]string {
	msg := make(map[string]string)
	if strings.TrimSpace(input.Email) == "" {
		msg["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(input.Email); err != nil {
		msg["email"] = "Invalid email format"
	}
	if len(input.Email) > 100 {
		msg["email"] = "Email must not exceed 100 characters"
	}

	if strings.TrimSpace(input.Password) == "" {
		msg["password"] = "Password is required"
	}
	if len(input.Password) > 100 {
		msg["password"] = "Password must not exceed 100 characters"
	}
	return msg
}
