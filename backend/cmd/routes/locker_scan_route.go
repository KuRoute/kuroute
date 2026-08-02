package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterLockerScanRoutes(router *mux.Router, lockerScanHandler *handler.LockerScanHandler) {
	lockerScans := router.PathPrefix("/api/v1/locker-scans").Subrouter()
	lockerScans.Use(middleware.AuthMiddleware)

	lockerScans.HandleFunc("/scan-in", lockerScanHandler.ScanIn).
		Methods(http.MethodPost, http.MethodOptions)
	lockerScans.HandleFunc("/scan-out", lockerScanHandler.ScanOut).
		Methods(http.MethodPost, http.MethodOptions)
	lockerScans.HandleFunc("/{id}", lockerScanHandler.GetLockerScanByID).
		Methods(http.MethodGet, http.MethodOptions)

	lockerByLocker := router.PathPrefix("/api/v1/lockers/{lockerId}/scans").Subrouter()
	lockerByLocker.Use(middleware.AuthMiddleware)
	lockerByLocker.HandleFunc("", lockerScanHandler.ListLockerScansByLocker).
		Methods(http.MethodGet, http.MethodOptions)

	packageHistory := router.PathPrefix("/api/v1/packages/{packageId}/locker-scans").Subrouter()
	packageHistory.Use(middleware.AuthMiddleware)
	packageHistory.HandleFunc("", lockerScanHandler.ListLockerScansByPackage).
		Methods(http.MethodGet, http.MethodOptions)

	courierHistory := router.PathPrefix("/api/v1/couriers").Subrouter()
	courierHistory.Use(middleware.AuthMiddleware)
	courierHistory.HandleFunc("/me/locker-scans", lockerScanHandler.ListMyLockerScans).
		Methods(http.MethodGet, http.MethodOptions)
}
