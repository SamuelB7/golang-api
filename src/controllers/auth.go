package controllers

import (
	"api/src/auth"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// AuthRequest represents the authentication request data
// @Description Authentication request with email and password
// @Schema AuthRequest
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// SignInRequest represents the registration request data
// @Description Registration request with name, email, and password
// @Schema SignInRequest
type SignInRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// Login godoc
// @Summary Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Login credentials"
// @Success 200 "Returns authentication token"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var authRequest AuthRequest
	if err = json.Unmarshal(body, &authRequest); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to unmarshal JSON"})
		return
	}

	err = Validate.Struct(authRequest)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Validation failed"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}
	defer db.Close(r.Context())

	repository := repositories.NewUsersRepository(db)

	user, err := repository.FindByEmail(authRequest.Email)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if user == nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "User or password is incorrect"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))

	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "User or password is incorrect"})
		return
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}

	responses.JsonResponse(w, http.StatusOK, map[string]string{"token": token})
}

// SignIn godoc
// @Summary Register a new user
// @Description Creates a new user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignInRequest true "Registration information"
// @Success 201 "User created successfully and returns authentication token"
// @Failure 400 "Bad request"
// @Failure 409 "User already exists"
// @Failure 500 "Internal server error"
// @Router /sign-in [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var signInRequest SignInRequest
	if err = json.Unmarshal(body, &signInRequest); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to unmarshal JSON"})
		return
	}

	err = Validate.Struct(signInRequest)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Validation failed"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signInRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate password hash"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}
	defer db.Close(r.Context())

	repository := repositories.NewUsersRepository(db)

	userExists, err := repository.FindByEmail(signInRequest.Email)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if userExists != nil {
		responses.JsonResponse(w, http.StatusConflict, map[string]string{"error": "User already registered"})
		return
	}

	userID, err := repository.Create(models.User{
		Name:     signInRequest.Name,
		Email:    signInRequest.Email,
		Password: string(passwordHash),
	})
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}

	responses.JsonResponse(w, http.StatusCreated, map[string]string{"token": token})
}
