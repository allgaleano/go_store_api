package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID			int64	`json:"id"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	Password	string	`json:"-"`
	CreatedAt	string	`json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context) error {
	return nil
}
