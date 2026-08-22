package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterLockerRoutes(router *mux.Router, lockerHandler *handler.LockerHandler) {
	clusterRoutes := router.PathPrefix("/api/v1/internal/lockers").Subrouter()
	clusterRoutes.Use(middleware.ServiceMiddleware)
	clusterRoutes.Use(middleware.NameServiceMiddleware(domain.ServiceCluster))
	clusterRoutes.HandleFunc("/hubs/{hubId}/lockers/available", lockerHandler.ListAvailableLockers).
		Methods(http.MethodGet, http.MethodOptions)
	clusterRoutes.HandleFunc("/clusters", lockerHandler.CreateClusters).
		Methods(http.MethodPost, http.MethodOptions)

	internalRoutes := router.PathPrefix("/internal/lockers/{lockerId}").Subrouter()
	internalRoutes.Use(middleware.ServiceMiddleware)
	internalRoutes.Use(middleware.NameServiceMiddleware(domain.ServiceRoute))
	internalRoutes.HandleFunc("/active-packages", lockerHandler.ListActivePackages).
		Methods(http.MethodGet, http.MethodOptions)

	lockers := router.PathPrefix("/api/v1/lockers").Subrouter()

	protected := lockers.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/{id}", lockerHandler.GetLockerByID).
		Methods(http.MethodGet, http.MethodOptions)

	protected.HandleFunc("", lockerHandler.CreateLocker).
		Methods(http.MethodPost, http.MethodOptions)

	protected.HandleFunc("/{id}", lockerHandler.UpdateLocker).
		Methods(http.MethodPut, http.MethodOptions)

	protected.HandleFunc("/{id}", lockerHandler.DeleteLocker).
		Methods(http.MethodDelete, http.MethodOptions)

	hubLockers := router.PathPrefix("/api/v1/hubs/{hubId}/lockers").Subrouter()
	hubLockers.Use(middleware.AuthMiddleware)
	hubLockers.HandleFunc("", lockerHandler.ListLockersByHub).
		Methods(http.MethodGet, http.MethodOptions)
}
