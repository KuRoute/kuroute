package domain

import (
	"time"

	"github.com/google/uuid"
)

// Enum
type DeliveryResult string

const (
	DeliveryResultSuccess              DeliveryResult = "success"
	DeliveryResultFailedRecipientAbsent DeliveryResult = "failed_recipient_absent"
	DeliveryResultFailedAddressNotFound DeliveryResult = "failed_address_not_found"
	DeliveryResultFailedRejected        DeliveryResult = "failed_rejected"
	DeliveryResultFailedOther           DeliveryResult = "failed_other"
)

// IsFailure returns true for any non-success result.
// Use this instead of string comparison throughout the codebase.
func (r DeliveryResult) IsFailure() bool {
	return r != DeliveryResultSuccess
}

// Model

// DeliveryReport logs each delivery attempt at a RouteStop.
// One-to-many: a stop can have multiple attempts (e.g. retry next day).
// Never overwrite — always insert a new row per attempt.
type DeliveryReport struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RouteStopID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"routeStopId"`
	Result        DeliveryResult `gorm:"type:delivery_result;not null" json:"result"`
	FailureReason *string        `gorm:"type:text"                json:"failureReason"` // nil when result=success
	PhotoURL      *string        `gorm:"type:text"                json:"photoUrl"`
	SignatureURL  *string        `gorm:"type:text"                json:"signatureUrl"`
	Notes         *string        `gorm:"type:text"                json:"notes"`
	ReportedAt    time.Time      `gorm:"not null;default:now()"   json:"reportedAt"`

	RouteStop RouteStop `gorm:"foreignKey:RouteStopID" json:"-"`
}

func (DeliveryReport) TableName() string { return "delivery_report" }

// Request

// SubmitDeliveryReportRequest is sent by the courier
// after attempting delivery at a stop.
// routeStopId comes from the URL param, not the body.
// reportedById is resolved from JWT claims.
type SubmitDeliveryReportRequest struct {
	Result        DeliveryResult `json:"result"        binding:"required,oneof=success failed_recipient_absent failed_address_not_found failed_rejected failed_other"`
	FailureReason *string        `json:"failureReason" binding:"omitempty"`
	PhotoURL      *string        `json:"photoUrl"      binding:"omitempty,url"`
	SignatureURL  *string        `json:"signatureUrl"  binding:"omitempty,url"`
	Notes         *string        `json:"notes"         binding:"omitempty"`
}

// Validate enforces the business rule that mirrors the DB CHECK constraint:
// failure_reason must be nil on success, and is allowed on failure.
// Call this in the service layer before persisting.
func (r *SubmitDeliveryReportRequest) Validate() error {
	if r.Result == DeliveryResultSuccess && r.FailureReason != nil {
		return ErrFailureReasonOnSuccess
	}
	return nil
}

var ErrFailureReasonOnSuccess = domainError("failure_reason must be nil when result is success")

type domainError string

func (e domainError) Error() string { return string(e) }

// Response
type DeliveryReportResponse struct {
	ID            uuid.UUID      `json:"id"`
	RouteStopID   uuid.UUID      `json:"routeStopId"`
	Result        DeliveryResult `json:"result"`
	FailureReason *string        `json:"failureReason"`
	PhotoURL      *string        `json:"photoUrl"`
	SignatureURL  *string        `json:"signatureUrl"`
	Notes         *string        `json:"notes"`
	ReportedAt    time.Time      `json:"reportedAt"`
}

func NewDeliveryReportResponse(dr *DeliveryReport) DeliveryReportResponse {
	return DeliveryReportResponse{
		ID:            dr.ID,
		RouteStopID:   dr.RouteStopID,
		Result:        dr.Result,
		FailureReason: dr.FailureReason,
		PhotoURL:      dr.PhotoURL,
		SignatureURL:  dr.SignatureURL,
		Notes:         dr.Notes,
		ReportedAt:    dr.ReportedAt,
	}
}
