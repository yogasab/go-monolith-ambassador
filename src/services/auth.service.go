package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type AuthService interface {
	RegisterAdmin(dto *dto.RegisterDTO) (*models.User, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) RegisterAdmin(dto *dto.RegisterDTO) (*models.User, error) {
	newUser := &models.User{}
	newUser.FirstName = dto.FirstName
	newUser.LastName = dto.LastName
	newUser.Email = dto.Email
	newUser.IsAmbassador = dto.IsAmbassador
	newUser.Password = newUser.HashPassword(dto.Password)

	u, err := s.userRepository.Create(newUser)
	if err != nil {
		return nil, err
	}
	return u, nil
}
