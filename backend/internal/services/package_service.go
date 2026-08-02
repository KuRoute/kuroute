package services

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
)

type PackageService struct {
	packageRepository *repository.PackageRepository
}

var (
	ErrPackageNotFound           = errors.New("package not found")
	ErrForbiddenStatusTransition = errors.New("forbidden status transition")
)

func NewPackageService(packageRepository *repository.PackageRepository) *PackageService {
	return &PackageService{
		packageRepository: packageRepository,
	}
}

func (s *PackageService) ListPackage() ([]domain.PackageResponse, error) {
	packs, err := s.packageRepository.ListPackage()
	if err != nil {
		return nil, err
	}

	resp := make([]domain.PackageResponse, 0, len(packs))
	for _, p := range packs {
		resp = append(resp, domain.NewPackageResponse(&p))
	}

	return resp, nil
}

func (s *PackageService) GetPackageByID(id uuid.UUID) (*domain.PackageResponse, error) {
	pack, err := s.packageRepository.GetPackageByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}

	resp := domain.NewPackageResponse(pack)
	return &resp, nil
}

func (s *PackageService) GetPackageByHubID(hubID uuid.UUID) ([]domain.PackageResponse, error) {
	packs, err := s.packageRepository.GetPackageByHubID(hubID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.PackageResponse, 0, len(packs))
	for _, p := range packs {
		resp = append(resp, domain.NewPackageResponse(&p))
	}

	return resp, nil
}

func (s *PackageService) GetStatusPackage(id uuid.UUID) (*domain.PackageStatus, error) {
	pack, err := s.packageRepository.GetStatusPackage(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("package with status not found")
		}
		return nil, err
	}

	return &pack.Status, nil
}

func (s *PackageService) GetPackageByHubAndStatus(hubID uuid.UUID, status domain.PackageStatus) ([]domain.PackageResponse, error) {
	packs, err := s.packageRepository.GetPackageByHubIDAndStatus(hubID, status)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.PackageResponse, 0, len(packs))
	for _, p := range packs {
		resp = append(resp, domain.NewPackageResponse(&p))
	}

	return resp, nil
}

func (s *PackageService) canUpdatePackageStatus(current, next domain.PackageStatus, role domain.UserRole) bool {
	if role == domain.UserRoleAdmin {
		return true
	}

	switch current {
	case domain.PackageStatusReceived:
		// received -> sorted: staff_sortir via LockerScanService.ScanIn (internal call)
		// dipicu saat staff memindai paket masuk ke locker hasil clustering,
		// bukan otomatis dari hasil clustering Python - clustering hanya
		// membuat rencana penempatan (Locker), belum mengubah status package.
		return next == domain.PackageStatusSorted && role == domain.UserRoleStaffSortir

	case domain.PackageStatusSorted:
		// sorted -> assigned: kurir via LockerScanService.ScanOut (internal call)
		// dipicu saat kurir memindai paket keluar dari locker untuk mulai rute.
		return next == domain.PackageStatusAssigned && role == domain.UserRoleKurir

	case domain.PackageStatusAssigned:
		// assigned -> in_delivery: kurir tap "mulai antar"
		return next == domain.PackageStatusInDelivery && role == domain.UserRoleKurir

	case domain.PackageStatusInDelivery:
		// in_delivery -> delivered/failed: kurir submit delivery report
		return (next == domain.PackageStatusDelivered ||
			next == domain.PackageStatusFailed) &&
			role == domain.UserRoleKurir

	default:
		// semua terminal state (delivered, failed) - tidak ada transisi lanjutan
		return false
	}
}

func (s *PackageService) CreatePackage(req *domain.CreatePackageRequest) (*domain.PackageResponse, error) {
	if req == nil {
		return nil, errors.New("request can't be empty")
	}

	if req.HubID == [16]byte{} || req.TrackingCode == "" || req.RecipientName == "" || req.AddressText == "" || req.Lat == 0 || req.Lng == 0 {
		return nil, errors.New("all fields are required")
	}

	pack := &domain.Package{
		HubID:         req.HubID,
		TrackingCode:  req.TrackingCode,
		RecipientName: req.RecipientName,
		AddressText:   req.AddressText,
		Lat:           req.Lat,
		Lng:           req.Lng,
	}

	if err := s.packageRepository.CreatePackage(pack); err != nil {
		return nil, err
	}

	resp := domain.NewPackageResponse(pack)
	return &resp, nil
}

func (s *PackageService) UpdatePackageStatus(id uuid.UUID, req *domain.UpdatePackageStatusRequest, callerRole domain.UserRole) (*domain.PackageResponse, error) {
	if req == nil || req.Status == "" {
		return nil, errors.New("status can't be empty")
	}

	pack, err := s.packageRepository.GetPackageByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}

	if !s.canUpdatePackageStatus(pack.Status, req.Status, callerRole) {
		return nil, ErrForbiddenStatusTransition
	}

	pack.Status = req.Status
	if err := s.packageRepository.UpdatePackageStatus(pack); err != nil {
		return nil, err
	}

	resp := domain.NewPackageResponse(pack)
	return &resp, nil
}

func (s *PackageService) DeletePackage(id uuid.UUID) error {
	return s.packageRepository.DeletePackage(id)
}
