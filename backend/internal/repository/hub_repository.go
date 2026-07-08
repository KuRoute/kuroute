package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HubRepository struct {
	db *database.DB
}

func NewHubRepository(db *database.DB) *HubRepository {
	return &HubRepository{db: db}
}

func (r *HubRepository) CreateHub(hub *domain.Hub) error {
	result := r.db.Create(hub)
	return result.Error
}

func (r *HubRepository) GetHubByID(id uuid.UUID) (*domain.Hub, error) {
	var hub domain.Hub

	result := r.db.First(&hub, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &hub, result.Error
}

func (r *HubRepository) UpdateHub(hub *domain.Hub) error {
	result := r.db.Model(&hub).Updates(hub)
	return result.Error
}

func (r *HubRepository) DeleteHub(id uuid.UUID) error {
	result := r.db.Delete(&domain.Hub{}, "id = ?", id)
	return result.Error
}

func (r *HubRepository) ListHubs() ([]domain.Hub, error) {
	var hubs []domain.Hub
	result := r.db.Find(&hubs)
	return hubs, result.Error
}
