package repositories

import (
	"fmt"
	"users-api/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	GetByID(id uint) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user domain.User) (domain.User, error) {
	// Implementar
	fmt.Printf("todo:implementar repository")
	return user, nil
}

func (r *userRepository) GetByID(id uint) (domain.User, error) {
	// Implementar
	return domain.User{}, nil
}

func (r *userRepository) GetByUsername(username string) (domain.User, error) {
	// Implementar
	return domain.User{}, nil
}

func (r *userRepository) GetByEmail(email string) (domain.User, error) {
	// Implementar
	return domain.User{}, nil
}
