package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handler.UserHandler) {
	user := router.PathPrefix("/api/v1/user").Subrouter()

	user.HandleFunc("/profile/{id}", userHandler.GetUserByID).
		Methods(http.MethodGet, http.MethodOptions)

	protected := user.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.RoleMiddleware(domain.UserRoleAdmin))

	protected.HandleFunc("/profile/{id}", userHandler.UpdateUser).
		Methods(http.MethodPut, http.MethodOptions)

	protected.HandleFunc("", userHandler.CreateUser).
		Methods(http.MethodPost, http.MethodOptions)

}