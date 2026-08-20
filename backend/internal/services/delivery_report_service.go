package services

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDeliveryReportNotFound   = errors.New("delivery report not found")
	ErrDeliveryReportForbidden  = errors.New("not authorized to access this delivery report")
	ErrRouteStopAccessForbidden = errors.New("not authorized to access this route stop")
	ErrBatchNotFound            = errors.New("courier batch not found")
	ErrBatchNotInProgress       = errors.New("courier batch must be in_progress before submitting a delivery report")
	ErrDeliveryStatusTransition = errors.New("package and route stop must be in_delivery")
	ErrUploadInvalidFileType    = errors.New("file type must be jpg or png")
	ErrUploadTooLarge           = errors.New("file size exceeds 5MB limit")
	ErrDeliverySubmitInvalid    = errors.New("invalid delivery report submission")
)

type DeliveryReportService struct {
	deliveryReportRepository *repository.DeliveryReportRepository
	routeStopRepository      *repository.RouteStopRepository
	routeRepository          *repository.RouteRepository
	courierBatchRepository   *repository.CourierBatchRepository
	packageRepository        *repository.PackageRepository
}

func NewDeliveryReportService(
	deliveryReportRepository *repository.DeliveryReportRepository,
	routeStopRepository *repository.RouteStopRepository,
	routeRepository *repository.RouteRepository,
	courierBatchRepository *repository.CourierBatchRepository,
	packageRepository *repository.PackageRepository,
) *DeliveryReportService {
	return &DeliveryReportService{
		deliveryReportRepository: deliveryReportRepository,
		routeStopRepository:      routeStopRepository,
		routeRepository:          routeRepository,
		courierBatchRepository:   courierBatchRepository,
		packageRepository:        packageRepository,
	}
}

func (s *DeliveryReportService) validateUpload(file multipart.File, header *multipart.FileHeader) error {
	if file == nil || header == nil {
		return ErrDeliverySubmitInvalid
	}
	if header.Size > 5*1024*1024 {
		return ErrUploadTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return ErrUploadInvalidFileType
	}
	return nil
}

func (s *DeliveryReportService) UploadProof(actor middleware.AuthUser, file multipart.File, header *multipart.FileHeader) (string, error) {
	if actor.Role != domain.UserRoleKurir {
		return "", ErrDeliveryReportForbidden
	}
	if err := s.validateUpload(file, header); err != nil {
		return "", err
	}

	storageDir := filepath.Join("uploads", "delivery-proof")
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		return "", err
	}

	filename := uuid.NewString() + strings.ToLower(filepath.Ext(header.Filename))
	storedPath := filepath.Join(storageDir, filename)
	storedFile, err := os.Create(storedPath)
	if err != nil {
		return "", err
	}
	defer storedFile.Close()

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err := io.Copy(storedFile, file); err != nil {
		return "", err
	}

	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.ReplaceAll(filepath.ToSlash(storedPath), "\\", "/")), nil
}

func (s *DeliveryReportService) SubmitDeliveryReport(actor middleware.AuthUser, stopID uuid.UUID, req *domain.SubmitDeliveryReportRequest) (*domain.DeliveryReportResponse, error) {
	if req == nil {
		return nil, ErrDeliverySubmitInvalid
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if actor.Role != domain.UserRoleKurir {
		return nil, ErrDeliveryReportForbidden
	}

	stop, err := s.routeStopRepository.FindByID(stopID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteStopNotFound
		}
		return nil, err
	}

	route, err := s.routeRepository.FindByID(stop.RouteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}
	if route.CourierBatch.BatchAssignment.CourierID != actor.UserID {
		return nil, ErrRouteStopAccessForbidden
	}

	batch, err := s.courierBatchRepository.FindByID(route.CourierBatchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBatchNotFound
		}
		return nil, err
	}
	if batch.Status != domain.BatchStatusInProgress {
		return nil, ErrBatchNotInProgress
	}
	if stop.Status != domain.PackageStatusInDelivery || stop.Package.Status != domain.PackageStatusInDelivery {
		return nil, ErrDeliveryStatusTransition
	}

	report := &domain.DeliveryReport{
		RouteStopID:   stop.ID,
		Result:        req.Result,
		FailureReason: req.FailureReason,
		PhotoURL:      req.PhotoURL,
		SignatureURL:  req.SignatureURL,
		Notes:         req.Notes,
		ReportedAt:    time.Now(),
	}

	newStopStatus := domain.PackageStatusDelivered
	if req.Result.IsFailure() {
		newStopStatus = domain.PackageStatusFailed
	}

	if err := s.deliveryReportRepository.CreateAndUpdateStatus(report, stop.ID, stop.PackageID, route.ID, route.CourierBatchID, newStopStatus); err != nil {
		if errors.Is(err, repository.ErrDeliveryStatusTransition) {
			return nil, ErrDeliveryStatusTransition
		}
		return nil, err
	}

	resp := domain.NewDeliveryReportResponse(report)
	return &resp, nil
}

func (s *DeliveryReportService) GetDeliveryReport(actor middleware.AuthUser, id uuid.UUID) (*domain.DeliveryReportResponse, error) {
	report, err := s.deliveryReportRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDeliveryReportNotFound
		}
		return nil, err
	}

	stop, err := s.routeStopRepository.FindByID(report.RouteStopID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRouteStopNotFound
		}
		return nil, err
	}

	route, err := s.routeRepository.FindByID(stop.RouteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir && route.CourierBatch.BatchAssignment.CourierID != actor.UserID {
		return nil, ErrDeliveryReportForbidden
	}

	resp := domain.NewDeliveryReportResponse(report)
	return &resp, nil
}

func (s *DeliveryReportService) ListDeliveryReportsByRoute(actor middleware.AuthUser, routeID uuid.UUID) ([]domain.DeliveryReportResponse, error) {
	route, err := s.routeRepository.FindByID(routeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir && route.CourierBatch.BatchAssignment.CourierID != actor.UserID {
		return nil, ErrDeliveryReportForbidden
	}

	reports, err := s.deliveryReportRepository.FindByRouteID(routeID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.DeliveryReportResponse, 0, len(reports))
	for _, report := range reports {
		resp = append(resp, domain.NewDeliveryReportResponse(&report))
	}
	return resp, nil
}

func (s *DeliveryReportService) ListDeliveryReportsByPackage(actor middleware.AuthUser, packageID uuid.UUID) ([]domain.DeliveryReportResponse, error) {
	if actor.Role != domain.UserRoleAdmin && actor.Role != domain.UserRoleStaffSortir {
		return nil, ErrDeliveryReportForbidden
	}

	if _, err := s.packageRepository.GetPackageByID(packageID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}

	reports, err := s.deliveryReportRepository.FindByPackageID(packageID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.DeliveryReportResponse, 0, len(reports))
	for _, report := range reports {
		resp = append(resp, domain.NewDeliveryReportResponse(&report))
	}
	return resp, nil
}
