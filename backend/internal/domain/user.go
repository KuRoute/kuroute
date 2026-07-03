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
	HubID              uuid.UUID `gorm:"type:uuid;not null"                       json:"hubId"`
	Name               string    `gorm:"type:varchar(200);not null"               json:"name"`
	Email              string    `gorm:"type:varchar(200);not null;uniqueIndex"   json:"email"`
	Phone              string    `gorm:"type:varchar(20);not null;uniqueIndex"    json:"phone"`
	Role               UserRole  `gorm:"type:user_role;not null;default:'kurir'"  json:"role"`
	PasswordHash       string    `gorm:"type:text;not null"                       json:"-"`

	Hub Hub `gorm:"foreignKey:HubID" json:"-"`
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
	Email    string    `json:"email"    binding:"required,email,max=200"`
	Phone    string    `json:"phone"    binding:"required,max=20"`
	Role     UserRole  `json:"role"     binding:"required,oneof=staff_sortir kurir"`
	Password string    `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email,max=200"`
	Password string `json:"password" binding:"required"`
}

type SetupPasswordRequest struct {
	SetupToken  string `json:"setupToken" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"  binding:"omitempty,max=200"`
	Email *string `json:"email" binding:"omitempty,email,max=200"`
	Phone *string `json:"phone" binding:"omitempty,max=20"`
}

// Response
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	HubID uuid.UUID `json:"hubId"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
	Role  UserRole  `json:"role"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int          `json:"expiresIn"`
	User         UserResponse `json:"user"`
}

type SetupPasswordResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int          `json:"expiresIn"`
	User         UserResponse `json:"user"`
}

type SetupTokenResponse struct {
	SetupToken string       `json:"setupToken"`
	User       UserResponse `json:"user"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func NewUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		HubID: u.HubID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Role:  u.Role,
	}
}