package domain

import (
	"time"

	"github.com/google/uuid"
)

// Model
type BatchAssignment struct {
	Base
	LockerID       uuid.UUID `gorm:"type:uuid;not null;index"   json:"lockerId"`
	CourierID      uuid.UUID `gorm:"type:uuid;not null;index"   json:"courierId"`
	AssignedByID   uuid.UUID `gorm:"type:uuid;not null;column:assigned_by" json:"assignedById"`
	BatchDate      time.Time `gorm:"type:date;not null;default:current_date" json:"batchDate"`
	BatchSequence  int16     `gorm:"not null;default:1"         json:"batchSequence"`
	AssignedAt     time.Time `gorm:"not null;default:now()"     json:"assignedAt"`

	Locker       Locker       `gorm:"foreignKey:LockerID"     json:"-"`
	Courier      User         `gorm:"foreignKey:CourierID"    json:"-"`
	AssignedBy   User         `gorm:"foreignKey:AssignedByID" json:"-"`
}

func (BatchAssignment) TableName() string { return "batch_assignment" }

// Request
type AssignLockerToCourierRequest struct {
	LockerID      uuid.UUID `json:"lockerId"      binding:"required"`
	CourierID     uuid.UUID `json:"courierId"     binding:"required"`
	BatchSequence int16     `json:"batchSequence" binding:"required,min=1"`
}

// Response
type BatchAssignmentResponse struct {
	ID            uuid.UUID `json:"id"`
	LockerID      uuid.UUID `json:"lockerId"`
	CourierID     uuid.UUID `json:"courierId"`
	AssignedByID  uuid.UUID `json:"assignedById"`
	BatchDate     time.Time `json:"batchDate"`
	BatchSequence int16     `json:"batchSequence"`
	AssignedAt    time.Time `json:"assignedAt"`
}

func NewBatchAssignmentResponse(ba *BatchAssignment) BatchAssignmentResponse {
	return BatchAssignmentResponse{
		ID:            ba.ID,
		LockerID:      ba.LockerID,
		CourierID:     ba.CourierID,
		AssignedByID:  ba.AssignedByID,
		BatchDate:     ba.BatchDate,
		BatchSequence: ba.BatchSequence,
		AssignedAt:    ba.AssignedAt,
	}
}
