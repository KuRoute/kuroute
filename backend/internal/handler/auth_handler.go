package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/jwt"
	"github.com/KuRoute/kuroute/backend/package/response"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.HubID == [16]byte{} || req.Name == "" || req.Email == "" || req.Phone == "" || req.Role == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "All fields are required")
		return
	}

	if req.Role != domain.UserRoleAdmin {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Role must be admin")
		return
	}

	if len(req.Password) < 8 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters")
		return
	}

	userResp, err := h.authService.RegisterUser(&req)
	if err != nil {
		response.Fail(w, http.StatusConflict, "REGISTRATION_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusCreated, userResp)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required")
		return
	}

	loginResp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, loginResp)
}

// POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Current password and new password are required")
		return
	}

	if len(req.NewPassword) < 8 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters")
		return
	}

	tokenString := jwt.ExtractToken(r)
	if tokenString == "" {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization token")
		return
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
		return
	}

	resp, err := h.authService.ChangePassword(claims.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "PASSWORD_CHANGE_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, resp)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required")
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, "REFRESH_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, resp)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required")
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "LOGOUT_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
