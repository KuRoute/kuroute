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

type LockerScanService struct {
	lockerScanRepository *repository.LockerScanRepository
	lockerRepository     *repository.LockerRepository
	packageRepository    *repository.PackageRepository
}

var (
	ErrLockerScanNotFound          = errors.New("locker scan not found")
	ErrLockerScanAlreadyCheckedIn  = errors.New("package is already scanned into a locker")
	ErrLockerScanAlreadyCheckedOut = errors.New("package is already scanned out from locker")
	ErrLockerScanHistoryForbidden  = errors.New("not authorized to view package locker history")
)

func NewLockerScanService( lockerScanRepository *repository.LockerScanRepository, lockerRepository *repository.LockerRepository, packageRepository *repository.PackageRepository) *LockerScanService {
	return &LockerScanService{
		lockerScanRepository: lockerScanRepository,
		lockerRepository:     lockerRepository,
		packageRepository:    packageRepository,
	}
}

var timeNow = time.Now()
var timeNowPtr = &timeNow

func (s *LockerScanService) ScanIn(actor middleware.AuthUser, req *domain.ScanPackageToLockerRequest) (*domain.LockerScanResponse, error) {
	if req == nil {
		return nil, errors.New("request can't be empty")
	}

	if  actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrForbiddenHubAccess
	}

	locker, err := s.lockerRepository.FindByID(req.LockerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleStaffSortir && locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	pkg, err := s.packageRepository.GetPackageByID(req.PackageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleStaffSortir && pkg.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	if _, err := s.lockerScanRepository.FindActiveByPackageID(req.PackageID); err == nil {
		return nil, ErrLockerScanAlreadyCheckedIn
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if pkg.Status == domain.PackageStatusReceived {
		pkg.Status = domain.PackageStatusSorted
		if err := s.packageRepository.UpdatePackageStatus(pkg); err != nil {
			return nil, err
		}
	}

	scan := &domain.LockerScan{
		LockerID:    req.LockerID,
		PackageID:   req.PackageID,
		ScannedByID: actor.UserID,
		ScannedAt:   time.Now(),
	}
	if err := s.lockerScanRepository.Create(scan); err != nil {
		return nil, err
	}

	resp := domain.NewLockerScanResponse(scan)
	return &resp, nil
}

func (s *LockerScanService) ScanOut(actor middleware.AuthUser, req *domain.ScanPackageToLockerRequest) (*domain.LockerScanResponse, error) {
	if req == nil {
		return nil, errors.New("request can't be empty")
	}

	if actor.Role != domain.UserRoleKurir {
		return nil, ErrForbiddenHubAccess
	}

	scan, err := s.lockerScanRepository.FindActiveByPackageID(req.PackageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerScanNotFound
		}
		return nil, err
	}

	locker, err := s.lockerRepository.FindByID(scan.LockerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleKurir && locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	if scan.ScannedOutAt != nil {
		return nil, ErrLockerScanAlreadyCheckedOut
	}

	scan.ScannedOutAt = timeNowPtr
	scan.ScannedOutByID = &actor.UserID
	if err := s.lockerScanRepository.Update(scan); err != nil {
		return nil, err
	}

	resp := domain.NewLockerScanResponse(scan)
	return &resp, nil
}

func (s *LockerScanService) GetLockerScan(actor middleware.AuthUser, id uuid.UUID) (*domain.LockerScanResponse, error) {
	scan, err := s.lockerScanRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerScanNotFound
		}
		return nil, err
	}

	locker, err := s.lockerRepository.FindByID(scan.LockerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	resp := domain.NewLockerScanResponse(scan)
	return &resp, nil
}

func (s *LockerScanService) ListLockerScansByLocker(actor middleware.AuthUser, lockerID uuid.UUID) ([]domain.LockerScanResponse, error) {
	locker, err := s.lockerRepository.FindByID(lockerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLockerNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir && actor.Role != domain.UserRoleKurir && locker.HubID != actor.HubID {
		return nil, ErrForbiddenHubAccess
	}

	scans, err := s.lockerScanRepository.FindByLockerID(lockerID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.LockerScanResponse, 0, len(scans))
	for _, scan := range scans {
		item := domain.NewLockerScanResponse(&scan)
		resp = append(resp, item)
	}
	return resp, nil
}

func (s *LockerScanService) ListLockerScansByPackage(actor middleware.AuthUser, packageID uuid.UUID) ([]domain.LockerScanResponse, error) {
	if actor.Role != domain.UserRoleAdmin {
		return nil, ErrLockerScanHistoryForbidden
	}

	scans, err := s.lockerScanRepository.FindByPackageID(packageID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.LockerScanResponse, 0, len(scans))
	for _, scan := range scans {
		item := domain.NewLockerScanResponse(&scan)
		resp = append(resp, item)
	}
	return resp, nil
}

func (s *LockerScanService) ListLockerScansByCourierToday(actor middleware.AuthUser) ([]domain.LockerScanResponse, error) {
	if actor.Role != domain.UserRoleKurir {
		return nil, ErrForbiddenHubAccess
	}

	scans, err := s.lockerScanRepository.FindByCourierToday(actor.UserID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.LockerScanResponse, 0, len(scans))
	for _, scan := range scans {
		item := domain.NewLockerScanResponse(&scan)
		resp = append(resp, item)
	}
	return resp, nil
}
