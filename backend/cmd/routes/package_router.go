package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterPackageRoutes(router *mux.Router, packageHandler *handler.PackageHandler) {
	packages := router.PathPrefix("/api/v1/packages").Subrouter()

	packages.HandleFunc("", packageHandler.ListPackages).
		Methods(http.MethodGet, http.MethodOptions)

	packages.HandleFunc("/{id}", packageHandler.GetPackageByID).
		Methods(http.MethodGet, http.MethodOptions)

	protected := packages.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("", packageHandler.CreatePackage).
		Methods(http.MethodPost, http.MethodOptions)

	protected.HandleFunc("/{id}/status", packageHandler.UpdatePackageStatus).
		Methods(http.MethodPut, http.MethodOptions)

	protected.HandleFunc("/{id}", packageHandler.DeletePackage).
		Methods(http.MethodDelete, http.MethodOptions)

	authorizedHubPackages := router.PathPrefix("/api/v1/hubs/{hubId}/packages").Subrouter()
	authorizedHubPackages.Use(middleware.AuthMiddleware)
	authorizedHubPackages.HandleFunc("", packageHandler.GetPackagesByHubID).
		Methods(http.MethodGet, http.MethodOptions).
		Queries("status", "{status}").Queries()
}
