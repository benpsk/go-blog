package post

import (
	"context"
	"net/http"
	"strconv"

	"github.com/benpsk/go-blog/internal/user"
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

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	var pageData pkg.PageData
	rows, err := h.db.Query(context.Background(), "SELECT p.id, p.title, p.excerpt, p.body, p.created_at, p.updated_at , u.name, u.email FROM posts p INNER JOIN users u ON p.user_id = u.id")
	defer rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	posts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Post, error) {
		var post Post
		var userRes user.UserResponse
		err := row.Scan(&post.Id, &post.Title, &post.Excerpt, &post.Body, &post.CreatedAt, &post.UpdatedAt, &userRes.Name, &userRes.Email)
		if err != nil {
			return Post{}, err
		}
		post.User = userRes
		return post, err
	})
	pageData.Data = posts
	pkg.RenderTemplate(w, r, "internal/templates/index.html", &pageData)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	var pageData pkg.PageData
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var post Post
	var userRes user.UserResponse
	err = h.db.QueryRow(context.Background(), "SELECT p.id, p.title, p.excerpt, p.body, u.name, u.email FROM posts p INNER JOIN users u ON p.user_id = u.id WHERE p.id=$1", id).Scan(&post.Id, &post.Title, &post.Excerpt, &post.Body, &userRes.Name, &userRes.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post.User = userRes
	pageData.Data = post
	pkg.RenderTemplate(w, r, "internal/templates/post/view.html", &pageData)
}
