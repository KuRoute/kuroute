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

func (r *CourierProfileRepository) FindCouriersByHubID(hubID uuid.UUID, vehicleType *domain.VehicleType) ([]domain.CourierProfile, error) {
	var profiles []domain.CourierProfile

	query := r.db.Joins("JOIN users ON users.id = courier_profile.user_id").Where("users.hub_id = ?", hubID)

	if vehicleType != nil {
		query = query.Where("courier_profile.vehicle_type = ?", *vehicleType)
	}

	result := query.Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}

	return profiles, nil
}

func (r *CourierProfileRepository) GetCourierProfileByUserID(userId uuid.UUID) (*domain.CourierProfile, error) {
	var profile domain.CourierProfile

	result := r.db.First(&profile, "user_id = ?", userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &profile, result.Error
}

func (r *CourierProfileRepository) UpdateCourierProfile(profile *domain.CourierProfile) error {
	result := r.db.Model(&domain.CourierProfile{}).Where("user_id = ?", profile.UserID).Updates(map[string]interface{}{
		"vehicle_type":  profile.VehicleType,
		"vehicle_plate": profile.VehiclePlate,
	})
	return result.Error
}

func (r *CourierProfileRepository) DeleteCourierProfile(userId uuid.UUID) error {
	result := r.db.Delete(&domain.CourierProfile{}, "id = ?", userId)
	return result.Error
}
