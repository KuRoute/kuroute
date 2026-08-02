package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PackageHandler struct {
	packageService *services.PackageService
}

func NewPackageHandler(packageService *services.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

func (h *PackageHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	packs, err := h.packageService.ListPackage()
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, packs)
}

func (h *PackageHandler) GetPackageByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	packResp, err := h.packageService.GetPackageByID(id)
	if err != nil {
		response.Fail(w, http.StatusNotFound, "PACKAGE_NOT_FOUND", err.Error())
		return
	}

	response.OK(w, http.StatusOK, packResp)
}

func (h *PackageHandler) GetPackagesByHubID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	hubIDStr := vars["hubId"]
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub ID")
		return
	}

	ctxValue := r.Context().Value(middleware.UserContextKey)
	if ctxValue == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	authUser, ok := ctxValue.(*middleware.AuthUser)
	if !ok || authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authentication context")
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		packs, err := h.packageService.GetPackageByHubID(hubID)
		if err != nil {
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
			return
		}

		response.OK(w, http.StatusOK, packs)
		return
	}

	status := domain.PackageStatus(statusStr)
	switch status {
	case	domain.PackageStatusReceived,
			domain.PackageStatusSorted,
			domain.PackageStatusAssigned,
			domain.PackageStatusInDelivery,
			domain.PackageStatusDelivered,
			domain.PackageStatusFailed:
	// valid
	default:
		response.Fail(w, http.StatusBadRequest, "INVALID_STATUS", "Invalid package status")
		return
	}

	packs, err := h.packageService.GetPackageByHubAndStatus(hubID, status)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, packs)
}

func (h *PackageHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctxValue := r.Context().Value(middleware.UserContextKey)
	if ctxValue == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	authUser, ok := ctxValue.(*middleware.AuthUser)
	if !ok || authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authentication context")
		return
	}

	if authUser.Role != domain.UserRoleAdmin {
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to create a new package")
		return
	}

	var reqs []domain.CreatePackageRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	if len(reqs) == 0 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Request body cannot be empty")
		return
	}

	var packages []domain.PackageResponse

	for _, req := range reqs {
		packageResp, err := h.packageService.CreatePackage(&req)
		if err != nil {
			response.Fail(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
			return
		}

		packages = append(packages, *packageResp)
	}

	response.OK(w, http.StatusCreated, packages)
}

func (h *PackageHandler) UpdatePackageStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctxValue := r.Context().Value(middleware.UserContextKey)
	if ctxValue == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	authUser, ok := ctxValue.(*middleware.AuthUser)
	if !ok || authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authentication context")
		return
	}

	allowedRoles := map[domain.UserRole]bool{
        domain.UserRoleAdmin:       true,
        domain.UserRoleStaffSortir: true,
        domain.UserRoleKurir:       true,
    }
    if !allowedRoles[authUser.Role] {
        response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Access denied")
        return
    }

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	var req domain.UpdatePackageStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	packResp, err := h.packageService.UpdatePackageStatus(id, &req, authUser.Role)
	if err != nil {
		if errors.Is(err, services.ErrPackageNotFound) {
			response.Fail(w, http.StatusNotFound, "PACKAGE_NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, services.ErrForbiddenStatusTransition) {
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to update to this status")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, packResp)
}

func (h *PackageHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctxValue := r.Context().Value(middleware.UserContextKey)
	if ctxValue == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	authUser, ok := ctxValue.(*middleware.AuthUser)
	if !ok || authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authentication context")
		return
	}

	if authUser.Role != domain.UserRoleAdmin {
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to delete a package")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	if err := h.packageService.DeletePackage(id); err != nil {
		response.Fail(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, map[string]string{"message": "package deleted"})
}
