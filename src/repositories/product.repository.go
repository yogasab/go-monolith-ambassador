package repositories

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]*models.Product, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{DB: DB}
}

func (r *productRepository) FindAll() ([]*models.Product, error) {
	var products []*models.Product
	if err := r.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
