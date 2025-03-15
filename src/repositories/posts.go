package repositories

import (
	"api/src/models"
	"context"
	"fmt"
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

func (repository posts) Update(id string, fields map[string]interface{}) (*models.Posts, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	query := "UPDATE posts SET"
	args := []interface{}{}
	argID := 1

	for key, value := range fields {
		if argID > 1 {
			query += ","
		}
		query += fmt.Sprintf(" %s = $%d", key, argID)
		args = append(args, value)
		argID++
	}

	query += fmt.Sprintf(", updated_at = NOW()")

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, title, content, user_id, created_at, updated_at", argID)
	args = append(args, id)
	var updatedPost models.Posts
	err := repository.db.QueryRow(context.Background(), query, args...).Scan(&updatedPost.ID, &updatedPost.Title, &updatedPost.Content, &updatedPost.UserID, &updatedPost.CreatedAt, &updatedPost.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

func (repository posts) FindById(id string) (*models.Posts, error) {
	var post models.Posts
	err := repository.db.QueryRow(context.Background(), "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (repository posts) FindManyByUserId(user_id string, limit int, offset int, filters map[string]interface{}) ([]models.Posts, error) {
	query := "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE user_id = $1"
	args := []interface{}{user_id}
	argID := 2

	for key, value := range filters {
		query += fmt.Sprintf(" AND %s ILIKE $%d", key, argID)
		args = append(args, fmt.Sprintf("%%%s%%", value))
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := repository.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Posts
	for rows.Next() {
		var post models.Posts
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository posts) Delete(id string) error {
	_, err := repository.db.Exec(context.Background(), "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
