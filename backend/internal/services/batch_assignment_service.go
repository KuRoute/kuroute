package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
)

type BatchAssignmentService struct {
	batchAssignmentRepository *repository.BatchAssignmentRepository
	lockerScanRepository      *repository.LockerScanRepository
	lockerRepository          *repository.LockerRepository
	packageRepository         *repository.PackageRepository
	userRepository            *repository.UserRepository
}

var (
	ErrBatchAssignmentNotFound = errors.New("batch assignment not found")
	ErrBatchAssignmentLocked   = errors.New("batch assignment is no longer editable")
	ErrBatchAssignmentConflict = errors.New("batch assignment conflict")
	ErrCourierAlreadyAssigned  = errors.New("courier already has an assignment for this batch sequence")
	ErrLockerAlreadyAssigned   = errors.New("locker is already assigned to a courier for the requested batch")
	ErrCourierNotFound         = errors.New("courier not found")
	ErrRouteNotFound           = errors.New("route not found")
	ErrBatchAlreadyStarted     = errors.New("batch has already started")
	ErrNoPackagesToStart       = errors.New("no assigned packages are ready to start")
	ErrLockerEmpty             = errors.New("locker has no packages ready for assignment")
)

func NewBatchAssignmentService(batchAssignmentRepository *repository.BatchAssignmentRepository, lockerScanRepository *repository.LockerScanRepository, lockerRepository *repository.LockerRepository, packageRepository *repository.PackageRepository, userRepository *repository.UserRepository) *BatchAssignmentService {
	return &BatchAssignmentService{
		batchAssignmentRepository: batchAssignmentRepository,
		lockerScanRepository:      lockerScanRepository,
		lockerRepository:          lockerRepository,
		packageRepository:         packageRepository,
		userRepository:            userRepository,
	}
}

func (s *BatchAssignmentService) CreateAssignment(actor middleware.AuthUser, req *domain.AssignLockerToCourierRequest) (*domain.BatchAssignmentResponse, error) {
	if req == nil {
		return nil, errors.New("request can't be empty")
	}

	if actor.Role != domain.UserRoleStaffSortir && actor.Role != domain.UserRoleAdmin {
		return nil, ErrForbiddenHubAccess
	}

	locker, err := s.lockerRepository.FindByID(req.LockerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	courier, err := s.userRepository.GetUserByID(req.CourierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}
	if courier.Role != domain.UserRoleKurir {
		return nil, ErrCourierNotFound
	}

	if actor.Role != domain.UserRoleAdmin && courier.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	hasActive, err := s.lockerRepository.HasActivePackages(locker.ID)
	if err != nil {
		return nil, err
	}
	if !hasActive {
		return nil, ErrLockerEmpty
	}

	batchDate := time.Now()
	if req.BatchDate != nil {
		batchDate = req.BatchDate.UTC().Truncate(24 * time.Hour)
	}

	if existing, err := s.batchAssignmentRepository.FindByLockerAndDateSequence(req.LockerID, batchDate, req.BatchSequence); err == nil && existing != nil {
		return nil, ErrLockerAlreadyAssigned
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existing, err := s.batchAssignmentRepository.FindByCourierAndDateSequence(req.CourierID, batchDate, req.BatchSequence); err == nil && existing != nil {
		return nil, ErrCourierAlreadyAssigned
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	assignment := &domain.BatchAssignment{
		LockerID:      req.LockerID,
		CourierID:     req.CourierID,
		AssignedByID:  actor.UserID,
		BatchDate:     batchDate,
		BatchSequence: req.BatchSequence,
		AssignedAt:    time.Now(),
	}

	if err := s.batchAssignmentRepository.Create(assignment); err != nil {
		return nil, err
	}

	courierBatch := &domain.CourierBatch{
		BatchAssignmentID: assignment.ID,
		Status:            domain.BatchStatusPendingRoute,
	}
	if err := s.batchAssignmentRepository.CreateCourierBatch(courierBatch); err != nil {
		return nil, err
	}

	resp := domain.NewBatchAssignmentResponse(assignment)
	return &resp, nil
}

func (s *BatchAssignmentService) GetAssignment(actor middleware.AuthUser, id uuid.UUID) (*domain.BatchAssignmentResponse, error) {
	assignment, err := s.batchAssignmentRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBatchAssignmentNotFound
		}
		return nil, err
	}

	if actor.Role == domain.UserRoleAdmin || actor.Role == domain.UserRoleStaffSortir {
		resp := domain.NewBatchAssignmentResponse(assignment)
		return &resp, nil
	}

	if actor.Role == domain.UserRoleKurir && assignment.CourierID == actor.UserID {
		resp := domain.NewBatchAssignmentResponse(assignment)
		return &resp, nil
	}

	return nil, ErrForbiddenHubAccess
}

func (s *BatchAssignmentService) ListAssignmentsByHub(actor middleware.AuthUser, hubID uuid.UUID, batchDate *time.Time, courierID *uuid.UUID) ([]domain.BatchAssignmentResponse, error) {
	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrForbiddenHubAccess
	}

	if actor.Role != domain.UserRoleAdmin && hubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	assignments, err := s.batchAssignmentRepository.FindByHubID(hubID, batchDate, courierID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.BatchAssignmentResponse, 0, len(assignments))
	for _, item := range assignments {
		mapped := domain.NewBatchAssignmentResponse(&item)
		resp = append(resp, mapped)
	}

	return resp, nil
}

