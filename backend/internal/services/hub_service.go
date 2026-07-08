package services

import (
	"database/sql"
	"errors"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/internal/repository"
	"github.com/google/uuid"
)

type HubService struct {
	hubRepository *repository.HubRepository
}

func NewHubService(hubRepository *repository.HubRepository) *HubService {
	return &HubService{
		hubRepository: hubRepository,
	}
}

func (s *HubService) CreateHub(req *domain.CreateHubRequest) (*domain.HubResponse, error) {
	if req.Name == "" || req.City == "" {
		return nil, errors.New("name and city are required")
	}

	hub := &domain.Hub{
		Name: req.Name,
		City: req.City,
	}

	if err := s.hubRepository.CreateHub(hub); err != nil {
		return nil, err
	}

	resp := domain.NewHubResponse(hub)
	return &resp, nil
}

func (s *HubService) GetHubByID(id uuid.UUID) (*domain.HubResponse, error) {
	hub, err := s.hubRepository.GetHubByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("hub not found")
		}
		return nil, err
	}
	resp := domain.NewHubResponse(hub)
	return &resp, nil
}

func (s *HubService) ListHubs() ([]domain.HubResponse, error) {
	hubs, err := s.hubRepository.ListHubs()
	if err != nil {
		return nil, err
	}

	resp := make([]domain.HubResponse, 0, len(hubs))
	for _, h := range hubs {
		resp = append(resp, domain.NewHubResponse(&h))
	}
	return resp, nil
}

func (s *HubService) UpdateHub(id uuid.UUID, req *domain.UpdateHubRequest) (*domain.HubResponse, error) {
	hub, err := s.hubRepository.GetHubByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("hub not found")
		}
		return nil, err
	}

	if req.Name != "" {
		hub.Name = req.Name
	}
	if req.City != "" {
		hub.City = req.City
	}

	if err := s.hubRepository.UpdateHub(hub); err != nil {
		return nil, err
	}

	resp := domain.NewHubResponse(hub)
	return &resp, nil
}

func (s *HubService) DeleteHub(id uuid.UUID) error {
	return s.hubRepository.DeleteHub(id)
}
