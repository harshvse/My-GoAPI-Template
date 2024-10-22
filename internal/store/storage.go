package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetPostByID(context.Context, int64) ([]PostWithComments, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentsStore{db: db},
	}
}
