package repositories

import (
	"errors"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type LinkRepository interface {
	FindAllUser(UserID int) ([]*models.Link, error)
	Create(link *models.Link) (*models.Link, error)
	FindByUserID(UserID int) ([]*models.Link, error)
	FindByCode(code string) (*models.Link, error)
}

type linkRepository struct {
	DB *gorm.DB
}

func NewLinkRepository(DB *gorm.DB) LinkRepository {
	return &linkRepository{DB: DB}
}

func (r *linkRepository) FindAllUser(UserID int) ([]*models.Link, error) {
	var links []*models.Link
	if err := r.DB.Where("user_id = ?", UserID).Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, errors.New("user links is not found")
	}
	return links, nil
}

func (r *linkRepository) Create(link *models.Link) (*models.Link, error) {
	if err := r.DB.Create(&link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

func (r *linkRepository) FindByUserID(UserID int) ([]*models.Link, error) {
	var links []*models.Link
	if err := r.DB.Find(&links, models.Link{UserID: uint(UserID)}).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, errors.New("links to user is not found")
	}
	return links, nil
}

func (r *linkRepository) FindByCode(code string) (*models.Link, error) {
	var link *models.Link
	if err := r.DB.Preload("User").Preload("Products").Where("code = ?", code).First(&link).Error; err != nil {
		if link.ID == 0 {
			return nil, errors.New("link is not found")
		}
		return nil, err
	}

	return link, nil
}
