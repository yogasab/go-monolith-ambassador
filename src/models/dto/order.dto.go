package dto

type OrderRankingDTO struct {
	Name    string  `json:"name"`
	Revenue float64 `json:"revenue"`
}

type CreateOrderDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Address   string `json:"address" validate:"required"`
	Country   string `json:"country" validate:"required"`
	City      string `json:"city" validate:"required"`
	ZIP       string `json:"zip" validate:"required"`
	Code      string `json:"code" validate:"required"`
	Products  []map[string]int
}
