package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterHubRoutes(router *mux.Router, hubHandler *handler.HubHandler) {
	hubs := router.PathPrefix("/api/v1/hubs").Subrouter()

	hubs.HandleFunc("", hubHandler.ListHubs).
		Methods(http.MethodGet, http.MethodOptions)

	hubs.HandleFunc("/{id}", hubHandler.GetHub).
		Methods(http.MethodGet, http.MethodOptions)

	protected := hubs.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.RoleMiddleware(domain.UserRoleAdmin))

	protected.HandleFunc("", hubHandler.CreateHub).
		Methods(http.MethodPost, http.MethodOptions)

	protected.HandleFunc("/{id}", hubHandler.UpdateHub).
		Methods(http.MethodPut, http.MethodOptions)

	protected.HandleFunc("/{id}", hubHandler.DeleteHub).
		Methods(http.MethodDelete, http.MethodOptions)
}
