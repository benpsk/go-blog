package handlers

import (
	"net/http"
	"strconv"

	"github.com/benpsk/go-blog/internal/components"
	"github.com/benpsk/go-blog/internal/components/post"
	"github.com/benpsk/go-blog/internal/models"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	var pageData models.HomeResponse
	posts, err := h.service.Home()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	auth, _ := r.Context().Value("user").(models.AuthUser)
	pageData.Data = *posts
	pageData.User = auth
	home := components.Home(&pageData)
	render(w, r, "Home Page", home)
}

func (h *handler) View(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.service.PostDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view := post.View(res)
	render(w, r, "Detail Page", view)
}
