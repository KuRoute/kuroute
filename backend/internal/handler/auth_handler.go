package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/services"
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
	var req domain.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.HubID == [16]byte{} || req.Name == "" || req.Email == "" || req.Phone == "" || req.Role == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "All fields are required")
		return
	}

	if req.Role != domain.UserRoleStaffSortir && req.Role != domain.UserRoleKurir {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Role must be staff_sortir or kurir")
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
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required")
		return
	}

	result, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}

	// Check if result is a setup token response (first-time login)
	if setupResp, ok := result.(*domain.SetupTokenResponse); ok {
		response.OK(w, http.StatusOK, map[string]interface{}{
			"setupToken":         setupResp.SetupToken,
			"user":               setupResp.User,
		})
		return
	}

	// Normal login response
	if loginResp, ok := result.(*domain.LoginResponse); ok {
		response.OK(w, http.StatusOK, loginResp)
		return
	}

	response.Fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected response type")
}

// POST /api/v1/auth/setup-password
func (h *AuthHandler) SetupPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.SetupPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.SetupToken == "" || req.NewPassword == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Setup token and new password are required")
		return
	}

	if len(req.NewPassword) < 8 {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters")
		return
	}

	resp, err := h.authService.SetupFirstPassword(req.SetupToken, req.NewPassword)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "SETUP_FAILED", err.Error())
		return
	}

	response.OK(w, http.StatusOK, resp)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
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