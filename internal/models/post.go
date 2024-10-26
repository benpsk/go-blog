package models

import (
	"time"
)

type Post struct {
	Id        int
	Title     string
	Excerpt   string
	Body      string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	User      UserResponse
}

type PageData struct {
	Data   interface{}
	Error  string
	Errors map[string]string
	User   AuthUser
}

type HomeResponse struct {
	User AuthUser
	Data []Post
}
