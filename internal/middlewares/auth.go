package middlewares

import (
	"context"
	"net/http"
	"regexp"

	"github.com/benpsk/go-blog/internal/models"
	"github.com/jackc/pgx/v5"
)

var postRegex = regexp.MustCompile(`^/post/([^/]+)$`)

func Auth(next http.HandlerFunc, db *pgx.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			if r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}
			if match := postRegex.FindStringSubmatch(r.URL.Path); match != nil {
				next.ServeHTTP(w, r)
				return
			}
		}
		// check the token
		var userId int
		err = db.QueryRow(context.Background(), "SELECT user_id FROM sessions WHERE token=$1", token.Value).Scan(&userId)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var user models.AuthUser
		err = db.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id=$1", userId).
			Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			http.Error(w, "User not found"+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
