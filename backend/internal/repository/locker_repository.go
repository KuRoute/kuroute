package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LockerRepository struct {
	db *database.DB
}

func NewLockerRepository(db *database.DB) *LockerRepository {
	return &LockerRepository{db: db}
}

func (r *LockerRepository) Create(locker *domain.Locker) error {
	result := r.db.Create(locker)
	return result.Error
}

func (r *LockerRepository) FindByID(id uuid.UUID) (*domain.Locker, error) {
	var locker domain.Locker
	result := r.db.First(&locker, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &locker, result.Error
}

func (r *LockerRepository) FindByHubID(hubID uuid.UUID) ([]domain.Locker, error) {
	var lockers []domain.Locker
	result := r.db.Where("hub_id = ?", hubID).Find(&lockers)
	return lockers, result.Error
}

func (r *LockerRepository) FindByHubIDAndLabel(hubID uuid.UUID, label string) (*domain.Locker, error) {
	var locker domain.Locker
	result := r.db.Where("hub_id = ? AND label = ?", hubID, label).First(&locker)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &locker, result.Error
}

func (r *LockerRepository) Update(locker *domain.Locker) error {
	result := r.db.Model(&domain.Locker{}).Where("id = ?", locker.ID).Updates(locker)
	return result.Error
}

func (r *LockerRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.Locker{}, "id = ?", id)
	return result.Error
}

func (r *LockerRepository) HasActivePackages(lockerID uuid.UUID) (bool, error) {
	var count int64
	result := r.db.Model(&domain.LockerScan{}).
		Where("locker_id = ? AND scanned_out_at IS NULL", lockerID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
