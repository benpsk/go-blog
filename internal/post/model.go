package post

import (
	"time"

	"github.com/benpsk/go-blog/internal/user"
)

type Post struct {
	Id        int
	Title     string
	Excerpt   string
	Body      string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	User      user.UserResponse
}
