package controllers

import (
	"api/src/auth"
	"api/src/controllers/dto"
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
	"github.com/rs/zerolog/log"
)

// Posts godoc
// @Summary Create a new post
// @Description Create a new post with title and content
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body dto.PostCreateDTO true "Post data"
// @Success 201 "Returns post_id and success message"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /posts [post]
// @Security ApiKeyAuth
func PostCreate(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.ExtractUserID(r)
	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var postDTO dto.PostCreateDTO
	if err = json.Unmarshal(body, &postDTO); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to unmarshal JSON"})
		return
	}

	err = Validate.Struct(postDTO)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation failed: %v", err)})
		return
	}

	post := models.Posts{
		Title:   postDTO.Title,
		Content: postDTO.Content,
		UserID:  userId,
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewPostsRepository(db)

	postID, err := repository.Create(post)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
		return
	}

	responses.JsonResponse(w, http.StatusCreated, map[string]interface{}{"message": "Post created successfully", "post_id": postID})
}

// Posts godoc
// @Summary Get all posts by user ID
// @Description Get all posts created by the authenticated user with pagination and filtering
// @Tags Posts
// @Accept json
// @Produce json
// @Param limit query int false "Number of posts to return (default 10)"
// @Param page query int false "Page number (default 1)"
// @Param title query string false "Filter by title"
// @Param content query string false "Filter by content"
// @Success 200 {array} models.Posts "List of posts"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /posts-by-user [get]
// @Security ApiKeyAuth
func PostGetAllByUserId(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 0 {
		limit = 10
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if title := queryParams.Get("title"); title != "" {
		filters["title"] = title
	}

	if content := queryParams.Get("content"); content != "" {
		filters["content"] = content
	}

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewPostsRepository(db)
	posts, err := repository.FindManyByUserId(userID, limit, offset, filters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch posts")
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch posts"})
		return
	}
	responses.JsonResponse(w, http.StatusOK, posts)
}

// Posts godoc
// @Summary Get a post by ID
// @Description Get the details of a specific post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.Posts "Post details"
// @Failure 500 "Internal server error"
// @Router /posts/{id} [get]
func PostGetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewPostsRepository(db)
	post, err := repository.FindById(id)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to find post"})
		return
	}
	responses.JsonResponse(w, http.StatusOK, post)
}

// Posts godoc
// @Summary Update a post
// @Description Update a post with the provided fields
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body dto.PostCreateDTO true "Fields to update" example({"title": "Updated Post Title", "content": "This is the updated content."})
// @Success 200 {object} models.Posts "Updated post"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /posts/{id} [put]
// @Security ApiKeyAuth
func PostUpdate(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractUserID(r)
	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

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

	repository := repositories.NewPostsRepository(db)
	updatedPost, err := repository.Update(id, fields)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update post"})
		return
	}
	responses.JsonResponse(w, http.StatusOK, updatedPost)
}

// Posts godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 204 "No content"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /posts/{id} [delete]
// @Security ApiKeyAuth
func PostDelete(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractUserID(r)
	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	db, err := database.Connect()
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return
	}

	repository := repositories.NewPostsRepository(db)
	err = repository.Delete(id)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete post"})
		return
	}
	responses.JsonResponse(w, http.StatusNoContent, nil)
}
