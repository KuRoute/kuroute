package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format")
		return
	}

	userResp, err := h.userService.GetUserByID(idUUID)
	if err != nil {
		log.Printf("[ERROR] GetUserByID failed: %v", err)
		response.Fail(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found or unavailable")
		return
	}

	response.OK(w, http.StatusOK, userResp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
        response.Fail(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to update user information")
        return
    }

    vars := mux.Vars(r)
    idStr := vars["id"]

    idUUID, err := uuid.Parse(idStr)
    if err != nil {
        response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format")
        return
    }

    var req domain.UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
        return
    }

	req.ID = idUUID

    err = h.userService.UpdateUser(&req, authUser.Role)
    if err != nil {
        log.Printf("[ERROR] UpdateUser failed for ID %s: %v", idUUID, err)
        response.Fail(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update user information")
        return
    }

    response.OK(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	var reqs []domain.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	if len(reqs) == 0 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Request body cannot be empty")
		return
	}

	var users []domain.UserResponse

	for _, req := range reqs {

		if req.HubID == [16]byte{} ||
			req.Name == "" ||
			req.Email == "" ||
			req.Phone == "" ||
			req.Role == "" ||
			req.Password == "" {

			response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "All fields are required")
			return
		}

		if len(req.Password) < 8 {
			response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters")
			return
		}

		userResp, err := h.userService.CreateUser(&req, authUser.Role)
		if err != nil {
			response.Fail(w, http.StatusConflict, "REGISTRATION_FAILED", err.Error())
			return
		}

		users = append(users, *userResp)
	}

	response.OK(w, http.StatusCreated, users)
}