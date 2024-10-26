package handlers

import (
	"net/http"

	"github.com/benpsk/go-blog/internal/components/auth"
	"github.com/benpsk/go-blog/internal/models"
	"github.com/benpsk/go-blog/internal/validations"
)

func (h *handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	login := auth.Login("", &models.LoginError{})
	render(w, r, "Login Page", login)
}
func (h *handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	register := auth.Register("", &models.RegisterError{})
	render(w, r, "Register Page", register)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &models.UserInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	validationErrors, ok := validations.Login(input)
	if !ok {
		login := auth.Login("", validationErrors)
		render(w, r, "Login Page", login)
		return
	}
	err := h.service.Login(w, input)
	if err != nil {
		login := auth.Login(err.Error(), validationErrors)
		render(w, r, "Login Page", login)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	input := &models.UserInput{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}
	validationErrors, ok := validations.Register(input)
	if !ok {
		register := auth.Register("", validationErrors)
		render(w, r, "Register Page", register)
		return
	}
	err := h.service.Register(w, input)
	if err != nil {
		register := auth.Register(err.Error(), validationErrors)
		render(w, r, "Register Page", register)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.service.Logout(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
