package domain

import (
	"time"

	"github.com/google/uuid"
)

// Model

// Locker represents a physical slot at a hub that holds
// packages belonging to the same geographic cluster.
// Lockers are created by the clustering algorithm, not manually.
type Locker struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	HubID       uuid.UUID `gorm:"type:uuid;not null;index"           json:"hubId"`
	Label       string    `gorm:"type:varchar(50);not null"          json:"label"`
	ClusterArea string    `gorm:"type:varchar(200);not null"         json:"clusterArea"`
	ClusteredAt time.Time `gorm:"not null;default:now()"             json:"clusteredAt"`

	Hub      Hub       `gorm:"foreignKey:HubID"    json:"-"`
	Packages []Package `gorm:"many2many:locker_scan;joinForeignKey:LockerID;joinReferences:PackageID" json:"-"`
}

func (Locker) TableName() string { return "locker" }

// Request

// CreateLockerRequest is called by the clustering service,
// not directly by a human operator.
type CreateLockerRequest struct {
	HubID       uuid.UUID `json:"hubId"       binding:"required"`
	Label       string    `json:"label"       binding:"required,max=50"`
	ClusterArea string    `json:"clusterArea" binding:"required,max=200"`
}

type UpdateLockerRequest struct {
	Label       *string `json:"label"       binding:"omitempty,max=50"`
	ClusterArea *string `json:"clusterArea" binding:"omitempty,max=200"`
}

// Response
type LockerResponse struct {
	ID          uuid.UUID `json:"id"`
	HubID       uuid.UUID `json:"hubId"`
	Label       string    `json:"label"`
	ClusterArea string    `json:"clusterArea"`
	ClusteredAt time.Time `json:"clusteredAt"`
}

func NewLockerResponse(l *Locker) LockerResponse {
	return LockerResponse{
		ID:          l.ID,
		HubID:       l.HubID,
		Label:       l.Label,
		ClusterArea: l.ClusterArea,
		ClusteredAt: l.ClusteredAt,
	}
}
