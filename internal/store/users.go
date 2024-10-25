package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username,email,password)
	VALUES($1,$2,$3) RETURNING id,created_at,updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userId int64) (*User, error) {
	query := `SELECT id,email,username,created_at,updated_at FROM users WHERE id=($1)`
	var user User

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
