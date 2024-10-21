package store

import "database/sql"

type Comment struct {
	ID      int64  `json:"id"`
	Post_ID int64  `json:"post_id"`
	User_ID int64  `json:"user_id"`
	Content string `json:"content"`
}

type CommentsStore struct {
	db *sql.DB
}
