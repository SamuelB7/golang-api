package repositories

import (
	"api/src/models"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	db *pgx.Conn
}

func NewUsersRepository(db *pgx.Conn) *users {
	return &users{db}
}

func (repository users) Create(user models.User) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(passwordHash)
	tx, err := repository.db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var userId string
	err = tx.QueryRow(context.Background(), "INSERT INTO users (id, name, email, password, created_at) VALUES (uuid_generate_v4(), $1, $2, $3, CURRENT_TIMESTAMP) RETURNING id", user.Name, user.Email, user.Password).Scan(&userId)
	if err != nil {
		tx.Rollback(context.Background())
		return "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (repository users) FindById(id string) (*models.User, error) {
	var user models.User
	err := repository.db.QueryRow(context.Background(), "SELECT id, name, email, created_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repository users) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repository.db.QueryRow(context.Background(), "SELECT id, name, email, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repository users) FindMany(limit int, offset int, filters map[string]interface{}) ([]models.User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE 1=1"
	args := []interface{}{}
	argID := 1

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

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository users) Update(id string, fields map[string]interface{}) (*models.User, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	query := "UPDATE users SET"
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

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, email, created_at", argID)
	args = append(args, id)

	var updatedUser models.User
	err := repository.db.QueryRow(context.Background(), query, args...).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (repository users) Delete(id string) (string, error) {
	query := "DELETE FROM users WHERE id = $1 RETURNING id"
	var deletedUser string
	err := repository.db.QueryRow(context.Background(), query, id).Scan(&deletedUser)
	if err != nil {
		return "", err
	}

	return deletedUser, nil
}
