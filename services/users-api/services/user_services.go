package services

import (
	"context"
	"fmt"
	"users-api/domain"
	"users-api/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, dto domain.CreateUserDTO) (domain.UserResponseDTO, error)
	GetUserByID(id uint) (domain.UserResponseDTO, error)
	Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{repository: repository}
}

// user_services.go
func (s *userService) CreateUser(ctx context.Context, dto domain.CreateUserDTO) (domain.UserResponseDTO, error) {
	created, err := s.repository.Create(ctx, dto) // <-- pasar dto real
	if err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to create user: %w", err)
	}
	return created, nil
}

func (s *userService) GetUserByID(id uint) (domain.UserResponseDTO, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (s *userService) Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error) {
	// Implementar
	return domain.LoginResponseDTO{}, nil
}
