package services

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrRouteStopNotFound = errors.New("route stop not found")
)

type RouteStopService struct {
	routeStopRepository *repository.RouteStopRepository
	routeRepository     *repository.RouteRepository
}

func NewRouteStopService(routeStopRepository *repository.RouteStopRepository, routeRepository *repository.RouteRepository) *RouteStopService {
	return &RouteStopService{
		routeStopRepository: routeStopRepository,
		routeRepository:     routeRepository,
	}
}

func canAccessRouteStop(actor middleware.AuthUser, ownerCourierID *uuid.UUID) bool {
	if actor.Role == domain.UserRoleAdmin || actor.Role == domain.UserRoleStaffSortir {
		return true
	}
	if actor.Role != domain.UserRoleKurir || ownerCourierID == nil {
		return false
	}
	return *ownerCourierID == actor.UserID
}

func (s *RouteStopService) GetRouteStop(actor middleware.AuthUser, id uuid.UUID) (*domain.RouteStopDetailResponse, error) {
	stop, err := s.routeStopRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteStopNotFound
		}
		return nil, err
	}

	route, err := s.routeRepository.FindByID(stop.RouteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	ownerCourierID := uuid.Nil
	if route.CourierBatch.BatchAssignment.ID != (uuid.UUID{}) {
		ownerCourierID = route.CourierBatch.BatchAssignment.CourierID
	}
	if !canAccessRouteStop(actor, &ownerCourierID) {
		return nil, ErrForbiddenHubAccess
	}

	resp := domain.NewRouteStopDetailResponse(stop)
	return &resp, nil
}

func (s *RouteStopService) ListRouteStops(actor middleware.AuthUser, routeID uuid.UUID) ([]domain.RouteStopResponse, error) {
	route, err := s.routeRepository.FindByID(routeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	ownerCourierID := uuid.Nil
	if route.CourierBatch.BatchAssignment.ID != (uuid.UUID{}) {
		ownerCourierID = route.CourierBatch.BatchAssignment.CourierID
	}
	if !canAccessRouteStop(actor, &ownerCourierID) {
		return nil, ErrForbiddenHubAccess
	}

	stops, err := s.routeStopRepository.FindByRouteID(routeID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.RouteStopResponse, 0, len(stops))
	for _, stop := range stops {
		resp = append(resp, domain.NewRouteStopResponse(&stop))
	}
	return resp, nil
}
