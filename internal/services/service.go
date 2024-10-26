package services

import "github.com/jackc/pgx/v5"

type Service struct {
	Db *pgx.Conn
}

func New(db *pgx.Conn) *Service {
	return &Service{
		Db: db,
	}
}
