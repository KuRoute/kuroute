package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LockerScanRepository struct {
	db *database.DB
}

func NewLockerScanRepository(db *database.DB) *LockerScanRepository {
	return &LockerScanRepository{db: db}
}

func (r *LockerScanRepository) Create(scan *domain.LockerScan) error {
	result := r.db.Create(scan)
	return result.Error
}

func (r *LockerScanRepository) FindByID(id uuid.UUID) (*domain.LockerScan, error) {
	var scan domain.LockerScan
	result := r.db.First(&scan, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &scan, result.Error
}

func (r *LockerScanRepository) FindActiveByPackageID(packageID uuid.UUID) (*domain.LockerScan, error) {
	var scan domain.LockerScan
	result := r.db.Where("package_id = ? AND scanned_out_at IS NULL", packageID).First(&scan)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &scan, result.Error
}

func (r *LockerScanRepository) FindByLockerID(lockerID uuid.UUID) ([]domain.LockerScan, error) {
	var scans []domain.LockerScan
	result := r.db.Where("locker_id = ?", lockerID).Order("scanned_at DESC").Find(&scans)
	return scans, result.Error
}

func (r *LockerScanRepository) FindByPackageID(packageID uuid.UUID) ([]domain.LockerScan, error) {
	var scans []domain.LockerScan
	result := r.db.Where("package_id = ?", packageID).Order("scanned_at DESC").Find(&scans)
	return scans, result.Error
}

func (r *LockerScanRepository) FindByCourierToday(courierID uuid.UUID) ([]domain.LockerScan, error) {
	var scans []domain.LockerScan
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	result := r.db.Where("scanned_by = ? AND scanned_at BETWEEN ? AND ?", courierID, startOfDay, endOfDay).
		Order("scanned_at DESC").
		Find(&scans)
	return scans, result.Error
}

func (r *LockerScanRepository) Update(scan *domain.LockerScan) error {
	result := r.db.Model(&domain.LockerScan{}).Where("id = ?", scan.ID).Updates(scan)
	return result.Error
}
