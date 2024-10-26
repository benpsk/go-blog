package handlers

import (
	"net/http"

	"github.com/benpsk/go-blog/internal/components/dashboard"
)

func (h *handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := dashboard.Dashboard()
	render(w, r, "Dashboard", dashboard)
}
