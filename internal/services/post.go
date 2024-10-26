package services

import (
	"context"
	"errors"

	"github.com/benpsk/go-blog/internal/models"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Home() (*[]models.Post, error) {
	rows, err := s.Db.Query(context.Background(), "SELECT p.id, p.title, p.excerpt, p.body, p.created_at, p.updated_at , u.name, u.email FROM posts p INNER JOIN users u ON p.user_id = u.id")
	defer rows.Close()
	if err != nil {
		return &[]models.Post{}, errors.New("Select query error: " + err.Error())
	}
	posts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Post, error) {
		var post models.Post
		var userRes models.UserResponse
		err := row.Scan(&post.Id, &post.Title, &post.Excerpt, &post.Body, &post.CreatedAt, &post.UpdatedAt, &userRes.Name, &userRes.Email)
		if err != nil {
			return models.Post{}, err
		}
		post.User = userRes
		return post, err
	})
	if err != nil {
		return &[]models.Post{}, errors.New("Row collections error: " + err.Error())
	}
	return &posts, nil
}

func (s *Service) PostDetail(id int) (*models.Post, error) {
	var model models.Post
	var userRes models.UserResponse
	err := s.Db.QueryRow(context.Background(), "SELECT p.id, p.title, p.excerpt, p.body, u.name, u.email FROM posts p INNER JOIN users u ON p.user_id = u.id WHERE p.id=$1", id).Scan(&model.Id, &model.Title, &model.Excerpt, &model.Body, &userRes.Name, &userRes.Email)
	if err != nil {
		return &models.Post{}, err
	}
	model.User = userRes
	return &model, nil
}
