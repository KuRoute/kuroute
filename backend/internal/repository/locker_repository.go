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

func (r *LockerRepository) FindAvailableByHubID(hubID uuid.UUID) ([]domain.Locker, error) {
	var lockers []domain.Locker
	result := r.db.
		Where("hub_id = ?", hubID).
		Where("NOT EXISTS (?)", r.db.Model(&domain.LockerScan{}).Select("1").Where("locker_scan.locker_id = locker.id AND locker_scan.scanned_out_at IS NULL")).
		Order("label ASC").
		Find(&lockers)
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

func (r *LockerRepository) FindActivePackages(lockerID uuid.UUID) ([]domain.Package, error) {
	var packages []domain.Package
	result := r.db.
		Joins("JOIN locker_scan ON locker_scan.package_id = package.id").
		Where("locker_scan.locker_id = ? AND locker_scan.scanned_out_at IS NULL", lockerID).
		Order("package.id ASC").
		Find(&packages)
	return packages, result.Error
}

func (r *LockerRepository) CreateClusters(assignments []domain.CreateClusterPackagesRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		lockerIDs := make(map[uuid.UUID]struct{}, len(assignments))
		packageIDs := make(map[uuid.UUID]struct{})

		type validatedAssignment struct {
			Locker   domain.Locker
			Request  domain.CreateClusterPackagesRequest
			Packages []uuid.UUID
		}

		validated := make([]validatedAssignment, 0, len(assignments))

		for _, assignment := range assignments {
			if _, exists := lockerIDs[assignment.LockerID]; exists {
				return ErrClusterLockerConflict
			}
			lockerIDs[assignment.LockerID] = struct{}{}

			var locker domain.Locker
			if err := tx.
				Where("id = ?", assignment.LockerID).
				First(&locker).Error; err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return sql.ErrNoRows
				}

				return err
			}

			for _, packageID := range assignment.PackageIDs {
				if _, exists := packageIDs[packageID]; exists {
					return ErrClusterPackageConflict
				}
				packageIDs[packageID] = struct{}{}

				var packageItem domain.Package

				if err := tx.
					Where(
						"id = ? AND hub_id = ? AND status = ?",
						packageID,
						locker.HubID,
						domain.PackageStatusReceived,
					).
					First(&packageItem).Error; err != nil {

					if errors.Is(err, gorm.ErrRecordNotFound) {
						return sql.ErrNoRows
					}

					return err
				}

				var activeScan int64
				if err := tx.
					Model(&domain.LockerScan{}).
					Where(
						"package_id = ? AND scanned_out_at IS NULL",
						packageID,
					).
					Count(&activeScan).Error; err != nil {
					return err
				}

				if activeScan != 0 {
					return ErrClusterPackageActive
				}
			}

			validated = append(validated, validatedAssignment{
				Locker:   locker,
				Request:  assignment,
				Packages: assignment.PackageIDs,
			})
		}

		for _, item := range validated {
			clusteredAt := time.Now()

			if err := tx.
				Model(&domain.Locker{}).
				Where("id = ?", item.Locker.ID).
				Updates(map[string]interface{}{
					"cluster_area": item.Request.ClusterArea,
					"clustered_at": clusteredAt,
				}).Error; err != nil {
				return err
			}

			for _, packageID := range item.Packages {
				assignment := &domain.LockerClusterAssignment{
					LockerID:    item.Locker.ID,
					PackageID:   packageID,
					ClusterArea: item.Request.ClusterArea,
					ClusteredAt: clusteredAt,
				}

				if err := tx.Create(assignment).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

var (
	ErrClusterLockerConflict  = errors.New("locker is assigned more than once")
	ErrClusterPackageConflict = errors.New("package is assigned more than once")
	ErrClusterPackageActive   = errors.New("package already assigned to an active locker")
)
