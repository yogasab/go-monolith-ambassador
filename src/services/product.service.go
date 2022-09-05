package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type ProductService interface {
	GetProducts() ([]*models.Product, error)
	GetProduct(ID int) (*models.Product, error)
	UpdateProduct(dto *dto.UpdateProductDTO) (*models.Product, error)
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

func (s *productService) GetProduct(ID int) (*models.Product, error) {
	product, err := s.productRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) UpdateProduct(dto *dto.UpdateProductDTO) (*models.Product, error) {
	product, err := s.productRepository.FindByID(dto.ID)
	if err != nil {
		return nil, err
	}

	product.Title = dto.Title
	product.Description = dto.Description
	product.Image = dto.ImageURL
	updatedProduct, err := s.productRepository.Update(dto.ID, product)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}
