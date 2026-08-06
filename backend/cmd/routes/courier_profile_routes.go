package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterCourierProfileRoutes(router *mux.Router, courierProfileHandler *handler.CourierProfileHandler) {
	courierRoutes := router.PathPrefix("/api/v1/couriers").Subrouter()
	courierRoutes.Use(middleware.AuthMiddleware)

	courierRoutes.HandleFunc("/{courierId}/profile", courierProfileHandler.GetCourierProfile).
		Methods(http.MethodGet, http.MethodOptions)

	adminProtected := courierRoutes.NewRoute().Subrouter()
	adminProtected.Use(middleware.RoleMiddleware(domain.UserRoleAdmin))
	adminProtected.HandleFunc("/{courierId}/profile", courierProfileHandler.CreateCourierProfile).
		Methods(http.MethodPost, http.MethodOptions)

	adminProtected.HandleFunc("/{courierId}/profile", courierProfileHandler.UpdateCourierProfile).
		Methods(http.MethodPut, http.MethodOptions)

	hubCouriers := router.PathPrefix("/api/v1/hubs/{hubId}/couriers").Subrouter()
	hubCouriers.Use(middleware.AuthMiddleware)
	hubCouriers.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	hubCouriers.HandleFunc("", courierProfileHandler.ListCouriersByHub).
		Methods(http.MethodGet, http.MethodOptions)
}
