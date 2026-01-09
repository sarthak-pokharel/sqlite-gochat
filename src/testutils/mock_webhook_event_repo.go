package testutils

import "github/sarthak-pokharel/sqlite-d1-gochat/src/models"

// MockWebhookEventRepository is a mock implementation of WebhookEventRepository
type MockWebhookEventRepository struct {
	Events      map[int64]*models.WebhookEvent
	NextID      int64
	CreateError error
	GetError    error
	ListError   error
	UpdateError error
}

func NewMockWebhookEventRepository() *MockWebhookEventRepository {
	return &MockWebhookEventRepository{
		Events: make(map[int64]*models.WebhookEvent),
		NextID: 1,
	}
}

func (m *MockWebhookEventRepository) Create(event *models.WebhookEvent) (*models.WebhookEvent, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	event.ID = m.NextID
	m.Events[event.ID] = event
	m.NextID++
	return event, nil
}

func (m *MockWebhookEventRepository) GetByID(id int64) (*models.WebhookEvent, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	event, ok := m.Events[id]
	if !ok {
		return nil, nil
	}
	return event, nil
}

func (m *MockWebhookEventRepository) ListUnprocessed(channelID int64, limit int) ([]*models.WebhookEvent, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.WebhookEvent, 0)
	for _, event := range m.Events {
		if event.ChannelID == channelID && event.ProcessedAt == nil {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *MockWebhookEventRepository) MarkProcessed(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockWebhookEventRepository) MarkFailed(id int64, errorMsg string) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	event, ok := m.Events[id]
	if ok {
		event.Error = &errorMsg
	}
	return nil
}
