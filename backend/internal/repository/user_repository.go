package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, "email = ?", email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &user, result.Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User

	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &user, result.Error
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":                user.Name,
		"email":               user.Email,
		"phone":               user.Phone,
		"password_hash":       user.PasswordHash,
	})
	return result.Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	result := r.db.Delete(&domain.User{}, "id = ?", id)
	return result.Error
}