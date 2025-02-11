package repositories

import (
	"api/src/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type posts struct {
	db *pgx.Conn
}

func NewPostsRepository(db *pgx.Conn) *posts {
	return &posts{db}
}

func (repository posts) Create(post models.Posts) (string, error) {
	var postId string

	tx, err := repository.db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = tx.QueryRow(context.Background(), "INSERT INTO posts (id, title, content, user_id, created_at) VALUES (uuid_generate_v4(), $1, $2, $3, CURRENT_TIMESTAMP) RETURNING id", post.Title, post.Content, post.UserID).Scan(&postId)
	if err != nil {
		tx.Rollback(context.Background())
		return "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}

	return postId, nil
}
