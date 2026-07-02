package domain

import (
	"time"

	"github.com/google/uuid"
)

// Model
type LockerScan struct {
	Base
	LockerID    uuid.UUID `gorm:"type:uuid;not null;index"                   json:"lockerId"`
	PackageID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"             json:"packageId"` // 1 package → 1 locker
	ScannedByID uuid.UUID `gorm:"type:uuid;not null;column:scanned_by"       json:"scannedById"`
	ScannedAt   time.Time `gorm:"not null;default:now()"                     json:"scannedAt"`

	Locker    Locker  `gorm:"foreignKey:LockerID"    json:"-"`
	Package   Package `gorm:"foreignKey:PackageID"   json:"-"`
	ScannedBy User    `gorm:"foreignKey:ScannedByID" json:"-"`
}

func (LockerScan) TableName() string { return "locker_scan" }

// Request

// ScanPackageToLockerRequest is the payload when staff
// physically scans a package barcode and assigns it to a locker.
// scannedById is resolved from JWT claims, not from request body.
type ScanPackageToLockerRequest struct {
	LockerID  uuid.UUID `json:"lockerId"  binding:"required"`
	PackageID uuid.UUID `json:"packageId" binding:"required"`
}

// Response
type LockerScanResponse struct {
	ID          uuid.UUID `json:"id"`
	LockerID    uuid.UUID `json:"lockerId"`
	PackageID   uuid.UUID `json:"packageId"`
	ScannedByID uuid.UUID `json:"scannedById"`
	ScannedAt   time.Time `json:"scannedAt"`
}

func NewLockerScanResponse(ls *LockerScan) LockerScanResponse {
	return LockerScanResponse{
		ID:          ls.ID,
		LockerID:    ls.LockerID,
		PackageID:   ls.PackageID,
		ScannedByID: ls.ScannedByID,
		ScannedAt:   ls.ScannedAt,
	}
}
