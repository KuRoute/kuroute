package repository

import (
	"errors"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteRepository struct {
	db *database.DB
}

func NewRouteRepository(db *database.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) CreateRoute(route *domain.Route) error {
	result := r.db.Create(route)
	return result.Error
}

func (r *RouteRepository) CreateRouteWithStops(route *domain.Route, stops []domain.RouteStop) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(route).Error; err != nil {
			return err
		}

		for i := range stops {
			stops[i].RouteID = route.ID
			if err := tx.Create(&stops[i]).Error; err != nil {
				return err
			}
		}

		return tx.Model(&domain.CourierBatch{}).Where("id = ?", route.CourierBatchID).Update("status", domain.BatchStatusRouteReady).Error
	})
}

func (r *RouteRepository) FindByID(id uuid.UUID) (*domain.Route, error) {
	var route domain.Route
	result := r.db.Preload("Stops").Preload("Stops.Package").First(&route, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &route, result.Error
}

func (r *RouteRepository) FindByHubID(hubID uuid.UUID, batchDate *time.Time) ([]domain.Route, error) {
	var routes []domain.Route
	query := r.db.Preload("Stops").Preload("Stops.Package").Joins("JOIN courier_batch ON courier_batch.id = route.courier_batch_id").Joins("JOIN batch_assignment ON batch_assignment.id = courier_batch.batch_assignment_id").Joins("JOIN locker ON locker.id = batch_assignment.locker_id").Where("locker.hub_id = ?", hubID)

	if batchDate != nil {
		query = query.Where("batch_assignment.batch_date = ?", batchDate.Format("2006-01-02"))
	}

	result := query.Order("batch_assignment.batch_date DESC, route.computed_at DESC").Find(&routes)
	return routes, result.Error
}
