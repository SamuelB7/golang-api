package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	repository := repositories.NewUsersRepository(db)
	userId, err := repository.Create(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created: %s", userId)))
}

func UserGetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get one user"))
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
