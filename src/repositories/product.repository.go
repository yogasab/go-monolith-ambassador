package repositories

import (
	"errors"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]*models.Product, error)
	FindByID(ID int) (*models.Product, error)
	Update(ID int, product *models.Product) (*models.Product, error)
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

func (r *productRepository) FindByID(ID int) (*models.Product, error) {
	var product *models.Product

	if err := r.DB.Where("id = ?", ID).First(&product).Error; err != nil {
		if product.ID == 0 {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Update(ID int, product *models.Product) (*models.Product, error) {
	if err := r.DB.Where("id = ?", ID).Save(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
