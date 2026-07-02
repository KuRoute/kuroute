package domain

import "github.com/google/uuid"

// Enum

type VehicleType string

const (
	VehicleTypeMotor  VehicleType = "motor"
	VehicleTypeMobil  VehicleType = "mobil"
	VehicleTypeVan    VehicleType = "van"
)

// Model
type CourierProfile struct {
	UserID      uuid.UUID   `gorm:"type:uuid;primaryKey"        json:"userId"`
	VehicleType	VehicleType `gorm:"type:vehicle_type;not null"  json:"vehicleType"`

	User		User		`gorm:"foreignKey:UserID"           json:"-"`
}

func (CourierProfile) TableName() string {
	return "courier_profile"
}

// Request
type CreateCourierProfileRequest struct {
	VehicleType VehicleType `json:"vehicleType" binding:"required,oneof=motor mobil van"`
}

type UpdateCourierProfileRequest struct {
	VehicleType *VehicleType `json:"vehicleType" binding:"omitempty,oneof=motor mobil van"`
}

// Response
type CourierProfileResponse struct {
	UserID      uuid.UUID   `json:"userId"`
	VehicleType VehicleType `json:"vehicleType"`
}

func NewCourierProfileResponse(cp *CourierProfile) CourierProfileResponse {
	return CourierProfileResponse{
		UserID:      cp.UserID,
		VehicleType: cp.VehicleType,
	}
}
