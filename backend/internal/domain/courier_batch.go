package domain

import (
	"time"

	"github.com/google/uuid"
)

// Enum
type BatchStatus string

const (
	BatchStatusPendingRoute BatchStatus = "pending_route" // assignment dibuat, solver belum jalan
	BatchStatusRouteReady   BatchStatus = "route_ready"   // rute sudah dihitung, kurir belum berangkat
	BatchStatusInProgress   BatchStatus = "in_progress"   // kurir sedang di jalan
	BatchStatusCompleted    BatchStatus = "completed"      // semua stop selesai
)

// Model
type CourierBatch struct {
	ID            		uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BatchAssignmentID 	uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex" json:"batchAssignmentId"`
	Status            	BatchStatus `gorm:"type:batch_status;not null;default:'pending_route'" json:"status"`
	StartedAt         	*time.Time  `gorm:"default:null"                   json:"startedAt"`   // nil until courier taps "mulai"
	CompletedAt       	*time.Time  `gorm:"default:null"                   json:"completedAt"` // nil until all stops done

	BatchAssignment BatchAssignment `gorm:"foreignKey:BatchAssignmentID" json:"-"`
}

func (CourierBatch) TableName() string { return "courier_batch" }

// Request
type StartBatchRequest struct{}

// Response
type CourierBatchResponse struct {
	ID                uuid.UUID   `json:"id"`
	BatchAssignmentID uuid.UUID   `json:"batchAssignmentId"`
	Status            BatchStatus `json:"status"`
	StartedAt         *time.Time  `json:"startedAt"`
	CompletedAt       *time.Time  `json:"completedAt"`
}

func NewCourierBatchResponse(cb *CourierBatch) CourierBatchResponse {
	return CourierBatchResponse{
		ID:                cb.ID,
		BatchAssignmentID: cb.BatchAssignmentID,
		Status:            cb.Status,
		StartedAt:         cb.StartedAt,
		CompletedAt:       cb.CompletedAt,
	}
}
