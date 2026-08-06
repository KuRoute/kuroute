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
	ErrCourierProfileNotFound      = errors.New("courier profile not found")
	ErrCourierProfileAlreadyExists = errors.New("courier profile already exists")
	ErrForbidden                   = errors.New("forbidden")
)

type CourierProfileService struct {
	courierProfileRepository *repository.CourierProfileRepository
	userRepository           *repository.UserRepository
}

func NewCourierProfileService(courierProfileRepository *repository.CourierProfileRepository, userRepository *repository.UserRepository) *CourierProfileService {
	return &CourierProfileService{
		courierProfileRepository: courierProfileRepository,
		userRepository:           userRepository,
	}
}

func (s *CourierProfileService) GetCourierProfile(actor *middleware.AuthUser, courierID uuid.UUID) (*domain.CourierProfileResponse, error) {
	courier, err := s.userRepository.GetUserByID(courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCourierProfileNotFound
		}
		return nil, err
	}

	if courier.Role != domain.UserRoleKurir {
		return nil, ErrCourierProfileNotFound
	}

	if actor.Role == domain.UserRoleKurir && actor.UserID != courierID {
		return nil, ErrForbidden
	}

	if actor.Role == domain.UserRoleStaffSortir && courier.HubID != actor.HubID {
		return nil, ErrForbidden
	}

	profile, err := s.courierProfileRepository.GetCourierProfileByUserID(courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCourierProfileNotFound
		}
		return nil, err
	}

	resp := domain.NewCourierProfileResponse(profile)
	return &resp, nil
}

func (s *CourierProfileService) CreateCourierProfile(courierID uuid.UUID, req *domain.CreateCourierProfileRequest) (*domain.CourierProfileResponse, error) {
	courier, err := s.userRepository.GetUserByID(courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCourierProfileNotFound
		}
		return nil, err
	}

	if courier.Role != domain.UserRoleKurir {
		return nil, ErrCourierProfileNotFound
	}

	courierProfile, err := s.courierProfileRepository.GetCourierProfileByUserID(courierID)
	if courierProfile != nil {
		return nil, ErrCourierProfileAlreadyExists
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	profile := &domain.CourierProfile{
		UserID:       courierID,
		VehicleType:  req.VehicleType,
		VehiclePlate: req.VehiclePlate,
	}

	if err := s.courierProfileRepository.CreateCourierProfile(profile); err != nil {
		return nil, err
	}

	resp := domain.NewCourierProfileResponse(profile)
	return &resp, nil
}

func (s *CourierProfileService) UpdateCourierProfile(courierID uuid.UUID, req *domain.UpdateCourierProfileRequest) (*domain.CourierProfileResponse, error) {
	courier, err := s.userRepository.GetUserByID(courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCourierProfileNotFound
		}
		return nil, err
	}

	if courier.Role != domain.UserRoleKurir {
		return nil, ErrCourierProfileNotFound
	}

	profile, err := s.courierProfileRepository.GetCourierProfileByUserID(courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			profile = &domain.CourierProfile{UserID: courierID}
		} else {
			return nil, err
		}
	}

	if req.VehicleType != nil {
		profile.VehicleType = *req.VehicleType
	}
	if req.VehiclePlate != nil {
		profile.VehiclePlate = *req.VehiclePlate
	}

	if profile.UserID == uuid.Nil {
		profile.UserID = courierID
	}

	if err := s.saveCourierProfile(profile); err != nil {
		return nil, err
	}

	resp := domain.NewCourierProfileResponse(profile)
	return &resp, nil
}

func (s *CourierProfileService) saveCourierProfile(profile *domain.CourierProfile) error {
	_, err := s.courierProfileRepository.GetCourierProfileByUserID(profile.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.courierProfileRepository.CreateCourierProfile(profile)
		}
		return err
	}

	return s.courierProfileRepository.UpdateCourierProfile(profile)
}

func (s *CourierProfileService) ListCouriersByHub(actor *middleware.AuthUser, hubID uuid.UUID, vehicleType *domain.VehicleType) ([]domain.CourierProfileResponse, error) {
	if actor.Role == domain.UserRoleStaffSortir && actor.HubID != hubID {
		return nil, ErrForbidden
	}

	profiles, err := s.courierProfileRepository.FindCouriersByHubID(hubID, vehicleType)
	if err != nil {
		return nil, err
	}

	results := make([]domain.CourierProfileResponse, 0, len(profiles))
	for _, profile := range profiles {
		results = append(results, domain.NewCourierProfileResponse(&profile))
	}
	return results, nil
}
