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

type LockerScanHandler struct {
	lockerScanService *services.LockerScanService
}

func NewLockerScanHandler(lockerScanService *services.LockerScanService) *LockerScanHandler {
	return &LockerScanHandler{lockerScanService: lockerScanService}
}

func (h *LockerScanHandler) ScanIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req domain.ScanPackageToLockerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.lockerScanService.ScanIn(*authUser, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrPackageNotFound):
			response.Fail(w, http.StatusNotFound, "PACKAGE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerScanAlreadyCheckedIn):
			response.Fail(w, http.StatusConflict, "LOCKER_SCAN_ALREADY_ACTIVE", err.Error())
		case errors.Is(err, services.ErrPackageStatusTransition):
			response.Fail(w, http.StatusConflict, "INVALID_PACKAGE_STATUS", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "SCAN_IN_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusCreated, res)
}

func (h *LockerScanHandler) ScanOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req domain.ScanPackageToLockerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.lockerScanService.ScanOut(*authUser, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerScanNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_SCAN_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerScanAlreadyCheckedOut):
			response.Fail(w, http.StatusConflict, "LOCKER_SCAN_ALREADY_CHECKED_OUT", err.Error())
		case errors.Is(err, services.ErrPackageStatusTransition):
			response.Fail(w, http.StatusConflict, "INVALID_PACKAGE_STATUS", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "SCAN_OUT_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *LockerScanHandler) GetLockerScanByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker scan ID")
		return
	}

	res, err := h.lockerScanService.GetLockerScan(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerScanNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_SCAN_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *LockerScanHandler) ListLockerScansByLocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	lockerIDStr := vars["lockerId"]
	lockerID, err := uuid.Parse(lockerIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	res, err := h.lockerScanService.ListLockerScansByLocker(*authUser, lockerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *LockerScanHandler) ListLockerScansByPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageId"]
	packageID, err := uuid.Parse(packageIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	res, err := h.lockerScanService.ListLockerScansByPackage(*authUser, packageID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLockerScanHistoryForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *LockerScanHandler) ListMyLockerScans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	res, err := h.lockerScanService.ListLockerScansByCourierToday(*authUser)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}
