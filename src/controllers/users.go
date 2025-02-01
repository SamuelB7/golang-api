package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to unmarshal JSON"})
		return
	}

	if err := validate.Struct(user); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation failed: %v", err)})
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)
	userId, err := repository.Create(user)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		return
	}

	responses.JsonResponse(w, http.StatusCreated, map[string]interface{}{"message": "User created successfully", "user_id": userId})
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
