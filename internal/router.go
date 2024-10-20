package internal

import (
	"net/http"

	"github.com/benpsk/go-blog/internal/dashboard"
	"github.com/benpsk/go-blog/internal/middlewares"
	"github.com/benpsk/go-blog/internal/post"
	"github.com/benpsk/go-blog/internal/user"
	"github.com/jackc/pgx/v5"
)

func Router(db *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()
	post := post.NewHandler(db)
	user := user.New(db)
	dashboard := dashboard.NewHandler(db)

	mux.HandleFunc("GET /", middlewares.Auth(post.Home, db))
	mux.HandleFunc("GET /post/{id}", post.View)

	mux.HandleFunc("GET /login", user.ShowLogin)
	mux.HandleFunc("POST /login", user.Login)
	mux.HandleFunc("POST /logout", user.Logout)
	mux.HandleFunc("GET /register", user.ShowRegister)
	mux.HandleFunc("POST /register", user.Register)

	mux.HandleFunc("GET /dashboard", middlewares.Auth(dashboard.Dashboard, db))
	return mux
}
