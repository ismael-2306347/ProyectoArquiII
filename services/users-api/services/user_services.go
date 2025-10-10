package services

import (
	"fmt"
	"users-api/domain"
	"users-api/repositories"
)

type UserService interface {
	CreateUser(dto domain.CreateUserDTO) (domain.UserResponseDTO, error)
	GetUserByID(id uint) (domain.UserResponseDTO, error)
	Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(dto domain.CreateUserDTO) (domain.UserResponseDTO, error) {
	// Implementar

	fmt.Printf("todo: implementar services ")
	return domain.UserResponseDTO{}, nil
}

func (s *userService) GetUserByID(id uint) (domain.UserResponseDTO, error) {
	// Implementar
	return domain.UserResponseDTO{}, nil
}

func (s *userService) Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error) {
	// Implementar
	return domain.LoginResponseDTO{}, nil
}
