package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type LockerHandler struct {
	lockerService *services.LockerService
}

func NewLockerHandler(lockerService *services.LockerService) *LockerHandler {
	return &LockerHandler{lockerService: lockerService}
}

func (h *LockerHandler) CreateLocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var reqs []domain.CreateLockerRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	var responses []domain.LockerResponse
	for _, req := range reqs {
		res, err := h.lockerService.CreateLocker(&req)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrForbiddenHubAccess):
				response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			case errors.Is(err, services.ErrHubNotFound):
				response.Fail(w, http.StatusNotFound, "HUB_NOT_FOUND", err.Error())
			case errors.Is(err, services.ErrLockerLabelTaken):
				response.Fail(w, http.StatusConflict, "LOCKER_LABEL_TAKEN", err.Error())
			default:
				response.Fail(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
			}
			return
		}
		responses = append(responses, *res)
	}

	response.OK(w, http.StatusCreated, responses)
}

func (h *LockerHandler) GetLockerByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	locker, err := h.lockerService.GetLocker(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, locker)
}

func (h *LockerHandler) ListLockersByHub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	hubIDStr := vars["hubId"]
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub ID")
		return
	}

	lockers, err := h.lockerService.ListLockersByHub(*authUser, hubID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, lockers)
}

func (h *LockerHandler) ListActivePackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	lockerID, err := uuid.Parse(vars["lockerId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	packages, err := h.lockerService.ListActivePackages(lockerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, packages)
}

func (h *LockerHandler) ListAvailableLockers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	service := middleware.GetAuthService(r)
	if service == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Service authentication required")
		return
	}
	hubID, err := uuid.Parse(mux.Vars(r)["hubId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub ID")
		return
	}
	lockers, err := h.lockerService.ListAvailableLockers(*service, hubID)
	if err != nil {
		if errors.Is(err, services.ErrForbiddenHubAccess) {
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}
	response.OK(w, http.StatusOK, lockers)
}

func (h *LockerHandler) CreateClusters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	service := middleware.GetAuthService(r)
	if service == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Service authentication required")
		return
	}
	var req domain.CreateClustersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}
	if err := h.lockerService.CreateClusters(*service, &req); err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrClusterPackageInvalid), errors.Is(err, services.ErrClusterLockerInvalid), errors.Is(err, repository.ErrClusterLockerConflict), errors.Is(err, repository.ErrClusterPackageConflict), errors.Is(err, repository.ErrClusterPackageActive):
			response.Fail(w, http.StatusConflict, "CLUSTER_NOT_ELIGIBLE", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "CLUSTER_FAILED", err.Error())
		}
		return
	}
	response.OK(w, http.StatusCreated, map[string]string{"message": "clusters created"})
}

func (h *LockerHandler) UpdateLocker(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	var req domain.UpdateLockerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	locker, err := h.lockerService.UpdateLocker(*authUser, id, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerLabelTaken):
			response.Fail(w, http.StatusConflict, "LOCKER_LABEL_TAKEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, locker)
}

func (h *LockerHandler) DeleteLocker(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	if err := h.lockerService.DeleteLocker(*authUser, id); err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerNotEmpty):
			response.Fail(w, http.StatusConflict, "LOCKER_NOT_EMPTY", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, map[string]string{"message": "locker deleted"})
}
