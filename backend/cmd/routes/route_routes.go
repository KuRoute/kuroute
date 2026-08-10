package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterRouteRoutes(router *mux.Router, routeHandler *handler.RouteHandler) {
	internalRoutes := router.PathPrefix("/api/v1/internal/routes").Subrouter()
	internalRoutes.Use(middleware.ServiceMiddleware)
	internalRoutes.Use(middleware.NameServiceMiddleware(domain.ServiceRoute))
	internalRoutes.HandleFunc("", routeHandler.CreateRoute).
		Methods(http.MethodPost, http.MethodOptions)

	userRoutes := router.PathPrefix("/api/v1/routes").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware)
	userRoutes.HandleFunc("/{id}", routeHandler.GetRouteByID).
		Methods(http.MethodGet, http.MethodOptions)

	hubRoutes := router.PathPrefix("/api/v1/hubs/{hubId}/routes").Subrouter()
	hubRoutes.Use(middleware.AuthMiddleware)
	hubRoutes.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	hubRoutes.HandleFunc("", routeHandler.ListRoutesByHub).
		Methods(http.MethodGet, http.MethodOptions)
}
