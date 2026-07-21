package repository

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackageRepository struct {
	db *database.DB
}

func NewPackageRepository(db *database.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (p *PackageRepository) CreatePackage(pack *domain.Package) error {
	result := p.db.Create(pack)

	return result.Error
}

func (p *PackageRepository) ListPackage() ([]domain.Package, error) {
	var packs []domain.Package

	result := p.db.Find(&packs)

	return packs, result.Error
}

func (p *PackageRepository) GetPackageByID(id uuid.UUID) (*domain.Package, error) {
	var pack domain.Package

	result := p.db.First(&pack, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &pack, result.Error
}

func (p *PackageRepository) GetPackageByHubID(hubId uuid.UUID) ([]domain.Package, error) {
	var packs []domain.Package

	result := p.db.Where("hub_id = ?", hubId).Find(&packs)
	return packs, result.Error
}

func (p *PackageRepository) GetStatusPackage(id uuid.UUID) (*domain.Package, error) {
	var pack domain.Package

	result := p.db.Select("status").Where("id = ?", id).First(&pack)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}

	return &pack, result.Error
}

func (p *PackageRepository) GetPackageByHubIDAndStatus(hubId uuid.UUID, status domain.PackageStatus) ([]domain.Package, error) {
	var packs []domain.Package

	result := p.db.Where("hub_id = ? AND status = ?", hubId, status).Find(&packs)
	return packs, result.Error
}

func (p *PackageRepository) UpdatePackageStatus(pack *domain.Package) error {
	result := p.db.Model(&domain.Package{}).Where("id = ?", pack.ID).Update("status", pack.Status)

	return result.Error
}

func (p *PackageRepository) DeletePackage(id uuid.UUID) error {
	result := p.db.Delete(&domain.Package{}, "id = ?", id)

	return result.Error
}
