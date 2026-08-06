package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CourierProfileHandler struct {
	courierProfileService *services.CourierProfileService
}

func NewCourierProfileHandler(courierProfileService *services.CourierProfileService) *CourierProfileHandler {
	return &CourierProfileHandler{courierProfileService: courierProfileService}
}

func (h *CourierProfileHandler) GetCourierProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	courierID, err := uuid.Parse(vars["courierId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid courier ID format")
		return
	}

	profile, err := h.courierProfileService.GetCourierProfile(authUser, courierID)
	if err != nil {
		if err == services.ErrForbidden {
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
			return
		}
		if err == services.ErrCourierProfileNotFound {
			response.Fail(w, http.StatusNotFound, "NOT_FOUND", "Courier profile not found")
			return
		}
		log.Printf("[ERROR] GetCourierProfile failed: %v", err)
		response.Fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve courier profile")
		return
	}

	response.OK(w, http.StatusOK, profile)
}

func (h *CourierProfileHandler) CreateCourierProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	if authUser.Role != domain.UserRoleAdmin {
		response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
		return
	}

	vars := mux.Vars(r)
	courierID, err := uuid.Parse(vars["courierId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid courier ID format")
		return
	}

	var req domain.CreateCourierProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	if req.VehicleType == "" || req.VehiclePlate == "" {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "vehicleType and vehiclePlate are required")
		return
	}

	createdProfile, err := h.courierProfileService.CreateCourierProfile(courierID, &req)
	if err != nil {
		if err == services.ErrCourierProfileNotFound {
			response.Fail(w, http.StatusNotFound, "NOT_FOUND", "Courier not found")
			return
		}
		if err == services.ErrCourierProfileAlreadyExists {
			response.Fail(w, http.StatusConflict, "CONFLICT", "Courier profile already exists")
			return
		}
		log.Printf("[ERROR] CreateCourierProfile failed: %v", err)
		response.Fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create courier profile")
		return
	}

	response.OK(w, http.StatusCreated, createdProfile)
}

func (h *CourierProfileHandler) UpdateCourierProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	courierID, err := uuid.Parse(vars["courierId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid courier ID format")
		return
	}

	var req domain.UpdateCourierProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	if req.VehicleType == nil && req.VehiclePlate == nil {
		response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "At least one field is required")
		return
	}

	profile, err := h.courierProfileService.UpdateCourierProfile(courierID, &req)
	if err != nil {
		if err == services.ErrCourierProfileNotFound {
			response.Fail(w, http.StatusNotFound, "NOT_FOUND", "Courier not found")
			return
		}
		log.Printf("[ERROR] UpdateCourierProfile failed: %v", err)
		response.Fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update courier profile")
		return
	}

	response.OK(w, http.StatusOK, profile)
}

func (h *CourierProfileHandler) ListCouriersByHub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	hubID, err := uuid.Parse(vars["hubId"])
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid hub ID format")
		return
	}

	var vehicleType *domain.VehicleType
	if vt := r.URL.Query().Get("vehicleType"); vt != "" {
		parsedType := domain.VehicleType(vt)
		switch parsedType {
		case domain.VehicleTypeMotor, domain.VehicleTypeMobil, domain.VehicleTypeVan:
			vehicleType = &parsedType
		default:
			response.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid vehicleType")
			return
		}
	}

	couriers, err := h.courierProfileService.ListCouriersByHub(authUser, hubID, vehicleType)
	if err != nil {
		if err == services.ErrForbidden {
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
			return
		}
		log.Printf("[ERROR] ListCouriersByHub failed: %v", err)
		response.Fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list couriers")
		return
	}

	response.OK(w, http.StatusOK, couriers)
}
