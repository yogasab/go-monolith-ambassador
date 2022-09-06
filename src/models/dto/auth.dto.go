package dto

type RegisterDTO struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	IsAmbassador bool   `json:"is_ambassador"`
}

type LoginDTO struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateProfileDTO struct {
	ID        int
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type UpdateProfilePassword struct {
	ID       int
	Password string `json:"password" validate:"required,min=8"`
}
