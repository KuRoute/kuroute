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
	packageHandler *handler.PackageHandler,
	lockerHandler *handler.LockerHandler,
	lockerScanHandler *handler.LockerScanHandler,
	batchAssignmentHandler *handler.BatchAssignmentHandler,
	courierProfileHandler *handler.CourierProfileHandler,
) *mux.Router {

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)

	RegisterAuthRoutes(router, authHandler)
	RegisterHubRoutes(router, hubHandler)
	RegisterAdminRoutes(router)
	RegisterUserRoutes(router, userHandler)
	RegisterPackageRoutes(router, packageHandler)
	RegisterLockerRoutes(router, lockerHandler)
	RegisterLockerScanRoutes(router, lockerScanHandler)
	RegisterBatchAssignmentRoutes(router, batchAssignmentHandler)
	RegisterCourierProfileRoutes(router, courierProfileHandler)

	return router
}
