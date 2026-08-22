package domain

import (
	"time"

	"github.com/google/uuid"
)

// Model

// LockerClusterAssignment represents the result of clustering:
// which package is assigned to which physical locker.
//
// This is a logical assignment created by the clustering service.
// It does not mean the package has physically entered the locker yet.
// Physical placement is recorded separately by LockerScan.
type LockerClusterAssignment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LockerID    uuid.UUID `gorm:"type:uuid;not null;index"                json:"lockerId"`
	PackageID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"          json:"packageId"`
	ClusterArea string    `gorm:"type:varchar(200);not null"              json:"clusterArea"`
	ClusteredAt time.Time `gorm:"not null;default:now()"                  json:"clusteredAt"`

	Locker  Locker  `gorm:"foreignKey:LockerID"  json:"-"`
	Package Package `gorm:"foreignKey:PackageID" json:"-"`
}

func (LockerClusterAssignment) TableName() string {
	return "locker_cluster_assignment"
}

// Response

// LockerClusterAssignmentResponse represents the locker destination
// for a package before the package is physically scanned into the locker.
type LockerClusterAssignmentResponse struct {
	ID            uuid.UUID `json:"id"`
	PackageID     uuid.UUID `json:"packageId"`
	TrackingCode  string    `json:"trackingCode"`
	RecipientName string    `json:"recipientName"`
	LockerID      uuid.UUID `json:"lockerId"`
	LockerLabel   string    `json:"lockerLabel"`
	ClusterArea   string    `json:"clusterArea"`
	ClusteredAt   time.Time `json:"clusteredAt"`
}

func NewLockerClusterAssignmentResponse(assignment *LockerClusterAssignment) LockerClusterAssignmentResponse {
	return LockerClusterAssignmentResponse{
		ID:            assignment.ID,
		PackageID:     assignment.PackageID,
		TrackingCode:  assignment.Package.TrackingCode,
		RecipientName: assignment.Package.RecipientName,
		LockerID:      assignment.LockerID,
		LockerLabel:   assignment.Locker.Label,
		ClusterArea:   assignment.Locker.ClusterArea,
		ClusteredAt:   assignment.ClusteredAt,
	}
}