func (s *BatchAssignmentService) ListMyAssignments(actor middleware.AuthUser, batchDate *time.Time) ([]domain.BatchAssignmentResponse, error) {
	if actor.Role != domain.UserRoleKurir {
		return nil, ErrForbiddenHubAccess
	}

	var date time.Time
	if batchDate != nil {
		date = batchDate.UTC().Truncate(24 * time.Hour)
	} else {
		date = time.Now().UTC().Truncate(24 * time.Hour)
	}

	assignments, err := s.batchAssignmentRepository.FindByCourierAndDate(actor.UserID, date)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.BatchAssignmentResponse, 0, len(assignments))
	for _, item := range assignments {
		mapped := domain.NewBatchAssignmentResponse(&item)
		resp = append(resp, mapped)
	}
	return resp, nil
}

func (s *BatchAssignmentService) UpdateAssignment(actor middleware.AuthUser, id uuid.UUID, req *domain.UpdateBatchAssignmentRequest) (*domain.BatchAssignmentResponse, error) {
	if req == nil {
		return nil, errors.New("request can't be empty")
	}
	if actor.Role != domain.UserRoleStaffSortir && actor.Role != domain.UserRoleAdmin {
		return nil, ErrForbiddenHubAccess
	}

	assignment, err := s.batchAssignmentRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBatchAssignmentNotFound
		}
		return nil, err
	}

	courierBatch, err := s.batchAssignmentRepository.FindCourierBatchByAssignmentID(assignment.ID)
	if err == nil && courierBatch != nil && (courierBatch.Status == domain.BatchStatusInProgress || courierBatch.Status == domain.BatchStatusCompleted) {
		return nil, ErrBatchAssignmentLocked
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if req.CourierID != nil {
		courier, err := s.userRepository.GetUserByID(*req.CourierID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCourierNotFound
			}
			return nil, err
		}
		if courier.Role != domain.UserRoleKurir {
			return nil, ErrCourierNotFound
		}

		batchDate := assignment.BatchDate
		if req.BatchDate != nil {
			batchDate = req.BatchDate.UTC().Truncate(24 * time.Hour)
		}
		sequence := assignment.BatchSequence
		if req.BatchSequence != nil {
			sequence = *req.BatchSequence
		}

		if existing, err := s.batchAssignmentRepository.FindByCourierAndDateSequence(*req.CourierID, batchDate, sequence); err == nil && existing != nil && existing.ID != assignment.ID {
			return nil, ErrCourierAlreadyAssigned
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		assignment.CourierID = *req.CourierID
	}

	if req.BatchSequence != nil {
		assignment.BatchSequence = *req.BatchSequence
	}
	if req.BatchDate != nil {
		assignment.BatchDate = req.BatchDate.UTC().Truncate(24 * time.Hour)
	}

	if req.BatchSequence != nil || req.BatchDate != nil {
		batchDate := assignment.BatchDate
		sequence := assignment.BatchSequence
		if existing, err := s.batchAssignmentRepository.FindByLockerAndDateSequence(assignment.LockerID, batchDate, sequence); err == nil && existing != nil && existing.ID != assignment.ID {
			return nil, ErrLockerAlreadyAssigned
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if err := s.batchAssignmentRepository.Update(assignment); err != nil {
		return nil, err
	}

	resp := domain.NewBatchAssignmentResponse(assignment)
	return &resp, nil
}

func (s *BatchAssignmentService) DeleteAssignment(actor middleware.AuthUser, id uuid.UUID) error {
	if actor.Role != domain.UserRoleStaffSortir && actor.Role != domain.UserRoleAdmin {
		return ErrForbiddenHubAccess
	}

	assignment, err := s.batchAssignmentRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBatchAssignmentNotFound
		}
		return err
	}

	courierBatch, err := s.batchAssignmentRepository.FindCourierBatchByAssignmentID(assignment.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if courierBatch != nil && courierBatch.Status != domain.BatchStatusPendingRoute {
		return ErrBatchAssignmentLocked
	}

	return s.batchAssignmentRepository.Delete(id)
}

func (s *BatchAssignmentService) StartBatch(actor middleware.AuthUser, id uuid.UUID) (*domain.CourierBatchResponse, error) {
	if actor.Role != domain.UserRoleKurir {
		return nil, ErrForbiddenHubAccess
	}

	assignment, err := s.batchAssignmentRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBatchAssignmentNotFound
		}
		return nil, err
	}
	if assignment.CourierID != actor.UserID {
		return nil, ErrForbiddenHubAccess
	}

	courierBatch, err := s.batchAssignmentRepository.FindCourierBatchByAssignmentID(assignment.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBatchAssignmentNotFound
		}
		return nil, err
	}
	if courierBatch.Status == domain.BatchStatusInProgress || courierBatch.Status == domain.BatchStatusCompleted {
		return nil, ErrBatchAlreadyStarted
	}
	if courierBatch.Status != domain.BatchStatusPendingRoute && courierBatch.Status != domain.BatchStatusRouteReady {
		return nil, ErrBatchAssignmentLocked
	}

	now := time.Now()
	updated, err := s.batchAssignmentRepository.StartCourierBatch(courierBatch.ID, assignment.LockerID, now)
	if err != nil {
		if errors.Is(err, repository.ErrBatchStatusTransition) {
			return nil, ErrBatchAssignmentLocked
		}
		return nil, err
	}
	if updated == 0 {
		return nil, ErrNoPackagesToStart
	}
	courierBatch.Status = domain.BatchStatusInProgress
	courierBatch.StartedAt = &now

	resp := domain.NewCourierBatchResponse(courierBatch)
	return &resp, nil
}

func (s *BatchAssignmentService) GetRouteForAssignment(actor middleware.AuthUser, id uuid.UUID) (*domain.RouteResponse, error) {
	assignment, err := s.batchAssignmentRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBatchAssignmentNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir && assignment.CourierID != actor.UserID {
		return nil, ErrForbiddenHubAccess
	}

	route, err := s.batchAssignmentRepository.FindRouteByAssignmentID(assignment.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	stops, err := s.batchAssignmentRepository.FindRouteStopsByRouteID(route.ID)
	if err != nil {
		return nil, err
	}

	stopResponses := make([]domain.RouteStopResponse, 0, len(stops))
	for _, stop := range stops {
		stopResponses = append(stopResponses, domain.NewRouteStopResponse(&stop))
	}

	resp := domain.NewRouteResponse(route, stopResponses)
	return &resp, nil
}
