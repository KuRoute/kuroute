package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RouteStopHandler struct {
	routeStopService *services.RouteStopService
}

func NewRouteStopHandler(routeStopService *services.RouteStopService) *RouteStopHandler {
	return &RouteStopHandler{routeStopService: routeStopService}
}

func (h *RouteStopHandler) GetRouteStopByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid route stop ID")
		return
	}

	res, err := h.routeStopService.GetRouteStop(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRouteStopNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_STOP_NOT_FOUND", err.Error())
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

func (h *RouteStopHandler) ListRouteStopsByRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	routeIDStr := vars["routeId"]
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid route ID")
		return
	}

	res, err := h.routeStopService.ListRouteStops(*authUser, routeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrForbiddenHubAccess):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *RouteStopHandler) GetRouteStops(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w)
	_ = h
	_ = r
}
