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

type BatchAssignmentRepository struct {
	db *database.DB
}

func NewBatchAssignmentRepository(db *database.DB) *BatchAssignmentRepository {
	return &BatchAssignmentRepository{db: db}
}

func (r *BatchAssignmentRepository) Create(assignment *domain.BatchAssignment) error {
	result := r.db.Create(assignment)
	return result.Error
}

func (r *BatchAssignmentRepository) FindByID(id uuid.UUID) (*domain.BatchAssignment, error) {
	var assignment domain.BatchAssignment
	result := r.db.First(&assignment, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &assignment, result.Error
}

func (r *BatchAssignmentRepository) Update(assignment *domain.BatchAssignment) error {
	result := r.db.Model(&domain.BatchAssignment{}).Where("id = ?", assignment.ID).Updates(assignment)
	return result.Error
}

func (r *BatchAssignmentRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.BatchAssignment{}, "id = ?", id)
	return result.Error
}

func (r *BatchAssignmentRepository) FindByLockerID(lockerID uuid.UUID) ([]domain.BatchAssignment, error) {
	var assignments []domain.BatchAssignment
	result := r.db.Where("locker_id = ?", lockerID).Order("assigned_at DESC").Find(&assignments)
	return assignments, result.Error
}

func (r *BatchAssignmentRepository) FindByCourierID(courierID uuid.UUID) ([]domain.BatchAssignment, error) {
	var assignments []domain.BatchAssignment
	result := r.db.Where("courier_id = ?", courierID).Order("batch_date DESC, batch_sequence ASC").Find(&assignments)
	return assignments, result.Error
}

func (r *BatchAssignmentRepository) FindByHubID(hubID uuid.UUID, batchDate *time.Time, courierID *uuid.UUID) ([]domain.BatchAssignment, error) {
	var assignments []domain.BatchAssignment
	query := r.db.Joins("JOIN locker ON locker.id = batch_assignment.locker_id").Where("locker.hub_id = ?", hubID)

	if batchDate != nil {
		query = query.Where("batch_assignment.batch_date = ?", batchDate.Format("2006-01-02"))
	}
	if courierID != nil {
		query = query.Where("batch_assignment.courier_id = ?", courierID)
	}

	result := query.Order("batch_assignment.batch_date DESC, batch_assignment.batch_sequence ASC").Find(&assignments)
	return assignments, result.Error
}

func (r *BatchAssignmentRepository) FindByCourierAndDate(courierID uuid.UUID, batchDate time.Time) ([]domain.BatchAssignment, error) {
	var assignments []domain.BatchAssignment
	result := r.db.Where("courier_id = ? AND batch_date = ?", courierID, batchDate.Format("2006-01-02")).Order("batch_sequence ASC").Find(&assignments)
	return assignments, result.Error
}

func (r *BatchAssignmentRepository) FindByLockerAndDateSequence(lockerID uuid.UUID, batchDate time.Time, batchSequence int16) (*domain.BatchAssignment, error) {
	var assignment domain.BatchAssignment
	result := r.db.Where("locker_id = ? AND batch_date = ? AND batch_sequence = ?", lockerID, batchDate.Format("2006-01-02"), batchSequence).First(&assignment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &assignment, result.Error
}

func (r *BatchAssignmentRepository) FindByCourierAndDateSequence(courierID uuid.UUID, batchDate time.Time, batchSequence int16) (*domain.BatchAssignment, error) {
	var assignment domain.BatchAssignment
	result := r.db.Where("courier_id = ? AND batch_date = ? AND batch_sequence = ?", courierID, batchDate.Format("2006-01-02"), batchSequence).First(&assignment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &assignment, result.Error
}

func (r *BatchAssignmentRepository) CreateCourierBatch(courierBatch *domain.CourierBatch) error {
	result := r.db.Create(courierBatch)
	return result.Error
}

func (r *BatchAssignmentRepository) FindCourierBatchByAssignmentID(batchAssignmentID uuid.UUID) (*domain.CourierBatch, error) {
	var courierBatch domain.CourierBatch
	result := r.db.Where("batch_assignment_id = ?", batchAssignmentID).First(&courierBatch)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &courierBatch, result.Error
}

func (r *BatchAssignmentRepository) UpdateCourierBatch(courierBatch *domain.CourierBatch) error {
	result := r.db.Model(&domain.CourierBatch{}).Where("id = ?", courierBatch.ID).Updates(courierBatch)
	return result.Error
}

func (r *BatchAssignmentRepository) FindRouteByAssignmentID(batchAssignmentID uuid.UUID) (*domain.Route, error) {
	var route domain.Route
	result := r.db.Table("route").
		Joins("JOIN courier_batch ON courier_batch.id = route.courier_batch_id").
		Joins("JOIN batch_assignment ON batch_assignment.id = courier_batch.batch_assignment_id").
		Where("batch_assignment.id = ?", batchAssignmentID).
		First(&route)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &route, result.Error
}

func (r *BatchAssignmentRepository) FindRouteStopsByRouteID(routeID uuid.UUID) ([]domain.RouteStop, error) {
	var stops []domain.RouteStop
	result := r.db.Where("route_id = ?", routeID).Order("stop_order ASC").Find(&stops)
	return stops, result.Error
}
