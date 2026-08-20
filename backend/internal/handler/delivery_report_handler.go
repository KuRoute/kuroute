package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DeliveryReportHandler struct {
	deliveryReportService *services.DeliveryReportService
}

func NewDeliveryReportHandler(deliveryReportService *services.DeliveryReportService) *DeliveryReportHandler {
	return &DeliveryReportHandler{deliveryReportService: deliveryReportService}
}

func (h *DeliveryReportHandler) UploadProof(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_UPLOAD", "Upload must be multipart form-data and max 5MB")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_UPLOAD", "Missing file in field 'file'")
		return
	}
	defer file.Close()

	url, err := h.deliveryReportService.UploadProof(*authUser, file, header)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDeliveryReportForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrUploadTooLarge):
			response.Fail(w, http.StatusBadRequest, "INVALID_UPLOAD", err.Error())
		case errors.Is(err, services.ErrUploadInvalidFileType):
			response.Fail(w, http.StatusBadRequest, "INVALID_UPLOAD", err.Error())
		case errors.Is(err, services.ErrDeliverySubmitInvalid):
			response.Fail(w, http.StatusBadRequest, "INVALID_UPLOAD", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusCreated, map[string]string{"url": url})
}

func (h *DeliveryReportHandler) SubmitDeliveryReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	stopIDStr := vars["id"]
	stopID, err := uuid.Parse(stopIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid route stop ID")
		return
	}

	var req domain.SubmitDeliveryReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON request body")
		return
	}

	res, err := h.deliveryReportService.SubmitDeliveryReport(*authUser, stopID, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRouteStopNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_STOP_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrRouteStopAccessForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, services.ErrBatchNotInProgress):
			response.Fail(w, http.StatusConflict, "BATCH_NOT_IN_PROGRESS", err.Error())
		case errors.Is(err, services.ErrDeliveryStatusTransition):
			response.Fail(w, http.StatusConflict, "INVALID_DELIVERY_STATUS", err.Error())
		case errors.Is(err, services.ErrDeliveryReportForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, domain.ErrFailureReasonOnSuccess):
			response.Fail(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "SUBMIT_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusCreated, res)
}

func (h *DeliveryReportHandler) GetDeliveryReportByID(w http.ResponseWriter, r *http.Request) {
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
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid delivery report ID")
		return
	}

	res, err := h.deliveryReportService.GetDeliveryReport(*authUser, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDeliveryReportNotFound):
			response.Fail(w, http.StatusNotFound, "DELIVERY_REPORT_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrRouteStopNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_STOP_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrDeliveryReportForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *DeliveryReportHandler) ListDeliveryReportsByRoute(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.deliveryReportService.ListDeliveryReportsByRoute(*authUser, routeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRouteNotFound):
			response.Fail(w, http.StatusNotFound, "ROUTE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrDeliveryReportForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}

func (h *DeliveryReportHandler) ListDeliveryReportsByPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authUser := middleware.GetAuthUser(r)
	if authUser == nil {
		response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageId"]
	packageID, err := uuid.Parse(packageIDStr)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID")
		return
	}

	res, err := h.deliveryReportService.ListDeliveryReportsByPackage(*authUser, packageID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPackageNotFound):
			response.Fail(w, http.StatusNotFound, "PACKAGE_NOT_FOUND", err.Error())
		case errors.Is(err, services.ErrDeliveryReportForbidden):
			response.Fail(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.Fail(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		}
		return
	}

	response.OK(w, http.StatusOK, res)
}
