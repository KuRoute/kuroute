package routes

import (
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/handler"
	"github.com/KuRoute/kuroute/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterDeliveryReportRoutes(router *mux.Router, deliveryReportHandler *handler.DeliveryReportHandler) {
	uploadRoutes := router.PathPrefix("/api/v1/uploads").Subrouter()
	uploadRoutes.Use(middleware.AuthMiddleware)
	uploadRoutes.Use(middleware.RoleMiddleware(domain.UserRoleKurir))
	uploadRoutes.HandleFunc("/delivery-proof", deliveryReportHandler.UploadProof).
		Methods(http.MethodPost, http.MethodOptions)

	stopReportRoutes := router.PathPrefix("/api/v1/route-stops/{id}").Subrouter()
	stopReportRoutes.Use(middleware.AuthMiddleware)
	stopReportRoutes.Use(middleware.RoleMiddleware(domain.UserRoleKurir))
	stopReportRoutes.HandleFunc("/delivery-report", deliveryReportHandler.SubmitDeliveryReport).
		Methods(http.MethodPost, http.MethodOptions)

	deliveryReportRoutes := router.PathPrefix("/api/v1/delivery-reports").Subrouter()
	deliveryReportRoutes.Use(middleware.AuthMiddleware)
	deliveryReportRoutes.HandleFunc("/{id}", deliveryReportHandler.GetDeliveryReportByID).
		Methods(http.MethodGet, http.MethodOptions)

	routeReportRoutes := router.PathPrefix("/api/v1/routes/{routeId}").Subrouter()
	routeReportRoutes.Use(middleware.AuthMiddleware)
	routeReportRoutes.HandleFunc("/delivery-reports", deliveryReportHandler.ListDeliveryReportsByRoute).
		Methods(http.MethodGet, http.MethodOptions)

	packageReportRoutes := router.PathPrefix("/api/v1/packages/{packageId}").Subrouter()
	packageReportRoutes.Use(middleware.AuthMiddleware)
	packageReportRoutes.Use(middleware.RoleMiddleware(domain.UserRoleAdmin, domain.UserRoleStaffSortir))
	packageReportRoutes.HandleFunc("/delivery-reports", deliveryReportHandler.ListDeliveryReportsByPackage).
		Methods(http.MethodGet, http.MethodOptions)
}
