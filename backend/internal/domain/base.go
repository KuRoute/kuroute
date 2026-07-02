package domain

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"not null;default:now()"                         json:"createdAt"`
}

type BaseWithUpdate struct {
	Base
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}