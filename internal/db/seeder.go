package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand"

	"github.com/jaswdr/faker/v2"

	"github.com/harshvse/go-api/internal/store"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	// Let's create users first
	users := GenerateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user: ", err)
			return
		}
	}
	tx.Commit()

	posts := GeneratePosts(users, 100)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts: ", err)
			return
		}
	}

	comments := GenerateComments(posts, 100)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comments", err)
			return
		}
	}
}

func GenerateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	fake := faker.New()

	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: fake.Person().Name(),
			Email:    fake.Internet().Email(),
		}
	}
	return users
}

func GeneratePosts(users []*store.User, n int) []*store.Post {
	posts := make([]*store.Post, n)
	fake := faker.New()
	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]
		tags := make([]string, 3)
		for j := 0; j < 3; j++ {
			tags[j] = fake.Lorem().Word()
		}
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   fake.Lorem().Sentence(10),
			Content: fake.Lorem().Paragraph(100),
			Tags:    tags,
		}
	}
	return posts
}

func GenerateComments(posts []*store.Post, n int) []*store.Comment {
	comments := make([]*store.Comment, n)
	fake := faker.New()
	for i := 0; i < n; i++ {
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			UserID:  post.UserID,
			PostID:  post.ID,
			Content: fake.Lorem().Paragraph(100),
		}
	}
	return comments
}
