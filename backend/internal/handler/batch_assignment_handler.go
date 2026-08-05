package handler

import (
	"encoding/json"
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

type BatchAssignmentHandler struct {
	batchAssignmentService *services.BatchAssignmentService
}

func NewBatchAssignmentHandler(batchAssignmentService *services.BatchAssignmentService) *BatchAssignmentHandler {
	return &BatchAssignmentHandler{batchAssignmentService: batchAssignmentService}
}

func (h *BatchAssignmentHandler) CreateBatchAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req domain.AssignLockerToCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.batchAssignmentService.CreateAssignment(*authUser, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrLockerNotFound):
			response.Fail(w, http.StatusNotFound, "LOCKER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrCourierNotFound):
			response.Fail(w, http.StatusNotFound, "COURIER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrLockerAlreadyAssigned), errors.Is(err, services.ErrCourierAlreadyAssigned):
			response.Fail(w, http.StatusConflict, "BATCH_ASSIGNMENT_CONFLICT", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusCreated, res)
}

func (h *BatchAssignmentHandler) GetBatchAssignmentByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid batch assignment ID")
		return
	}

	res, err := h.batchAssignmentService.GetAssignment(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "BATCH_ASSIGNMENT_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *BatchAssignmentHandler) ListBatchAssignmentsByHub(w http.ResponseWriter, r *http.Request) {
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

	var batchDate *time.Time
	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.Fail(w, http.StatusBadRequest, "INVALID_DATE", "Invalid date format, expected YYYY-MM-DD")
			return
		}
		batchDate = &parsedDate
	}

	var courierID *uuid.UUID
	if courierIDStr := r.URL.Query().Get("courierId"); courierIDStr != "" {
		parsedCourierID, err := uuid.Parse(courierIDStr)
		if err != nil {
			response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid courier ID")
			return
		}
		courierID = &parsedCourierID
	}

	res, err := h.batchAssignmentService.ListAssignmentsByHub(*authUser, hubID, batchDate, courierID)
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

func (h *BatchAssignmentHandler) ListMyBatchAssignments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
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

	res, err := h.batchAssignmentService.ListMyAssignments(*authUser, batchDate)
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

func (h *BatchAssignmentHandler) UpdateBatchAssignment(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid batch assignment ID")
		return
	}

	var req domain.UpdateBatchAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.batchAssignmentService.UpdateAssignment(*authUser, id, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "BATCH_ASSIGNMENT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrCourierNotFound):
			response.Fail(w, http.StatusNotFound, "COURIER_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentLocked):
			response.Fail(w, http.StatusConflict, "BATCH_ASSIGNMENT_LOCKED", err.Error())
		case errors.Is(err, services.ErrLockerAlreadyAssigned), errors.Is(err, services.ErrCourierAlreadyAssigned):
			response.Fail(w, http.StatusConflict, "BATCH_ASSIGNMENT_CONFLICT", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *BatchAssignmentHandler) DeleteBatchAssignment(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid batch assignment ID")
		return
	}

	if err := h.batchAssignmentService.DeleteAssignment(*authUser, id); err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "BATCH_ASSIGNMENT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentLocked):
			response.Fail(w, http.StatusConflict, "BATCH_ASSIGNMENT_LOCKED", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, map[string]string{"message": "batch assignment deleted"})
}

func (h *BatchAssignmentHandler) StartBatch(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid batch assignment ID")
		return
	}

	res, err := h.batchAssignmentService.StartBatch(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "BATCH_ASSIGNMENT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrBatchAlreadyStarted):
			response.Fail(w, http.StatusConflict, "BATCH_ALREADY_STARTED", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "START_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *BatchAssignmentHandler) GetBatchAssignmentRoute(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid batch assignment ID")
		return
	}

	res, err := h.batchAssignmentService.GetRouteForAssignment(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchAssignmentNotFound):
			response.Fail(w, http.StatusNotFound, "BATCH_ASSIGNMENT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_ROUTE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}
