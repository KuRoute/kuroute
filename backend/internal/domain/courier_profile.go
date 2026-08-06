package domain

import "github.com/google/uuid"

// Enum

type VehicleType string

const (
	VehicleTypeMotor VehicleType = "motor"
	VehicleTypeMobil VehicleType = "mobil"
	VehicleTypeVan   VehicleType = "van"
)

// Model
type CourierProfile struct {
	UserID       uuid.UUID   `gorm:"type:uuid;primaryKey"        json:"userId"`
	VehicleType  VehicleType `gorm:"type:vehicle_type;not null"  json:"vehicleType"`
	VehiclePlate string      `gorm:"type:varchar(20);not null"   json:"vehiclePlate"`

	User User `gorm:"foreignKey:UserID"           json:"-"`
}

func (CourierProfile) TableName() string {
	return "courier_profile"
}

// Request
type CreateCourierProfileRequest struct {
	VehicleType  VehicleType `json:"vehicleType" binding:"required,oneof=motor mobil van"`
	VehiclePlate string      `json:"vehiclePlate" binding:"required,max=20"`
}

type UpdateCourierProfileRequest struct {
	VehicleType  *VehicleType `json:"vehicleType" binding:"omitempty,oneof=motor mobil van"`
	VehiclePlate *string      `json:"vehiclePlate" binding:"omitempty,max=20"`
}

// Response
type CourierProfileResponse struct {
	UserID       uuid.UUID   `json:"userId"`
	VehicleType  VehicleType `json:"vehicleType"`
	VehiclePlate string      `json:"vehiclePlate"`
}

func NewCourierProfileResponse(cp *CourierProfile) CourierProfileResponse {
	return CourierProfileResponse{
		UserID:       cp.UserID,
		VehicleType:  cp.VehicleType,
		VehiclePlate: cp.VehiclePlate,
	}
}
