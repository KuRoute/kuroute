package domain

import (
	"time"

	"github.com/google/uuid"
)

// Model

// Route holds the result of a single solver run for a CourierBatch.
// algorithm_used is TEXT (not enum) because the algorithm choice
// is under active research — see ADR in docs/adr/002-routing-algorithm.md.
// Once the algorithm is finalized, add a CHECK constraint via migration.
type Route struct {
	Base
	CourierBatchID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"courierBatchId"`
	AlgorithmUsed    string    `gorm:"type:text;not null;default:'unknown'" json:"algorithmUsed"`
	TotalDistanceKm  float64   `gorm:"type:numeric(8,3);not null"     json:"totalDistanceKm"`
	TotalStops       int16     `gorm:"not null"                       json:"totalStops"`
	ComputedAt       time.Time `gorm:"not null;default:now()"         json:"computedAt"`

	CourierBatch CourierBatch `gorm:"foreignKey:CourierBatchID" json:"-"`
	Stops        []RouteStop  `gorm:"foreignKey:RouteID"        json:"-"`
}

func (Route) TableName() string { return "route" }

// Request

// ComputeRouteRequest is sent by the solver service
// after a CourierBatch reaches status 'pending_route'.
// Not a public client endpoint — called internally after batch assignment.
type ComputeRouteRequest struct {
	CourierBatchID  uuid.UUID `json:"courierBatchId"  binding:"required"`
	AlgorithmUsed   string    `json:"algorithmUsed"   binding:"required"`
	TotalDistanceKm float64   `json:"totalDistanceKm" binding:"required,gt=0"`
	Stops           []RouteStopInput `json:"stops" binding:"required,min=1,dive"`
}

// RouteStopInput is a single stop as produced by the solver.
type RouteStopInput struct {
	PackageID uuid.UUID `json:"packageId" binding:"required"`
	StopOrder int16     `json:"stopOrder" binding:"required,min=1"`
	Lat       float64   `json:"lat"       binding:"required"`
	Lng       float64   `json:"lng"       binding:"required"`
}

// Response

type RouteResponse struct {
	ID              uuid.UUID         `json:"id"`
	CourierBatchID  uuid.UUID         `json:"courierBatchId"`
	AlgorithmUsed   string            `json:"algorithmUsed"`
	TotalDistanceKm float64           `json:"totalDistanceKm"`
	TotalStops      int16             `json:"totalStops"`
	ComputedAt      time.Time         `json:"computedAt"`
	Stops           []RouteStopResponse `json:"stops,omitempty"`
}

func NewRouteResponse(r *Route, stops []RouteStopResponse) RouteResponse {
	return RouteResponse{
		ID:              r.ID,
		CourierBatchID:  r.CourierBatchID,
		AlgorithmUsed:   r.AlgorithmUsed,
		TotalDistanceKm: r.TotalDistanceKm,
		TotalStops:      r.TotalStops,
		ComputedAt:      r.ComputedAt,
		Stops:           stops,
	}
}
