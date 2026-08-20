package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterRouteStopRoutes(router *mux.Router, routeStopHandler *handler.RouteStopHandler) {
	stopRoutes := router.PathPrefix("/api/v1/route-stops").Subrouter()
	stopRoutes.Use(middleware.AuthMiddleware)
	stopRoutes.HandleFunc("/{id}", routeStopHandler.GetRouteStopByID).
		Methods(http.MethodGet, http.MethodOptions)

	routeRoutes := router.PathPrefix("/api/v1/routes/{routeId}/stops").Subrouter()
	routeRoutes.Use(middleware.AuthMiddleware)
	routeRoutes.HandleFunc("", routeStopHandler.ListRouteStopsByRoute).
		Methods(http.MethodGet, http.MethodOptions)
}
