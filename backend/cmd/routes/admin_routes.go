package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(router *mux.Router) {

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RoleMiddleware(domain.UserRoleStaffSortir))

	admin.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin users"))
	}).Methods("GET")
}