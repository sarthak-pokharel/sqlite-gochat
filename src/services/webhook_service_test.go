package services

import (
	"errors"
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMessageService is a local mock for MessageService to avoid import cycles
type mockMessageService struct {
	ProcessedMessages  []*ProcessIncomingMessageRequest
	SentMessages       []*SendOutgoingMessageRequest
	DeliveredMessages  []int64
	ReadMessages       []int64
	ProcessError       error
	SendError          error
	MarkDeliveredError error
	MarkReadError      error
	ReturnMessage      *models.Message
}

func newMockMessageService() *mockMessageService {
	return &mockMessageService{
		ProcessedMessages: make([]*ProcessIncomingMessageRequest, 0),
		SentMessages:      make([]*SendOutgoingMessageRequest, 0),
		DeliveredMessages: make([]int64, 0),
		ReadMessages:      make([]int64, 0),
		ReturnMessage: &models.Message{
			ID:      1,
			Content: "test message",
		},
	}
}

func (m *mockMessageService) ProcessIncomingMessage(req *ProcessIncomingMessageRequest) (*models.Message, error) {
	if m.ProcessError != nil {
		return nil, m.ProcessError
	}
	m.ProcessedMessages = append(m.ProcessedMessages, req)
	return m.ReturnMessage, nil
}

func (m *mockMessageService) SendOutgoingMessage(req *SendOutgoingMessageRequest) (*models.Message, error) {
	if m.SendError != nil {
		return nil, m.SendError
	}
	m.SentMessages = append(m.SentMessages, req)
	return m.ReturnMessage, nil
}

func (m *mockMessageService) GetMessageHistory(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	return nil, nil
}

func (m *mockMessageService) MarkDelivered(messageID int64) error {
	if m.MarkDeliveredError != nil {
		return m.MarkDeliveredError
	}
	m.DeliveredMessages = append(m.DeliveredMessages, messageID)
	return nil
}

func (m *mockMessageService) MarkRead(messageID int64) error {
	if m.MarkReadError != nil {
		return m.MarkReadError
	}
	m.ReadMessages = append(m.ReadMessages, messageID)
	return nil
}

func TestWebhookService_ProcessWebhook_MessageEvent(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id":   "msg-123",
		"user_id":      "user-456",
		"user_name":    "John Doe",
		"content":      "Hello from webhook!",
		"message_type": "text",
	}

	err := service.ProcessWebhook(1, "message", payload)
	require.NoError(t, err)

	// Verify webhook event was stored
	assert.Len(t, eventRepo.Events, 1)

	// Verify message was processed
	assert.Len(t, msgService.ProcessedMessages, 1)
	assert.Equal(t, "user-456", msgService.ProcessedMessages[0].PlatformUserID)
	assert.Equal(t, "Hello from webhook!", msgService.ProcessedMessages[0].Content)
}

func TestWebhookService_ProcessWebhook_StatusUpdateDelivered(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id": float64(123),
		"status":     "delivered",
	}

	err := service.ProcessWebhook(1, "status_update", payload)
	require.NoError(t, err)

	// Verify webhook event was stored
	assert.Len(t, eventRepo.Events, 1)

	// Verify mark delivered was called
	assert.Len(t, msgService.DeliveredMessages, 1)
	assert.Equal(t, int64(123), msgService.DeliveredMessages[0])
}

func TestWebhookService_ProcessWebhook_StatusUpdateRead(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id": float64(123),
		"status":     "read",
	}

	err := service.ProcessWebhook(1, "status_update", payload)
	require.NoError(t, err)

	// Verify mark read was called
	assert.Len(t, msgService.ReadMessages, 1)
	assert.Equal(t, int64(123), msgService.ReadMessages[0])
}

func TestWebhookService_ProcessWebhook_UnknownEventType(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"data": "test",
	}

	err := service.ProcessWebhook(1, "unknown_event", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")

	// Verify event was marked as failed
	assert.Len(t, eventRepo.Events, 1)
}

func TestWebhookService_ProcessWebhook_EventRepoError(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	eventRepo.CreateError = errors.New("database error")
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"data": "test",
	}

	err := service.ProcessWebhook(1, "message", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store webhook event")
}

func TestWebhookService_ProcessWebhook_MessageProcessError(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	msgService.ProcessError = errors.New("processing failed")
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id": "msg-123",
		"user_id":    "user-456",
		"content":    "Hello!",
	}

	err := service.ProcessWebhook(1, "message", payload)
	assert.Error(t, err)

	// Verify event was marked as failed
	event := eventRepo.Events[int64(1)]
	assert.NotNil(t, event.Error)
}

func TestWebhookService_ProcessWebhook_InvalidPayloadFormat(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	// Pass a non-map payload
	payload := "invalid"

	err := service.ProcessWebhook(1, "message", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload format")
}

func TestWebhookService_ProcessWebhook_MissingRequiredFields(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	// Payload missing user_id and content
	payload := map[string]interface{}{
		"message_id": "msg-123",
	}

	err := service.ProcessWebhook(1, "message", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestWebhookService_ProcessWebhook_StatusUpdateMissingMessageID(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"status": "delivered",
	}

	err := service.ProcessWebhook(1, "status_update", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing message_id")
}

func TestWebhookService_ProcessWebhook_MarkDeliveredError(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	msgService.MarkDeliveredError = errors.New("mark delivered failed")
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id": float64(123),
		"status":     "delivered",
	}

	err := service.ProcessWebhook(1, "status_update", payload)
	assert.Error(t, err)
}

func TestWebhookService_ProcessWebhook_MarkReadError(t *testing.T) {
	eventRepo := testutils.NewMockWebhookEventRepository()
	msgService := newMockMessageService()
	msgService.MarkReadError = errors.New("mark read failed")
	service := NewWebhookService(eventRepo, msgService)

	payload := map[string]interface{}{
		"message_id": float64(123),
		"status":     "read",
	}

	err := service.ProcessWebhook(1, "status_update", payload)
	assert.Error(t, err)
}
