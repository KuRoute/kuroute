package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HubHandler struct {
	hubService *services.HubService
}

func NewHubHandler(hubService *services.HubService) *HubHandler {
	return &HubHandler{hubService: hubService}
}

// POST /api/v1/hubs
func (h *HubHandler) CreateHub(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to create a new user")
		return
	}

	var reqs []domain.CreateHubRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if len(reqs) == 0 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Request body cannot be empty")
		return
	}

	var hubs []domain.HubResponse

	for _, req := range reqs {
		if req.Name == "" || req.City == "" {
			response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "name and city are required")
			return
		}

		hubResp, err := h.hubService.CreateHub(&req)
		if err != nil {
			response.Fail(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
			return
		}

		hubs = append(hubs, *hubResp)
	}

	response.OK(w, http.StatusCreated, hubs)
}

// GET /api/v1/hubs
func (h *HubHandler) ListHubs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hubs, err := h.hubService.ListHubs()
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}
	response.OK(w, http.StatusOK, hubs)
}

// GET /api/v1/hubs/{id}
func (h *HubHandler) GetHub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub id")
		return
	}

	hubResp, err := h.hubService.GetHubByID(id)
	if err != nil {
		response.Fail(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	response.OK(w, http.StatusOK, hubResp)
}

// PUT /api/v1/hubs/{id}
func (h *HubHandler) UpdateHub(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to create a new user")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub id")
		return
	}

	var req domain.UpdateHubRequest

	if req.Name == "" && req.City == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "At least one of name or city must be provided")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	hubResp, err := h.hubService.UpdateHub(id, &req)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, hubResp)
}

// DELETE /api/v1/hubs/{id}
func (h *HubHandler) DeleteHub(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to create a new user")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid hub id")
		return
	}

	if err := h.hubService.DeleteHub(id); err != nil {
		response.Fail(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, map[string]string{"message": "hub deleted"})
}
