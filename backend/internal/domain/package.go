package domain

import (
	"time"

	"github.com/google/uuid"
)

// Enum
type PackageStatus string

const (
	PackageStatusReceived   PackageStatus = "received"    // masuk hub, belum di-scan ke locker
	PackageStatusSorted     PackageStatus = "sorted"      // sudah masuk locker hasil clustering
	PackageStatusAssigned   PackageStatus = "assigned"    // locker sudah di-bind ke kurir
	PackageStatusInDelivery PackageStatus = "in_delivery" // kurir sedang antar
	PackageStatusDelivered  PackageStatus = "delivered"   // berhasil diterima
	PackageStatusFailed     PackageStatus = "failed"      // gagal antar (semua attempt)
)

// Model
type Package struct {
	Base
	HubID         uuid.UUID     `gorm:"type:uuid;not null;index"              json:"hubId"`
	TrackingCode  string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"trackingCode"`
	RecipientName string        `gorm:"type:varchar(200);not null"            json:"recipientName"`
	AddressText   string        `gorm:"type:text;not null"                    json:"addressText"`
	Lat           float64       `gorm:"type:numeric(9,6);not null"            json:"lat"`
	Lng           float64       `gorm:"type:numeric(9,6);not null"            json:"lng"`
	Status        PackageStatus `gorm:"type:package_status;not null;default:'received'" json:"status"`
	ReceivedAt    time.Time     `gorm:"not null;default:now()"               json:"receivedAt"`

	Hub Hub `gorm:"foreignKey:HubID" json:"-"`
}

func (Package) TableName() string { return "package" }

// Request

// CreatePackageRequest is called when a package arrives at the hub.
// lat/lng should be resolved via geocoding before this request is saved —
// do NOT rely on the client to provide coordinates; geocode server-side.
type CreatePackageRequest struct {
	TrackingCode  string  `json:"trackingCode"  binding:"required,max=100"`
	RecipientName string  `json:"recipientName" binding:"required,max=200"`
	AddressText   string  `json:"addressText"   binding:"required"`
	Lat           float64 `json:"lat"           binding:"required"`
	Lng           float64 `json:"lng"           binding:"required"`
}

// UpdatePackageStatusRequest is used internally by the system
// (route solver, delivery report handler) to advance package state machine.
// Not exposed as a public endpoint — status transitions are side effects
// of other operations (scan, assign, deliver).
type UpdatePackageStatusRequest struct {
	Status PackageStatus `json:"status" binding:"required,oneof=received sorted assigned in_delivery delivered failed"`
}

// Response
type PackageResponse struct {
	ID            uuid.UUID     `json:"id"`
	HubID         uuid.UUID     `json:"hubId"`
	TrackingCode  string        `json:"trackingCode"`
	RecipientName string        `json:"recipientName"`
	AddressText   string        `json:"addressText"`
	Lat           float64       `json:"lat"`
	Lng           float64       `json:"lng"`
	Status        PackageStatus `json:"status"`
	ReceivedAt    time.Time     `json:"receivedAt"`
}

func NewPackageResponse(p *Package) PackageResponse {
	return PackageResponse{
		ID:            p.ID,
		HubID:         p.HubID,
		TrackingCode:  p.TrackingCode,
		RecipientName: p.RecipientName,
		AddressText:   p.AddressText,
		Lat:           p.Lat,
		Lng:           p.Lng,
		Status:        p.Status,
		ReceivedAt:    p.ReceivedAt,
	}
}
