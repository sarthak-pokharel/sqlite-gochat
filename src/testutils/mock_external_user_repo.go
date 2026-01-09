package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockExternalUserRepository is a mock implementation of ExternalUserRepository
type MockExternalUserRepository struct {
	Users       map[int64]*models.ExternalUser
	NextID      int64
	CreateError error
	GetError    error
	UpdateError error
}

func NewMockExternalUserRepository() *MockExternalUserRepository {
	return &MockExternalUserRepository{
		Users:  make(map[int64]*models.ExternalUser),
		NextID: 1,
	}
}

func (m *MockExternalUserRepository) Create(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	user := &models.ExternalUser{
		ID:             m.NextID,
		ChannelID:      req.ChannelID,
		PlatformUserID: req.PlatformUserID,
	}
	if req.DisplayName != nil {
		user.DisplayName = req.DisplayName
	}
	m.Users[user.ID] = user
	m.NextID++
	return user, nil
}

func (m *MockExternalUserRepository) GetByID(id int64) (*models.ExternalUser, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	user, ok := m.Users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *MockExternalUserRepository) GetByPlatformUser(channelID int64, platformUserID string) (*models.ExternalUser, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	for _, user := range m.Users {
		if user.ChannelID == channelID && user.PlatformUserID == platformUserID {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockExternalUserRepository) FindOrCreate(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	user, err := m.GetByPlatformUser(req.ChannelID, req.PlatformUserID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return m.Create(req)
}

func (m *MockExternalUserRepository) Update(id int64, req *models.UpdateExternalUserRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockExternalUserRepository) UpdateLastSeen(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}
