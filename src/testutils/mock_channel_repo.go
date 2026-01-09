package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockChannelRepository is a mock implementation of ChannelRepository
type MockChannelRepository struct {
	Channels    map[int64]*models.ChatChannel
	NextID      int64
	CreateError error
	GetError    error
	ListError   error
	UpdateError error
	DeleteError error
}

func NewMockChannelRepository() *MockChannelRepository {
	return &MockChannelRepository{
		Channels: make(map[int64]*models.ChatChannel),
		NextID:   1,
	}
}

func (m *MockChannelRepository) Create(req *models.CreateChannelRequest) (*models.ChatChannel, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	channel := &models.ChatChannel{
		ID:                m.NextID,
		OrganizationID:    req.OrganizationID,
		Platform:          req.Platform,
		Name:              req.Name,
		AccountIdentifier: req.AccountIdentifier,
		Status:            models.ChannelStatusPending,
		IsActive:          true,
	}
	m.Channels[channel.ID] = channel
	m.NextID++
	return channel, nil
}

func (m *MockChannelRepository) GetByID(id int64) (*models.ChatChannel, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	channel, ok := m.Channels[id]
	if !ok {
		return nil, nil
	}
	return channel, nil
}

func (m *MockChannelRepository) ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.ChatChannel, 0)
	for _, ch := range m.Channels {
		if ch.OrganizationID == orgID {
			result = append(result, ch)
		}
	}
	return result, nil
}

func (m *MockChannelRepository) Update(id int64, req *models.UpdateChannelRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	channel, ok := m.Channels[id]
	if !ok {
		return nil
	}
	if req.Name != nil {
		channel.Name = *req.Name
	}
	if req.Status != nil {
		channel.Status = *req.Status
	}
	return nil
}

func (m *MockChannelRepository) UpdateStatus(id int64, status models.ChannelStatus) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	channel, ok := m.Channels[id]
	if ok {
		channel.Status = status
	}
	return nil
}

func (m *MockChannelRepository) Delete(id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Channels, id)
	return nil
}
