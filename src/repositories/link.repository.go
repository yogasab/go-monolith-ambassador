package repositories

import (
	"errors"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type LinkRepository interface {
	FindAllUser(UserID int) ([]*models.Link, error)
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
