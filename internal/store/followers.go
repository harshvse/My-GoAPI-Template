package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, follower_id int64, followed_id int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *FollowerStore) UnFollow(ctx context.Context, follower_id int64, followed_id int64) error {
	query := `DELETE FROM followers WHERE user_id = ($1) AND follower_id = ($2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}
