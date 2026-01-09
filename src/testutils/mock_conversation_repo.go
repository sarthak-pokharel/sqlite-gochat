package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockConversationRepository is a mock implementation of ConversationRepository
type MockConversationRepository struct {
	Conversations map[int64]*models.Conversation
	NextID        int64
	CreateError   error
	GetError      error
	ListError     error
	UpdateError   error
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{
		Conversations: make(map[int64]*models.Conversation),
		NextID:        1,
	}
}

func (m *MockConversationRepository) Create(req *models.CreateConversationRequest) (*models.Conversation, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	conv := &models.Conversation{
		ID:             m.NextID,
		ChannelID:      req.ChannelID,
		ExternalUserID: req.ExternalUserID,
		Status:         models.ConversationStatusOpen,
		Priority:       req.Priority,
	}
	m.Conversations[conv.ID] = conv
	m.NextID++
	return conv, nil
}

func (m *MockConversationRepository) GetByID(id int64) (*models.Conversation, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	conv, ok := m.Conversations[id]
	if !ok {
		return nil, nil
	}
	return conv, nil
}

func (m *MockConversationRepository) GetOrCreateByUser(channelID, externalUserID int64) (*models.Conversation, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	for _, conv := range m.Conversations {
		if conv.ChannelID == channelID && conv.ExternalUserID == externalUserID {
			return conv, nil
		}
	}
	return m.Create(&models.CreateConversationRequest{
		ChannelID:      channelID,
		ExternalUserID: externalUserID,
		Priority:       models.PriorityNormal,
	})
}

func (m *MockConversationRepository) List(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Conversation, 0)
	for _, conv := range m.Conversations {
		if conv.ChannelID == channelID {
			if status == nil || conv.Status == *status {
				result = append(result, conv)
			}
		}
	}
	return result, nil
}

func (m *MockConversationRepository) Update(id int64, req *models.UpdateConversationRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	conv, ok := m.Conversations[id]
	if !ok {
		return nil
	}
	if req.Status != nil {
		conv.Status = *req.Status
	}
	if req.Priority != nil {
		conv.Priority = *req.Priority
	}
	if req.AssignedToExternalID != nil {
		conv.AssignedToExternalID = req.AssignedToExternalID
	}
	return nil
}

func (m *MockConversationRepository) UpdateLastMessage(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}
