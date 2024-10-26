package validations

import (
	"net/mail"
	"strings"

	"github.com/benpsk/go-blog/internal/models"
)

func Register(input *models.UserInput) (*models.RegisterError, bool) {
	var msg models.RegisterError
	ok := true
	if strings.TrimSpace(input.Name) == "" {
		msg.Name = "Name is required"
		ok = false
	}
	if len(input.Name) > 100 {
		msg.Name = "Name must not exceed 100 characters"
		ok = false
	}
	if strings.TrimSpace(input.Email) == "" {
		msg.Email = "Email is required"
		ok = false
	} else if _, err := mail.ParseAddress(input.Email); err != nil {
		msg.Email = "Invalid email format"
		ok = false
	}
	if len(input.Email) > 100 {
		msg.Email = "Email must not exceed 100 characters"
		ok = false
	}

	if strings.TrimSpace(input.Password) == "" {
		msg.Password = "Password is required"
		ok = false
	}
	if len(input.Password) > 100 {
		msg.Password = "Password must not exceed 100 characters"
		ok = false
	}
	if strings.TrimSpace(input.ConfirmPassword) == "" {
		msg.ConfirmPassword = "Confirm password is required"
		ok = false
	}
	if len(input.ConfirmPassword) > 100 {
		msg.ConfirmPassword = "Confirm password must not exceed 100 characters"
		ok = false
	}
	if input.Password != input.ConfirmPassword {
		msg.Password = "Password does not match"
		ok = false
	}
	return &msg, ok
}

func Login(input *models.UserInput) (*models.LoginError, bool) {
	var msg models.LoginError
	ok := true
	if strings.TrimSpace(input.Email) == "" {
		msg.Email = "Email is required"
		ok = false
	} else if _, err := mail.ParseAddress(input.Email); err != nil {
		msg.Email = "Invalid email format"
		ok = false
	}
	if len(input.Email) > 100 {
		msg.Email = "Email must not exceed 100 characters"
		ok = false
	}

	if strings.TrimSpace(input.Password) == "" {
		msg.Password = "Password is required"
		ok = false
	}
	if len(input.Password) > 100 {
		msg.Password = "Password must not exceed 100 characters"
		ok = false
	}
	return &msg, ok
}
