package repositories

import (
	"errors"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(ID int) (*models.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userRepository{DB: DB}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.DB.Create(&user).Error; err != nil {
		// if strings.Contains(err.Error(), "1062") {
		// 	return nil, err
		// }
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.DB.Find(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user is not registered")
	}
	return user, nil
}

func (r *userRepository) FindByID(ID int) (*models.User, error) {
	var user *models.User
	if err := r.DB.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
