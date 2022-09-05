package dto

type UpdateProductDTO struct {
	ID          int
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	ImageURL    string `json:"image_url" validate:"required"`
}
