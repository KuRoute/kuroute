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

var ErrDeliveryStatusTransition = errors.New("delivery status transition failed")

type DeliveryReportRepository struct {
	db *database.DB
}

func NewDeliveryReportRepository(db *database.DB) *DeliveryReportRepository {
	return &DeliveryReportRepository{db: db}
}

func (r *DeliveryReportRepository) Create(report *domain.DeliveryReport) error {
	return r.db.Create(report).Error
}

func (r *DeliveryReportRepository) CreateAndUpdateStatus(report *domain.DeliveryReport, routeStopID uuid.UUID, packageID uuid.UUID, routeID uuid.UUID, courierBatchID uuid.UUID, newStopStatus domain.PackageStatus) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(report).Error; err != nil {
			return err
		}
		result := tx.Model(&domain.RouteStop{}).Where("id = ? AND status = ?", routeStopID, domain.PackageStatusInDelivery).Update("status", newStopStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrDeliveryStatusTransition
		}
		result = tx.Model(&domain.Package{}).Where("id = ? AND status = ?", packageID, domain.PackageStatusInDelivery).Update("status", newStopStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrDeliveryStatusTransition
		}

		var remaining int64
		if err := tx.Model(&domain.RouteStop{}).Where("route_id = ? AND status NOT IN (?, ?)", routeID, domain.PackageStatusDelivered, domain.PackageStatusFailed).Count(&remaining).Error; err != nil {
			return err
		}
		if remaining == 0 {
			if err := tx.Model(&domain.CourierBatch{}).Where("id = ?", courierBatchID).Updates(map[string]any{"status": domain.BatchStatusCompleted, "completed_at": time.Now()}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DeliveryReportRepository) FindByID(id uuid.UUID) (*domain.DeliveryReport, error) {
	var report domain.DeliveryReport
	result := r.db.First(&report, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &report, result.Error
}

func (r *DeliveryReportRepository) FindByRouteStopID(routeStopID uuid.UUID) ([]domain.DeliveryReport, error) {
	var reports []domain.DeliveryReport
	result := r.db.Where("route_stop_id = ?", routeStopID).Order("reported_at DESC").Find(&reports)
	return reports, result.Error
}

func (r *DeliveryReportRepository) FindByPackageID(packageID uuid.UUID) ([]domain.DeliveryReport, error) {
	var reports []domain.DeliveryReport
	result := r.db.Table("delivery_report").
		Joins("JOIN route_stop ON route_stop.id = delivery_report.route_stop_id").
		Where("route_stop.package_id = ?", packageID).
		Order("delivery_report.reported_at DESC").
		Find(&reports)
	return reports, result.Error
}

func (r *DeliveryReportRepository) FindByRouteID(routeID uuid.UUID) ([]domain.DeliveryReport, error) {
	var reports []domain.DeliveryReport
	result := r.db.Table("delivery_report").
		Joins("JOIN route_stop ON route_stop.id = delivery_report.route_stop_id").
		Where("route_stop.route_id = ?", routeID).
		Order("delivery_report.reported_at DESC").
		Find(&reports)
	return reports, result.Error
}

func (r *DeliveryReportRepository) CountCompletedStops(routeID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.Model(&domain.RouteStop{}).Where("route_id = ? AND status IN (?, ?, ?)", routeID, domain.PackageStatusDelivered, domain.PackageStatusFailed, domain.PackageStatusInDelivery).Count(&count)
	return count, result.Error
}

func (r *DeliveryReportRepository) UpdateRouteStopStatus(routeStopID uuid.UUID, status domain.PackageStatus) error {
	return r.db.Model(&domain.RouteStop{}).Where("id = ?", routeStopID).Update("status", status).Error
}

func (r *DeliveryReportRepository) UpdatePackageStatus(packageID uuid.UUID, status domain.PackageStatus) error {
	return r.db.Model(&domain.Package{}).Where("id = ?", packageID).Update("status", status).Error
}

func (r *DeliveryReportRepository) UpdateCourierBatchCompleted(courierBatchID uuid.UUID, completedAt time.Time) error {
	return r.db.Model(&domain.CourierBatch{}).Where("id = ?", courierBatchID).Updates(map[string]any{"status": domain.BatchStatusCompleted, "completed_at": completedAt}).Error
}

func (r *DeliveryReportRepository) FindRouteStopByID(id uuid.UUID) (*domain.RouteStop, error) {
	var stop domain.RouteStop
	result := r.db.Preload("Package").First(&stop, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &stop, result.Error
}
