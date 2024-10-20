package dashboard

import (
	"net/http"

	"github.com/benpsk/go-blog/pkg"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	pkg.RenderTemplate(w, r, "internal/templates/dashboard/index.html", &pkg.PageData{})
}
