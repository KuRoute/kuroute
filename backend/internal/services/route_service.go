package services

import (
	"errors"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrForbiddenServiceAccess = errors.New("forbidden service access")
)

type RouteService struct {
	routeRepository        *repository.RouteRepository
	courierBatchRepository *repository.CourierBatchRepository
}

func NewRouteService(routeRepository *repository.RouteRepository, courierBatchRepository *repository.CourierBatchRepository) *RouteService {
	return &RouteService{
		routeRepository:        routeRepository,
		courierBatchRepository: courierBatchRepository,
	}
}

func (s *RouteService) CreateRoute(actor middleware.AuthService, req *domain.ComputeRouteRequest) (*domain.RouteResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if actor.Service != domain.ServiceRoute {
		return nil, ErrForbiddenServiceAccess
	}

	if _, err := s.courierBatchRepository.FindByID(req.CourierBatchID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourierBatchNotFound
		}
		return nil, err
	}

	stops := make([]domain.RouteStop, 0, len(req.Stops))
	for _, item := range req.Stops {
		stops = append(stops, domain.RouteStop{
			PackageID: item.PackageID,
			StopOrder: item.StopOrder,
			Lat:       item.Lat,
			Lng:       item.Lng,
			Status:    domain.PackageStatusAssigned,
		})
	}

	route := &domain.Route{
		CourierBatchID:  req.CourierBatchID,
		AlgorithmUsed:   req.AlgorithmUsed,
		TotalDistanceKm: req.TotalDistanceKm,
		TotalStops:      int16(len(req.Stops)),
	}

	if err := s.routeRepository.CreateRouteWithStops(route, stops); err != nil {
		return nil, err
	}

	createdStops := make([]domain.RouteStopResponse, 0, len(stops))
	for _, stop := range stops {
		createdStops = append(createdStops, domain.RouteStopResponse{
			ID: stop.ID,
			RouteID: stop.RouteID,
			PackageID: stop.PackageID,
			StopOrder: stop.StopOrder,
			Lat:       stop.Lat,
			Lng:       stop.Lng,
			Status:    stop.Status,
			RecipientName: stop.Package.RecipientName,
			AddressText: stop.Package.AddressText,
		})
	}

	resp := domain.NewRouteResponse(route, createdStops)
	return &resp, nil
}

func (s *RouteService) GetRoute(actor middleware.AuthUser, id uuid.UUID) (*domain.RouteResponse, error) {
	route, err := s.routeRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	if actor.Role == domain.UserRoleAdmin || actor.Role == domain.UserRoleStaffSortir {
		stops := make([]domain.RouteStopResponse, 0, len(route.Stops))
		for _, stop := range route.Stops {
			stops = append(stops, domain.NewRouteStopResponse(&stop))
		}
		resp := domain.NewRouteResponse(route, stops)
		return &resp, nil
	}

	if actor.Role == domain.UserRoleKurir {
		if route.CourierBatch.BatchAssignment.CourierID == actor.UserID {
			stops := make([]domain.RouteStopResponse, 0, len(route.Stops))
			for _, stop := range route.Stops {
				stops = append(stops, domain.NewRouteStopResponse(&stop))
			}
			resp := domain.NewRouteResponse(route, stops)
			return &resp, nil
		}
	}

	return nil, ErrForbiddenHubAccess
}

func (s *RouteService) ListRoutesByHub(actor middleware.AuthUser, hubID uuid.UUID, batchDate *time.Time) ([]domain.RouteResponse, error) {
	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrForbiddenHubAccess
	}

	if actor.Role != domain.UserRoleAdmin && hubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	routes, err := s.routeRepository.FindByHubID(hubID, batchDate)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.RouteResponse, 0, len(routes))
	for _, item := range routes {
		stops := make([]domain.RouteStopResponse, 0, len(item.Stops))
		for _, stop := range item.Stops {
			stops = append(stops, domain.NewRouteStopResponse(&stop))
		}
		mapped := domain.NewRouteResponse(&item, stops)
		resp = append(resp, mapped)
	}

	return resp, nil
}
