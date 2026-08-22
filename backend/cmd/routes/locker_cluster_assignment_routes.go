package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterLockerClusterAssignmentRoutes(router *mux.Router, assignmentHandler *handler.LockerClusterAssignmentHandler) {
	routes := router.PathPrefix("/api/v1/locker-cluster-assignments").Subrouter()
	routes.Use(middleware.AuthMiddleware)
	routes.Use(middleware.RoleMiddleware(domain.UserRoleStaffSortir))
	routes.HandleFunc("/packages/{packageId}", assignmentHandler.GetByPackageID).
		Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("", assignmentHandler.ListByLockerID).
		Methods(http.MethodGet, http.MethodOptions)
}
