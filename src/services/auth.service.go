package services

import (
	"errors"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type AuthService interface {
	Register(dto *dto.RegisterDTO) (*models.User, error)
	Login(dto *dto.LoginDTO) (*models.User, error)
	GetProfile(ID int) (*models.User, error)
	UpdateProfile(dto *dto.UpdateProfileDTO) (*models.User, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) Register(dto *dto.RegisterDTO) (*models.User, error) {
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

func (s *authService) Login(dto *dto.LoginDTO) (*models.User, error) {
	registeredUser, err := s.userRepository.FindByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if ok := registeredUser.ComparePassword(registeredUser.Password, dto.Password); !ok {
		return nil, errors.New("invalid credentials")
	}
	return registeredUser, nil
}

func (s *authService) GetProfile(ID int) (*models.User, error) {
	user, err := s.userRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) UpdateProfile(dto *dto.UpdateProfileDTO) (*models.User, error) {
	// check user if it exists
	user, err := s.userRepository.FindByID(dto.ID)
	if err != nil {
		return nil, err
	}
	// check if email is registered
	if _, err = s.userRepository.FindByEmail(dto.Email); err != nil {
		return nil, err
	}
	user.Email = dto.Email
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	updatedUser, err := s.userRepository.Update(dto.ID, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
