package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterCourierBatchRoutes(router *mux.Router, courierBatchHandler *handler.CourierBatchHandler) {
	courierBatchRoutes := router.PathPrefix("/api/v1").Subrouter()
	courierBatchRoutes.Use(middleware.AuthMiddleware)

	courierBatchRoutes.HandleFunc("/courier-batches/{id}", courierBatchHandler.GetCourierBatchByID).
		Methods(http.MethodGet, http.MethodOptions)

	hubCourierBatches := router.PathPrefix("/api/v1/hubs/{hubId}/courier-batches").Subrouter()
	hubCourierBatches.Use(middleware.AuthMiddleware)
	hubCourierBatches.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	hubCourierBatches.HandleFunc("", courierBatchHandler.ListCourierBatchesByHub).
		Methods(http.MethodGet, http.MethodOptions)

	courierMeBatches := router.PathPrefix("/api/v1/couriers").Subrouter()
	courierMeBatches.Use(middleware.AuthMiddleware)
	courierMeBatches.Use(middleware.RoleMiddleware(domain.UserRoleKurir))
	courierMeBatches.HandleFunc("/me/courier-batches", courierBatchHandler.ListMyCourierBatches).
		Methods(http.MethodGet, http.MethodOptions)
}
