package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterBatchAssignmentRoutes(router *mux.Router, batchAssignmentHandler *handler.BatchAssignmentHandler) {
	batchAssignments := router.PathPrefix("/api/v1/batch-assignments").Subrouter()
	batchAssignments.Use(middleware.AuthMiddleware)

	batchAssignments.HandleFunc("/{id}", batchAssignmentHandler.GetBatchAssignmentByID).
		Methods(http.MethodGet, http.MethodOptions)
	batchAssignments.HandleFunc("/{id}/route", batchAssignmentHandler.GetBatchAssignmentRoute).
		Methods(http.MethodGet, http.MethodOptions)
	batchAssignments.HandleFunc("/{id}/start", batchAssignmentHandler.StartBatch).
		Methods(http.MethodPost, http.MethodOptions)

	staffProtected := batchAssignments.NewRoute().Subrouter()
	staffProtected.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	staffProtected.HandleFunc("", batchAssignmentHandler.CreateBatchAssignment).
		Methods(http.MethodPost, http.MethodOptions)
	staffProtected.HandleFunc("/{id}", batchAssignmentHandler.UpdateBatchAssignment).
		Methods(http.MethodPut, http.MethodOptions)
	staffProtected.HandleFunc("/{id}", batchAssignmentHandler.DeleteBatchAssignment).
		Methods(http.MethodDelete, http.MethodOptions)

	hubBatchAssignments := router.PathPrefix("/api/v1/hubs/{hubId}/batch-assignments").Subrouter()
	hubBatchAssignments.Use(middleware.AuthMiddleware)
	hubBatchAssignments.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	hubBatchAssignments.HandleFunc("", batchAssignmentHandler.ListBatchAssignmentsByHub).
		Methods(http.MethodGet, http.MethodOptions)

	courierAssignments := router.PathPrefix("/api/v1/couriers").Subrouter()
	courierAssignments.Use(middleware.AuthMiddleware)
	courierAssignments.Use(middleware.RoleMiddleware(domain.UserRoleKurir))
	courierAssignments.HandleFunc("/me/batch-assignments", batchAssignmentHandler.ListMyBatchAssignments).
		Methods(http.MethodGet, http.MethodOptions)
}
