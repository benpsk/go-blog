package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/benpsk/go-blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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

func (s *Service) storeSession(user_id int) (*models.Session, error) {
	token, err := generateToken()
	if err != nil {
		return &models.Session{}, err
	}
	ctx := context.Background()
	var exists bool
	var sessionId int
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id=$1)`
	err = s.Db.QueryRow(context.Background(), query, user_id).Scan(&exists)
	if err != nil {
		return &models.Session{}, err
	}
	if exists {
		err = s.Db.QueryRow(ctx, "UPDATE sessions SET token=$1 WHERE user_id=$2 RETURNING id", token, user_id).Scan(&sessionId)
		if err != nil {
			return &models.Session{}, err
		}
	} else {
		err = s.Db.QueryRow(ctx, "INSERT INTO sessions (user_id, token) VALUES($1, $2) RETURNING id", user_id, token).Scan(&sessionId)
		if err != nil {
			return &models.Session{}, err
		}
	}
	return &models.Session{
		Id:     sessionId,
		Token:  token,
		UserId: user_id,
	}, nil
}

func (s *Service) setSession(user_id int, w http.ResponseWriter) error {
	session, err := s.storeSession(user_id)
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

func (s *Service) Login(w http.ResponseWriter, input *models.UserInput) error {
	ctx := context.Background()
	var user models.User
	err := s.Db.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", input.Email).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		return errors.New("Select query error: " + err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return errors.New("Password does not match")
	}
	err = s.setSession(user.Id, w)
	return err
}

func (s *Service) Register(w http.ResponseWriter, input *models.UserInput) error {
	ctx := context.Background()
	var exists bool
	var err error
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err = s.Db.QueryRow(ctx, query, input.Email).Scan(&exists)
	if err != nil {
		return errors.New("User exist query error: " + err.Error())
	}
	if exists {
		return errors.New("User already exists")
	}
	password, err := hashPassword(input.Password)
	if err != nil {
		return errors.New("Hash password failed")
	}
	input.Password = password
	// insert to table
	var user_id int
	err = s.Db.QueryRow(ctx, `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`, input.Name, input.Email, input.Password).Scan(&user_id)
	if err != nil {
		return errors.New("Insert query error: " + err.Error())
	}
	err = s.setSession(user_id, w)
	if err != nil {
		return errors.New("Set session error: " + err.Error())
	}
	return nil
}

func (s *Service) Logout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,  // Prevents JS access to the Cookie
		Secure:   false, // set to true when using https
		SameSite: http.SameSiteStrictMode,
	})
}
