package domain

import "github.com/google/uuid"

// Model
type Hub struct {
	Base
	Name string `gorm:"type:varchar(300);not null" json:"name"`
	City string `gorm:"type:varchar(300);not null" json:"city"`

	Users   []User   `gorm:"foreignKey:HubID" json:"-"`
	Lockers []Locker `gorm:"foreignKey:HubID" json:"-"`
}

func (Hub) TableName() string { return "hub" }

// Request
type CreateHubRequest struct {
	Name string `json:"name" binding:"required,max=300"`
	City string `json:"city" binding:"required,max=300"`
}

type UpdateHubRequest struct {
	Name  string `json:"name" binding:"omitempty,max=300"`
	City  string `json:"city" binding:"omitempty,max=300"`
}

// Response
type HubResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	City string    `json:"city"`
}

func NewHubResponse(h *Hub) HubResponse {
	return HubResponse{
		ID:   h.ID,
		Name: h.Name,
		City: h.City,
	}
}