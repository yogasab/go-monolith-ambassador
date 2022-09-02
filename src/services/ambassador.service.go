package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type AmbassadorService interface {
	GetAmbassadors() ([]*models.User, error)
}

type ambassadorService struct {
	userRepository repositories.UserRepository
}

func NewAmbassadorService(userRepository repositories.UserRepository) AmbassadorService {
	return &ambassadorService{userRepository: userRepository}
}

func (s *ambassadorService) GetAmbassadors() ([]*models.User, error) {
	ambassadors, err := s.userRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return ambassadors, nil
}
