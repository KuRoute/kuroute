package routes

import (
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	hubHandler *handler.HubHandler,
	userHandler *handler.UserHandler,
) *mux.Router {

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)

	RegisterAuthRoutes(router, authHandler)
	RegisterHubRoutes(router, hubHandler)
	RegisterAdminRoutes(router)
	RegisterUserRoutes(router, userHandler)

	return router
}
