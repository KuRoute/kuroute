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

type CourierBatchService struct {
	courierBatchRepository *repository.CourierBatchRepository
}

var (
	ErrCourierBatchNotFound = errors.New("courier batch not found")
	ErrInvalidBatchStatus   = errors.New("invalid courier batch status")
)

func NewCourierBatchService(courierBatchRepository *repository.CourierBatchRepository) *CourierBatchService {
	return &CourierBatchService{courierBatchRepository: courierBatchRepository}
}

func (s *CourierBatchService) GetCourierBatch(actor middleware.AuthUser, id uuid.UUID) (*domain.CourierBatchResponse, error) {
	courierBatch, err := s.courierBatchRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourierBatchNotFound
		}
		return nil, err
	}

	if actor.Role == domain.UserRoleAdmin || actor.Role == domain.UserRoleStaffSortir {
		resp := domain.NewCourierBatchResponse(courierBatch)
		return &resp, nil
	}

	if actor.Role == domain.UserRoleKurir && courierBatch.BatchAssignment.CourierID == actor.UserID {
		resp := domain.NewCourierBatchResponse(courierBatch)
		return &resp, nil
	}

	return nil, ErrForbiddenHubAccess
}

func (s *CourierBatchService) ListCourierBatchesByHub(actor middleware.AuthUser, hubID uuid.UUID, status *domain.BatchStatus, batchDate *time.Time) ([]domain.CourierBatchResponse, error) {
	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrForbiddenHubAccess
	}

	if actor.Role != domain.UserRoleAdmin && hubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	batches, err := s.courierBatchRepository.FindByHubID(hubID, status, batchDate)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.CourierBatchResponse, 0, len(batches))
	for _, item := range batches {
		mapped := domain.NewCourierBatchResponse(&item)
		resp = append(resp, mapped)
	}

	return resp, nil
}

func (s *CourierBatchService) ListMyCourierBatches(actor middleware.AuthUser) ([]domain.CourierBatchResponse, error) {
	if actor.Role != domain.UserRoleKurir {
		return nil, ErrForbiddenHubAccess
	}

	batches, err := s.courierBatchRepository.FindByCourierID(actor.UserID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.CourierBatchResponse, 0, len(batches))
	for _, item := range batches {
		mapped := domain.NewCourierBatchResponse(&item)
		resp = append(resp, mapped)
	}

	return resp, nil
}

func (s *CourierBatchService) ListInternalCourierBatches(status domain.BatchStatus) ([]domain.InternalCourierBatchResponse, error) {
	switch status {
	case domain.BatchStatusPendingRoute, domain.BatchStatusRouteReady, domain.BatchStatusInProgress, domain.BatchStatusCompleted:
	default:
		return nil, ErrInvalidBatchStatus
	}

	batches, err := s.courierBatchRepository.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.InternalCourierBatchResponse, 0, len(batches))
	for _, batch := range batches {
		resp = append(resp, domain.NewInternalCourierBatchResponse(&batch))
	}
	return resp, nil
}
