package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handler.AuthHandler) {

	auth := router.PathPrefix("/api/v1/auth").Subrouter()

	auth.HandleFunc("/login", authHandler.Login).
		Methods(http.MethodPost, http.MethodOptions)

	auth.HandleFunc("/register", authHandler.Register).
		Methods(http.MethodPost, http.MethodOptions)

	auth.HandleFunc("/change-password", authHandler.ChangePassword).
		Methods(http.MethodPost, http.MethodOptions)

	auth.HandleFunc("/refresh", authHandler.RefreshToken).
		Methods(http.MethodPost, http.MethodOptions)

	auth.HandleFunc("/logout", authHandler.Logout).
		Methods(http.MethodPost, http.MethodOptions)

}