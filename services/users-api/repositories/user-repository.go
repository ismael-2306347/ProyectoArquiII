package repositories

import (
	"context"
	"errors"
	"fmt"
	"users-api/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.CreateUserDTO) (domain.UserResponseDTO, error)
	GetByID(id uint) (domain.UserResponseDTO, error)
	GetByUsername(username string) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {

	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user domain.CreateUserDTO) (domain.UserResponseDTO, error) {
	// Create the user record in the database
	newUser := domain.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      domain.RoleNormal, // Default role
		// Add other fields as needed
	}
	if err := r.db.Create(&newUser).Error; err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to create user: %w", err)
	}

	response := domain.UserResponseDTO{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Role:      newUser.Role,

		// Add other fields as needed
	}

	return response, nil
}

func (r *userRepository) GetByID(id uint) (domain.UserResponseDTO, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserResponseDTO{}, gorm.ErrRecordNotFound // <-- propagar
		}
		return domain.UserResponseDTO{}, fmt.Errorf("failed to get user by id: %w", err)
	}
	return domain.UserResponseDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}, nil
}

func (r *userRepository) GetByUsername(username string) (domain.User, error) {
	var u domain.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, gorm.ErrRecordNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	return u, nil
}

func (r *userRepository) GetByEmail(email string) (domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, gorm.ErrRecordNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return u, nil
}
