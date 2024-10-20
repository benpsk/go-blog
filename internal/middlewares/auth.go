package middlewares

import (
	"context"
	"net/http"

	"github.com/benpsk/go-blog/pkg"
	"github.com/jackc/pgx/v5"
)

func Auth(next http.HandlerFunc, db *pgx.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("token")
		if err != nil {
			if r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// check the token
		var userId int
		err = db.QueryRow(context.Background(), "SELECT user_id FROM sessions WHERE token=$1", token.Value).Scan(&userId)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var user pkg.AuthUser
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
