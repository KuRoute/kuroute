package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourierProfileRepository struct {
	db *database.DB
}

func NewCourierProfileRepository(db *database.DB) *CourierProfileRepository {
	return &CourierProfileRepository{db: db}
}

func (r *CourierProfileRepository) CreateCourierProfile(profile *domain.CourierProfile) error {
	result := r.db.Create(profile)
	return result.Error
}

func (r *CourierProfileRepository) GetCourierProfileByUserID(userId uuid.UUID) (*domain.CourierProfile, error) {
	var profile domain.CourierProfile

	result := r.db.First(&profile, "id = ?", userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &profile, result.Error
}

func (r *CourierProfileRepository) UpdateCourierProfile(profile *domain.CourierProfile) error {
	result := r.db.Model(&domain.CourierProfile{}).Where("id = ?", profile.UserID).Updates(map[string]interface{}{
		"user_id":        profile.UserID,
		"vehicle_type":   profile.VehicleType,
	})
	return result.Error
}

func (r *CourierProfileRepository) DeleteCourierProfile(userId uuid.UUID) error {
	result := r.db.Delete(&domain.CourierProfile{}, "id = ?", userId)
	return result.Error
}