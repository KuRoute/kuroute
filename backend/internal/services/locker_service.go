package services

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
)

type LockerService struct {
	lockerRepository *repository.LockerRepository
	hubRepository    *repository.HubRepository
}

var (
	ErrLockerNotFound     = errors.New("locker not found")
	ErrLockerLabelTaken   = errors.New("locker label already used in this hub")
	ErrLockerNotEmpty     = errors.New("locker still has active packages")
	ErrHubNotFound        = errors.New("hub not found")
	ErrForbiddenHubAccess = errors.New("not authorized to access this hub's data")
)

func NewLockerService(lockerRepository *repository.LockerRepository, hubRepository *repository.HubRepository) *LockerService {
	return &LockerService{
		lockerRepository: lockerRepository,
		hubRepository:    hubRepository,
	}
}

func (s *LockerService) CreateLocker(req *domain.CreateLockerRequest) (*domain.LockerResponse, error) {
	if _, err := s.hubRepository.GetHubByID(req.HubID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHubNotFound
		}
		return nil, err
	}

	if existing, err := s.lockerRepository.FindByHubIDAndLabel(req.HubID, req.Label); err == nil && existing != nil {
		return nil, ErrLockerLabelTaken
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	locker := &domain.Locker{
		HubID:       req.HubID,
		Label:       req.Label,
		ClusterArea: req.ClusterArea,
	}

	if err := s.lockerRepository.Create(locker); err != nil {
		return nil, err
	}

	resp := domain.NewLockerResponse(locker)
	return &resp, nil
}

func (s *LockerService) GetLocker(actor middleware.AuthUser, id uuid.UUID) (*domain.LockerResponse, error) {
	locker, err := s.lockerRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	resp := domain.NewLockerResponse(locker)
	return &resp, nil
}

func (s *LockerService) ListLockersByHub(actor middleware.AuthUser, hubID uuid.UUID) ([]domain.LockerResponse, error) {
	if actor.Role != domain.UserRoleAdmin && hubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	lockers, err := s.lockerRepository.FindByHubID(hubID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.LockerResponse, 0, len(lockers))
	for _, locker := range lockers {
		resp = append(resp, domain.NewLockerResponse(&locker))
	}

	return resp, nil
}

func (s *LockerService) UpdateLocker(actor middleware.AuthUser, id uuid.UUID, req domain.UpdateLockerRequest) (*domain.LockerResponse, error) {
	if actor.Role != domain.UserRoleAdmin {
		return nil, ErrForbiddenHubAccess
	}

	locker, err := s.lockerRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if req.Label != nil && *req.Label != locker.Label {
		if existing, err := s.lockerRepository.FindByHubIDAndLabel(locker.HubID, *req.Label); err == nil && existing != nil && existing.ID != locker.ID {
			return nil, ErrLockerLabelTaken
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		locker.Label = *req.Label
	}

	if req.ClusterArea != nil {
		locker.ClusterArea = *req.ClusterArea
	}

	if err := s.lockerRepository.Update(locker); err != nil {
		return nil, err
	}

	resp := domain.NewLockerResponse(locker)
	return &resp, nil
}

func (s *LockerService) DeleteLocker(actor middleware.AuthUser, id uuid.UUID) error {
	if actor.Role != domain.UserRoleAdmin {
		return ErrForbiddenHubAccess
	}

	locker, err := s.lockerRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLockerNotFound
		}
		return err
	}

	hasActive, err := s.lockerRepository.HasActivePackages(locker.ID)
	if err != nil {
		return err
	}
	if hasActive {
		return ErrLockerNotEmpty
	}

	return s.lockerRepository.Delete(id)
}
