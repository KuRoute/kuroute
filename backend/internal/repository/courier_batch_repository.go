package repository

import (
	"errors"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourierBatchRepository struct {
	db *database.DB
}

func NewCourierBatchRepository(db *database.DB) *CourierBatchRepository {
	return &CourierBatchRepository{db: db}
}

func (r *CourierBatchRepository) FindByID(id uuid.UUID) (*domain.CourierBatch, error) {
	var courierBatch domain.CourierBatch
	result := r.db.Preload("BatchAssignment").First(&courierBatch, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &courierBatch, result.Error
}

func (r *CourierBatchRepository) FindByHubID(hubID uuid.UUID, status *domain.BatchStatus, batchDate *time.Time) ([]domain.CourierBatch, error) {
	var batches []domain.CourierBatch
	query := r.db.Preload("BatchAssignment").Joins("JOIN batch_assignment ON batch_assignment.id = courier_batch.batch_assignment_id").Joins("JOIN locker ON locker.id = batch_assignment.locker_id").Where("locker.hub_id = ?", hubID)

	if status != nil {
		query = query.Where("courier_batch.status = ?", *status)
	}
	if batchDate != nil {
		query = query.Where("batch_assignment.batch_date = ?", batchDate.Format("2006-01-02"))
	}

	result := query.Order("batch_assignment.batch_date DESC, courier_batch.started_at DESC").Find(&batches)
	return batches, result.Error
}

func (r *CourierBatchRepository) FindByCourierID(courierID uuid.UUID) ([]domain.CourierBatch, error) {
	var batches []domain.CourierBatch
	result := r.db.Preload("BatchAssignment").Joins("JOIN batch_assignment ON batch_assignment.id = courier_batch.batch_assignment_id").Where("batch_assignment.courier_id = ?", courierID).Order("batch_assignment.batch_date DESC, courier_batch.started_at DESC").Find(&batches)
	return batches, result.Error
}

func (r *CourierBatchRepository) FindByStatus(status domain.BatchStatus) ([]domain.CourierBatch, error) {
	var batches []domain.CourierBatch
	result := r.db.
		Preload("BatchAssignment").
		Where("status = ?", status).
		Order("id ASC").
		Find(&batches)
	return batches, result.Error
}
