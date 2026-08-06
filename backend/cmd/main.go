package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KuRoute/kuroute/backend/cmd/routes"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/KuRoute/kuroute/backend/internal/services"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	hubRepo := repository.NewHubRepository(db)
	packageRepo := repository.NewPackageRepository(db)
	lockerRepo := repository.NewLockerRepository(db)
	lockerScanRepo := repository.NewLockerScanRepository(db)
	batchAssignmentRepo := repository.NewBatchAssignmentRepository(db)
	courierProfileRepo := repository.NewCourierProfileRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	hubService := services.NewHubService(hubRepo)
	userService := services.NewUserService(userRepo)
	packageService := services.NewPackageService(packageRepo)
	lockerService := services.NewLockerService(lockerRepo, hubRepo)
	lockerScanService := services.NewLockerScanService(lockerScanRepo, lockerRepo, packageRepo)
	batchAssignmentService := services.NewBatchAssignmentService(batchAssignmentRepo, lockerScanRepo, lockerRepo, packageRepo, userRepo)
	courierProfileService := services.NewCourierProfileService(courierProfileRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	hubHandler := handler.NewHubHandler(hubService)
	userHandler := handler.NewUserHandler(userService)
	packageHandler := handler.NewPackageHandler(packageService)
	lockerHandler := handler.NewLockerHandler(lockerService)
	lockerScanHandler := handler.NewLockerScanHandler(lockerScanService)
	batchAssignmentHandler := handler.NewBatchAssignmentHandler(batchAssignmentService)
	courierProfileHandler := handler.NewCourierProfileHandler(courierProfileService)

	// Setup router
	router := routes.SetupRouter(authHandler, hubHandler, userHandler, packageHandler, lockerHandler, lockerScanHandler, batchAssignmentHandler, courierProfileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
