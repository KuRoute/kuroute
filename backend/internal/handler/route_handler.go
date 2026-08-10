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

type RouteHandler struct {
	routeService *services.RouteService
}

func NewRouteHandler(routeService *services.RouteService) *RouteHandler {
	return &RouteHandler{routeService: routeService}
}

func (h *RouteHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authService := middleware.GetAuthService(r)
	if authService == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Service authentication required")
		return
	}

	var req domain.ComputeRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.routeService.CreateRoute(*authService, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenServiceAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrCourierBatchNotFound):
			response.Fail(w, http.StatusNotFound, "COURIER_BATCH_NOT_FOUND", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusCreated, res)
}

func (h *RouteHandler) GetRouteByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid route ID")
		return
	}

	res, err := h.routeService.GetRoute(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *RouteHandler) ListRoutesByHub(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.routeService.ListRoutesByHub(*authUser, hubID, batchDate)
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
