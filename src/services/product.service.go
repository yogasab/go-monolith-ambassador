package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type ProductService interface {
	GetProducts() ([]*models.Product, error)
}

type productService struct {
	productRepository repositories.ProductRepository
}

func NewProductsService(productRepository repositories.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) GetProducts() ([]*models.Product, error) {
	products, err := s.productRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}
