package repositories

import (
	"api/src/models"
	"context"
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
