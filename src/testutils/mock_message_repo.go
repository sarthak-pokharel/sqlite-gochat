package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockMessageRepository is a mock implementation of MessageRepository
type MockMessageRepository struct {
	Messages    map[int64]*models.Message
	NextID      int64
	CreateError error
	GetError    error
	ListError   error
	UpdateError error
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{
		Messages: make(map[int64]*models.Message),
		NextID:   1,
	}
}

func (m *MockMessageRepository) Create(msg *models.Message) (*models.Message, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	msg.ID = m.NextID
	m.Messages[msg.ID] = msg
	m.NextID++
	return msg, nil
}

func (m *MockMessageRepository) GetByID(id int64) (*models.Message, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	msg, ok := m.Messages[id]
	if !ok {
		return nil, nil
	}
	return msg, nil
}

func (m *MockMessageRepository) ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Message, 0)
	for _, msg := range m.Messages {
		if msg.ConversationID == conversationID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (m *MockMessageRepository) UpdateStatus(id int64, status models.MessageStatus) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	msg, ok := m.Messages[id]
	if ok {
		msg.Status = status
	}
	return nil
}
