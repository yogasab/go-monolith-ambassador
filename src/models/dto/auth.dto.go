package dto

type RegisterDTO struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	IsAmbassador bool   `json:"is_ambassador" validate:"required"`
}

type LoginDTO struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
