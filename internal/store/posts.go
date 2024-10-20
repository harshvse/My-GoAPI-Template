package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags)
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostStore) Get(ctx context.Context, userId int64) error {
	query := `SELECT * FROM posts WHERE user_id $1`
	var post Post
	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Tags,
	)
	if err != nil {
		return err
	}
	return nil
}
