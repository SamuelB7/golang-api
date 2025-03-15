package controllers

import (
	"api/src/auth"
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

func PostCreate(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractUserID(r)
	if err != nil {
		responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var post models.Posts
	if err = json.Unmarshal(body, &post); err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to unmarshal JSON"})
		return
	}

	err = Validate.Struct(post)
	if err != nil {
		responses.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation failed: %v", err)})
		return
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
	if err != nil || page < 0 {
		page = 0
	}

	if page > 0 {
		page = page - 1
	}

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
	posts, err := repository.FindManyByUserId(userID, limit, page, filters)
	if err != nil {
		responses.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch posts"})
		return
	}
	responses.JsonResponse(w, http.StatusOK, posts)
}

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
