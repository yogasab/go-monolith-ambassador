package dto

type CreateLinkDTO struct {
	UserID   int
	Products []int `json:"products" validate:"required"`
}
