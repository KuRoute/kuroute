package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CourierBatchHandler struct {
	courierBatchService *services.CourierBatchService
}

func NewCourierBatchHandler(courierBatchService *services.CourierBatchService) *CourierBatchHandler {
	return &CourierBatchHandler{courierBatchService: courierBatchService}
}

func (h *CourierBatchHandler) GetCourierBatchByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid courier batch ID")
		return
	}

	res, err := h.courierBatchService.GetCourierBatch(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCourierBatchNotFound):
			response.Fail(w, http.StatusNotFound, "COURIER_BATCH_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *CourierBatchHandler) ListCourierBatchesByHub(w http.ResponseWriter, r *http.Request) {
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

	var status *domain.BatchStatus
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		parsed := domain.BatchStatus(statusStr)
		switch parsed {
		case domain.BatchStatusPendingRoute, domain.BatchStatusRouteReady, domain.BatchStatusInProgress, domain.BatchStatusCompleted:
			status = &parsed
		default:
			response.Fail(w, http.StatusBadRequest, "INVALID_STATUS", "Invalid status filter")
			return
		}
	}

	var batchDate *time.Time
	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.Fail(w, http.StatusBadRequest, "INVALID_DATE", "Invalid date format, expected YYYY-MM-DD")
			return
		}
		batchDate = &parsedDate
	}

	res, err := h.courierBatchService.ListCourierBatchesByHub(*authUser, hubID, status, batchDate)
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

func (h *CourierBatchHandler) ListMyCourierBatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	res, err := h.courierBatchService.ListMyCourierBatches(*authUser)
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
