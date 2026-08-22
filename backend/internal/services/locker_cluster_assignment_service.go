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
	ErrLockerClusterAssignmentNotFound  = errors.New("locker cluster assignment not found")
	ErrLockerClusterAssignmentForbidden = errors.New("not authorized to access locker cluster assignments")
)

type LockerClusterAssignmentService struct {
	repository *repository.LockerClusterAssignmentRepository
}

func NewLockerClusterAssignmentService(repository *repository.LockerClusterAssignmentRepository) *LockerClusterAssignmentService {
	return &LockerClusterAssignmentService{repository: repository}
}

func (s *LockerClusterAssignmentService) GetByPackageID(actor middleware.AuthUser, packageID uuid.UUID) (*domain.LockerClusterAssignmentResponse, error) {
	if actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrLockerClusterAssignmentForbidden
	}

	assignment, err := s.repository.FindByPackageID(packageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerClusterAssignmentNotFound
		}
		return nil, err
	}

	resp := domain.NewLockerClusterAssignmentResponse(assignment)
	return &resp, nil
}

func (s *LockerClusterAssignmentService) ListByLockerID(actor middleware.AuthUser, lockerID uuid.UUID) ([]domain.LockerClusterAssignmentResponse, error) {
	if actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrLockerClusterAssignmentForbidden
	}

	assignments, err := s.repository.FindByLockerID(lockerID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.LockerClusterAssignmentResponse, 0, len(assignments))
	for _, assignment := range assignments {
		resp = append(resp, domain.NewLockerClusterAssignmentResponse(&assignment))
	}
	return resp, nil
}
