package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LockerClusterAssignmentRepository struct {
	db *database.DB
}

func NewLockerClusterAssignmentRepository(db *database.DB) *LockerClusterAssignmentRepository {
	return &LockerClusterAssignmentRepository{db: db}
}

func (r *LockerClusterAssignmentRepository) FindByPackageID(packageID uuid.UUID) (*domain.LockerClusterAssignment, error) {
	var assignment domain.LockerClusterAssignment
	result := r.db.Preload("Package").Preload("Locker").First(&assignment, "package_id = ?", packageID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &assignment, result.Error
}

func (r *LockerClusterAssignmentRepository) FindByLockerID(lockerID uuid.UUID) ([]domain.LockerClusterAssignment, error) {
	var assignments []domain.LockerClusterAssignment
	result := r.db.Preload("Package").Preload("Locker").Where("locker_id = ?", lockerID).Order("clustered_at ASC, id ASC").Find(&assignments)
	return assignments, result.Error
}
