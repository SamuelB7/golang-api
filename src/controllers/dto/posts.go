package dto

type PostCreateDTO struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
