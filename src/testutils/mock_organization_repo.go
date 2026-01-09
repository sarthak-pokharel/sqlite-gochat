package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockOrganizationRepository is a mock implementation of OrganizationRepository
type MockOrganizationRepository struct {
	Organizations       map[int64]*models.Organization
	OrganizationsBySlug map[string]*models.Organization
	NextID              int64
	CreateError         error
	GetError            error
	ListError           error
	UpdateError         error
	DeleteError         error
}

func NewMockOrganizationRepository() *MockOrganizationRepository {
	return &MockOrganizationRepository{
		Organizations:       make(map[int64]*models.Organization),
		OrganizationsBySlug: make(map[string]*models.Organization),
		NextID:              1,
	}
}

func (m *MockOrganizationRepository) Create(req *models.CreateOrganizationRequest) (*models.Organization, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	org := &models.Organization{
		ID:   m.NextID,
		Name: req.Name,
		Slug: req.Slug,
	}
	m.Organizations[org.ID] = org
	m.OrganizationsBySlug[org.Slug] = org
	m.NextID++
	return org, nil
}

func (m *MockOrganizationRepository) GetByID(id int64) (*models.Organization, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	org, ok := m.Organizations[id]
	if !ok {
		return nil, nil
	}
	return org, nil
}

func (m *MockOrganizationRepository) GetBySlug(slug string) (*models.Organization, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	org, ok := m.OrganizationsBySlug[slug]
	if !ok {
		return nil, nil
	}
	return org, nil
}

func (m *MockOrganizationRepository) List(limit, offset int) ([]*models.Organization, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Organization, 0)
	for _, org := range m.Organizations {
		result = append(result, org)
	}
	return result, nil
}

func (m *MockOrganizationRepository) Update(id int64, req *models.UpdateOrganizationRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	org, ok := m.Organizations[id]
	if !ok {
		return nil
	}
	if req.Name != nil {
		org.Name = *req.Name
	}
	return nil
}

func (m *MockOrganizationRepository) Delete(id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	org, ok := m.Organizations[id]
	if ok {
		delete(m.OrganizationsBySlug, org.Slug)
		delete(m.Organizations, id)
	}
	return nil
}
