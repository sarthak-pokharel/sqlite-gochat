package services

import (
	"github/sarthak-pokharel/sqlite-d1-gochat/src/events"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
)

// OrganizationService handles organization business logic
type OrganizationService interface {
	Create(req *models.CreateOrganizationRequest) (*models.Organization, error)
	GetByID(id int64) (*models.Organization, error)
	GetBySlug(slug string) (*models.Organization, error)
	List(limit, offset int) ([]*models.Organization, error)
	Update(id int64, req *models.UpdateOrganizationRequest) error
	Delete(id int64) error
}

type organizationService struct {
	repo    repositories.OrganizationRepository
	emitter events.Emitter
}

func NewOrganizationService(repo repositories.OrganizationRepository, emitter events.Emitter) OrganizationService {
	return &organizationService{
		repo:    repo,
		emitter: emitter,
	}
}

func (s *organizationService) Create(req *models.CreateOrganizationRequest) (*models.Organization, error) {
	org, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Emit event (non-blocking, fire and forget)
	go s.emitter.Emit("organization.created", map[string]interface{}{
		"organization_id": org.ID,
		"slug":            org.Slug,
		"name":            org.Name,
	})

	return org, nil
}

func (s *organizationService) GetByID(id int64) (*models.Organization, error) {
	return s.repo.GetByID(id)
}

func (s *organizationService) GetBySlug(slug string) (*models.Organization, error) {
	return s.repo.GetBySlug(slug)
}

func (s *organizationService) List(limit, offset int) ([]*models.Organization, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.List(limit, offset)
}

func (s *organizationService) Update(id int64, req *models.UpdateOrganizationRequest) error {
	if err := s.repo.Update(id, req); err != nil {
		return err
	}

	go s.emitter.Emit("organization.updated", map[string]interface{}{
		"organization_id": id,
	})

	return nil
}

func (s *organizationService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	go s.emitter.Emit("organization.deleted", map[string]interface{}{
		"organization_id": id,
	})

	return nil
}
