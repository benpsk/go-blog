package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/benpsk/go-blog/internal/components/layouts"
	"github.com/benpsk/go-blog/internal/models"
	"github.com/benpsk/go-blog/internal/services"
)

type handler struct {
	service *services.Service
}

func New(s *services.Service) *handler {
	return &handler{
		service: s,
	}
}

func render(w http.ResponseWriter, r *http.Request, title string, component templ.Component) {
	auth, _ := r.Context().Value("user").(models.AuthUser)
	layouts.Layout(title, auth, component).Render(context.Background(), w)
}
