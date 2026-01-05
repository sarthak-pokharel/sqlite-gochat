package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(org *models.CreateOrganizationRequest) (*models.Organization, error)
	GetByID(id int64) (*models.Organization, error)
	GetBySlug(slug string) (*models.Organization, error)
	List(limit, offset int) ([]*models.Organization, error)
	Update(id int64, req *models.UpdateOrganizationRequest) error
	Delete(id int64) error
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(req *models.CreateOrganizationRequest) (*models.Organization, error) {
	org := &models.Organization{
		Name:     req.Name,
		Slug:     req.Slug,
		Metadata: req.Metadata,
		IsActive: true,
	}
	
	if err := r.db.Create(org).Error; err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	
	return org, nil
}

func (r *organizationRepository) GetByID(id int64) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.First(&org, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(slug string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.Where("slug = ?", slug).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) List(limit, offset int) ([]*models.Organization, error) {
	var orgs []*models.Organization
	if err := r.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	return orgs, nil
}

func (r *organizationRepository) Update(id int64, req *models.UpdateOrganizationRequest) error {
	updates := make(map[string]interface{})
	
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}
	
	result := r.db.Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	
	return nil
}

func (r *organizationRepository) Delete(id int64) error {
	result := r.db.Model(&models.Organization{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	
	return nil
}
