package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Enum
type UserRole string

const (
	UserRoleStaffSortir UserRole = "staff_sortir"
	UserRoleKurir       UserRole = "kurir"
)

// Model
type User struct {
	BaseWithUpdate
	HubID        uuid.UUID `gorm:"type:uuid;not null"                       json:"hubId"`
	Name         string    `gorm:"type:varchar(200);not null"               json:"name"`
	Phone        string    `gorm:"type:varchar(20);not null;uniqueIndex"    json:"phone"`
	Role         UserRole  `gorm:"type:user_role;not null;default:'kurir'"  json:"role"`
	PasswordHash string    `gorm:"type:text;not null"                       json:"-"`

	Hub          Hub       `gorm:"foreignKey:HubID"                         json:"-"`
}

func (User) TableName() string { return "users" }

// JWT Claims
type JWTClaims struct {
	UserID uuid.UUID `json:"userId"`
	Role   UserRole  `json:"role"`
	HubID  uuid.UUID `json:"hubId"`
	jwt.RegisteredClaims
}

// Request
type RegisterUserRequest struct {
	HubID    uuid.UUID `json:"hubId"    binding:"required"`
	Name     string    `json:"name"     binding:"required,max=200"`
	Phone    string    `json:"phone"    binding:"required,max=20"`
	Role     UserRole  `json:"role"     binding:"required,oneof=staff_sortir kurir"`
	Password string    `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"  binding:"omitempty,max=200"`
	Phone *string `json:"phone" binding:"omitempty,max=20"`
}

// Response
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	HubID uuid.UUID `json:"hubId"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Role  UserRole  `json:"role"`
}

type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user"`
}

func NewUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		HubID: u.HubID,
		Name:  u.Name,
		Phone: u.Phone,
		Role:  u.Role,
	}
}