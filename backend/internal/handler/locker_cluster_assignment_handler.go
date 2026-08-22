package handler

import (
	"errors"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type LockerClusterAssignmentHandler struct {
	service *services.LockerClusterAssignmentService
}

func NewLockerClusterAssignmentHandler(service *services.LockerClusterAssignmentService) *LockerClusterAssignmentHandler {
	return &LockerClusterAssignmentHandler{service: service}
}

func (h *LockerClusterAssignmentHandler) GetByPackageID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	packageID, err := uuid.Parse(mux.Vars(r)["packageId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	result, err := h.service.GetByPackageID(*authUser, packageID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLockerClusterAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "ASSIGNMENT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerClusterAssignmentForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}
	response.OK(w, http.StatusOK, result)
}

func (h *LockerClusterAssignmentHandler) ListByLockerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	lockerID, err := uuid.Parse(r.URL.Query().Get("lockerId"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid locker ID")
		return
	}

	result, err := h.service.ListByLockerID(*authUser, lockerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLockerClusterAssignmentForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}
	response.OK(w, http.StatusOK, result)
}
