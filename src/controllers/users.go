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
	"strconv"

	"github.com/gorilla/mux"
)

// Users godoc
// @Summary Create a new user
// @Description Create a new user with name, email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.User true "User data"
// @Success 201 {object} map[string]interface{} "Returns user_id and success message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
// @Security ApiKeyAuth
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

	err = Validate.Struct(user)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation failed: %v", err)})
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)

	userExists, err := repository.FindByEmail(user.Email)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check if user exists"})
		return
	}

	if userExists != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "User already exists!"})
		return
	}

	userId, err := repository.Create(user)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		return
	}

	responses.JsonResponse(w, http.StatusCreated, map[string]interface{}{"message": "User created successfully", "user_id": userId})
}

// Users godoc
// @Summary Get all users
// @Description Get all users with pagination and filtering
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to return (default 10)"
// @Param page query int false "Page number (default 1)"
// @Param name query string false "Filter by name"
// @Param email query string false "Filter by email"
// @Success 200 {array} models.User "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
// @Security ApiKeyAuth
func UserGetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 0 {
		page = 0
	}

	if page > 0 {
		page = page - 1
	}

	filters := make(map[string]interface{})
	if name := queryParams.Get("name"); name != "" {
		filters["name"] = name
	}

	if email := queryParams.Get("email"); email != "" {
		filters["email"] = email
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)
	users, err := repository.FindMany(limit, page, filters)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve users"})
		return
	}

	responses.JsonResponse(w, http.StatusOK, users)
}

// Users godoc
// @Summary Get a user by ID
// @Description Get the details of a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
// @Security ApiKeyAuth
func UserGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)
	user, err := repository.FindById(id)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
		return
	}

	if user == nil {
		responses.JsonResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	responses.JsonResponse(w, http.StatusOK, user)
}

// Users godoc
// @Summary Update a user
// @Description Update a user with the provided fields
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body map[string]interface{} true "Fields to update"
// @Success 200 {object} models.User "Updated user"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
// @Security ApiKeyAuth
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var fields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)
	updatedUser, err := repository.Update(id, fields)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		return
	}

	responses.JsonResponse(w, http.StatusOK, updatedUser)
}

// Users godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "Success message and deleted user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
func UserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewUsersRepository(db)
	deletedUser, err := repository.Delete(id)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete the user"})
		return
	}

	responses.JsonResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully", "user_id": deletedUser})
}
