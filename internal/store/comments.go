package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id,user_id, content) VALUES ($1,$2,$3) RETURNING id,created_at`
	err := s.db.QueryRowContext(
		ctx, query, comment.PostID, comment.UserID, comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

type PostWithComments struct {
	UserID           int64    `json:"user_id"`
	PostID           int64    `json:"post_id"`
	PostTitle        string   `json:"post_title"`
	PostContent      string   `json:"post_content"`
	PostTags         []string `json:"tags"`
	PostCreatedAt    string   `json:"post_created_at"`
	PostUpdatedAt    string   `json:"post_updated_at"`
	CommentID        int64    `json:"comment_id"`
	CommentContent   string   `json:"comment_content"`
	CommentCreatedAt string   `json:"comment_created_at"`
}

func (s *CommentsStore) GetPostByID(ctx context.Context, postId int64) (*[]PostWithComments, error) {
	var postWithComments PostWithComments
	query := `SELECT
		c.user_id as user_id,
		c.post_id as post_id,
		p.title as post_title,
		p.content as post_content,
		p.tags as tags,
		p.created_at as post_created_at,
		p.updated_at as post_updated_at,
		c.id as comment_id,
		c.content as comment_content,
		c.created_at as content_created_at
	FROM comments AS c 
	INNER JOIN posts as p 
	ON p.id=c.post_id 
	WHERE c.post_id=($1)`
	rows, err := s.db.QueryContext(ctx, query, postId)

	if err != nil {
		return nil, err
	}
	return &postWithComments, nil
}
