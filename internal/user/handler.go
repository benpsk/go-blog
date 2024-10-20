package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/benpsk/go-blog/pkg"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Handler {
	return &Handler{
		db: db,
	}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func generateToken() (string, error) {
	b := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *Handler) storeSession(user_id int) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return &Session{}, err
	}
	ctx := context.Background()
	var exists bool
	var sessionId int
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id=$1)`
	err = h.db.QueryRow(context.Background(), query, user_id).Scan(&exists)
	if err != nil {
		return &Session{}, err
	}
	if exists {
		err = h.db.QueryRow(ctx, "UPDATE sessions SET token=$1 WHERE user_id=$2 RETURNING id", token, user_id).Scan(&sessionId)
		if err != nil {
			return &Session{}, err
		}
	} else {
		err = h.db.QueryRow(ctx, "INSERT INTO sessions (user_id, token) VALUES($1, $2) RETURNING id", user_id, token).Scan(&sessionId)
		if err != nil {
			return &Session{}, err
		}
	}
	return &Session{
		Id:     sessionId,
		Token:  token,
		UserId: user_id,
	}, nil
}

func (h *Handler) setSession(user_id int, w http.ResponseWriter) error {
	session, err := h.storeSession(user_id)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    session.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,  // Prevents JS access to the Cookie
		Secure:   false, // set to true when using https
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	pkg.RenderTemplate(w, r, "internal/templates/auth/login.html", &pkg.PageData{})
}
func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pkg.PageData{})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var pageData pkg.PageData
	input := &UserInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	validationErrors := validateLogin(input)
	if len(validationErrors) > 0 {
		pageData.Errors = validationErrors
		pkg.RenderTemplate(w, r, "internal/templates/auth/login.html", &pageData)
		return
	}
	ctx := context.Background()
	var user User
	err := h.db.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", input.Email).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		pageData.Error = "Select query error: " + err.Error()
		pkg.RenderTemplate(w, r, "internal/templates/auth/login.html", &pageData)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		pageData.Error = "Password does not match"
		pkg.RenderTemplate(w, r, "internal/templates/auth/login.html", &pageData)
		return
	}
	err = h.setSession(user.Id, w)
	if err != nil {
		pageData.Error = "Set session error: " + err.Error()
		pkg.RenderTemplate(w, r, "internal/templates/auth/login.html", &pageData)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var pageData pkg.PageData
	input := &UserInput{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}
	validationErrors := validateRegister(input)
	if len(validationErrors) > 0 {
		pageData.Errors = validationErrors
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	ctx := context.Background()
	var exists bool
	var err error
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err = h.db.QueryRow(context.Background(), query, input.Email).Scan(&exists)
	if err != nil {
		pageData.Error = "User exist query error: " + err.Error()
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	if exists {
		pageData.Error = "User already exists"
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	password, err := hashPassword(input.Password)
	if err != nil {
		pageData.Error = "Hash password failed"
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	input.Password = password
	// insert to table
	var user_id int
	err = h.db.QueryRow(ctx, `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`, input.Name, input.Email, input.Password).Scan(&user_id)
	if err != nil {
		pageData.Error = "Insert query error: " + err.Error()
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	err = h.setSession(user_id, w)
	if err != nil {
		pageData.Error = "Set session error: " + err.Error()
		pkg.RenderTemplate(w, r, "internal/templates/auth/register.html", &pageData)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,  // Prevents JS access to the Cookie
		Secure:   false, // set to true when using https
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
