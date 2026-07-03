package routes

import (
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
) *mux.Router {

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)

	RegisterAuthRoutes(router, authHandler)
	RegisterAdminRoutes(router)

	return router
}