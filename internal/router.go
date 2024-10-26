package internal

import (
	"net/http"

	"github.com/benpsk/go-blog/internal/handlers"
	"github.com/benpsk/go-blog/internal/middlewares"
	"github.com/benpsk/go-blog/internal/services"
	"github.com/jackc/pgx/v5"
)

func Router(db *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()

	service := services.New(db)
	handler := handlers.New(service)

	mux.HandleFunc("GET /{$}", middlewares.Auth(handler.Home, db))
	mux.HandleFunc("GET /post/{id}", middlewares.Auth(handler.View, db))

	mux.HandleFunc("GET /login", handler.ShowLogin)
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /logout", handler.Logout)
	mux.HandleFunc("GET /register", handler.ShowRegister)
	mux.HandleFunc("POST /register", handler.Register)

	mux.HandleFunc("GET /dashboard", middlewares.Auth(handler.Dashboard, db))

	return mux
}
