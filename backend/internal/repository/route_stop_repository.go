package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteStopRepository struct {
	db *database.DB
}

func NewRouteStopRepository(db *database.DB) *RouteStopRepository {
	return &RouteStopRepository{db: db}
}

func (r *RouteStopRepository) FindByID(id uuid.UUID) (*domain.RouteStop, error) {
	var stop domain.RouteStop
	result := r.db.Preload("Package").Preload("DeliveryReports", func(db *gorm.DB) *gorm.DB {
		return db.Order("reported_at DESC")
	}).First(&stop, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &stop, result.Error
}

func (r *RouteStopRepository) FindByRouteID(routeID uuid.UUID) ([]domain.RouteStop, error) {
	var stops []domain.RouteStop
	result := r.db.Preload("Package").Where("route_id = ?", routeID).Order("stop_order ASC").Find(&stops)
	return stops, result.Error
}

func (r *RouteStopRepository) FindRouteByStopID(stopID uuid.UUID) (*domain.Route, error) {
	var route domain.Route
	result := r.db.Table("route").
		Joins("JOIN route_stop ON route_stop.route_id = route.id").
		Where("route_stop.id = ?", stopID).
		First(&route)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &route, result.Error
}
