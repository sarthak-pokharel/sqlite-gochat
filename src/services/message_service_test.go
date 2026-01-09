package services

import (
	"errors"
	"testing"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageService_ProcessIncomingMessage(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	req := &ProcessIncomingMessageRequest{
		ChannelID:         1,
		PlatformMessageID: "msg-123",
		PlatformUserID:    "user-456",
		UserDisplayName:   "John Doe",
		Content:           "Hello, world!",
		MessageType:       models.MessageTypeText,
	}

	msg, err := service.ProcessIncomingMessage(req)
	require.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "Hello, world!", msg.Content)
	assert.Equal(t, models.DirectionInbound, msg.Direction)
	assert.Equal(t, models.SenderExternal, msg.SenderType)

	// Verify user was created
	assert.Len(t, userRepo.Users, 1)

	// Verify conversation was created
	assert.Len(t, convRepo.Conversations, 1)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestMessageService_ProcessIncomingMessage_ExistingUser(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Create existing user
	displayName := "John Doe"
	userRepo.Create(&models.CreateExternalUserRequest{
		ChannelID:      1,
		PlatformUserID: "user-456",
		DisplayName:    &displayName,
	})

	req := &ProcessIncomingMessageRequest{
		ChannelID:         1,
		PlatformMessageID: "msg-123",
		PlatformUserID:    "user-456",
		UserDisplayName:   "John Doe",
		Content:           "Hello again!",
		MessageType:       models.MessageTypeText,
	}

	msg, err := service.ProcessIncomingMessage(req)
	require.NoError(t, err)
	assert.NotNil(t, msg)

	// Verify no new user was created
	assert.Len(t, userRepo.Users, 1)
}

func TestMessageService_ProcessIncomingMessage_UserRepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	userRepo.GetError = errors.New("database error")
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	req := &ProcessIncomingMessageRequest{
		ChannelID:         1,
		PlatformMessageID: "msg-123",
		PlatformUserID:    "user-456",
		UserDisplayName:   "John Doe",
		Content:           "Hello!",
		MessageType:       models.MessageTypeText,
	}

	_, err := service.ProcessIncomingMessage(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find/create user")
}

func TestMessageService_ProcessIncomingMessage_ConvRepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	convRepo.GetError = errors.New("database error")
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	req := &ProcessIncomingMessageRequest{
		ChannelID:         1,
		PlatformMessageID: "msg-123",
		PlatformUserID:    "user-456",
		UserDisplayName:   "John Doe",
		Content:           "Hello!",
		MessageType:       models.MessageTypeText,
	}

	_, err := service.ProcessIncomingMessage(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get/create conversation")
}

func TestMessageService_ProcessIncomingMessage_MsgRepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	msgRepo.CreateError = errors.New("database error")
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	req := &ProcessIncomingMessageRequest{
		ChannelID:         1,
		PlatformMessageID: "msg-123",
		PlatformUserID:    "user-456",
		UserDisplayName:   "John Doe",
		Content:           "Hello!",
		MessageType:       models.MessageTypeText,
	}

	_, err := service.ProcessIncomingMessage(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create message")
}

func TestMessageService_SendOutgoingMessage(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Create a conversation first
	conv, _ := convRepo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	req := &SendOutgoingMessageRequest{
		ConversationID: conv.ID,
		Content:        "Hello from agent!",
		MessageType:    models.MessageTypeText,
	}

	msg, err := service.SendOutgoingMessage(req)
	require.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "Hello from agent!", msg.Content)
	assert.Equal(t, models.DirectionOutbound, msg.Direction)
	assert.Equal(t, models.SenderInternal, msg.SenderType)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestMessageService_SendOutgoingMessage_ConversationNotFound(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	req := &SendOutgoingMessageRequest{
		ConversationID: 999,
		Content:        "Hello!",
		MessageType:    models.MessageTypeText,
	}

	_, err := service.SendOutgoingMessage(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestMessageService_SendOutgoingMessage_MsgRepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	msgRepo.CreateError = errors.New("database error")
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Create a conversation first
	conv, _ := convRepo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	req := &SendOutgoingMessageRequest{
		ConversationID: conv.ID,
		Content:        "Hello!",
		MessageType:    models.MessageTypeText,
	}

	_, err := service.SendOutgoingMessage(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create message")
}

func TestMessageService_GetMessageHistory(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Add some messages
	msgRepo.Create(&models.Message{
		ConversationID: 1,
		Content:        "Message 1",
		MessageType:    models.MessageTypeText,
	})
	msgRepo.Create(&models.Message{
		ConversationID: 1,
		Content:        "Message 2",
		MessageType:    models.MessageTypeText,
	})

	msgs, err := service.GetMessageHistory(1, 10, 0, nil)
	require.NoError(t, err)
	assert.Len(t, msgs, 2)
}

func TestMessageService_GetMessageHistory_DefaultLimit(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Test with invalid limit (should default to 50)
	msgs, err := service.GetMessageHistory(1, 0, 0, nil)
	require.NoError(t, err)
	assert.NotNil(t, msgs)

	// Test with limit > 100 (should default to 50)
	msgs, err = service.GetMessageHistory(1, 200, 0, nil)
	require.NoError(t, err)
	assert.NotNil(t, msgs)
}

func TestMessageService_MarkDelivered(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Create a message first
	msg, _ := msgRepo.Create(&models.Message{
		ConversationID: 1,
		Content:        "Test message",
		Status:         models.MessageStatusSent,
	})

	err := service.MarkDelivered(msg.ID)
	require.NoError(t, err)

	// Verify status update
	updated, _ := msgRepo.GetByID(msg.ID)
	assert.Equal(t, models.MessageStatusDelivered, updated.Status)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestMessageService_MarkDelivered_RepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	msgRepo.UpdateError = errors.New("update failed")
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	err := service.MarkDelivered(1)
	assert.Error(t, err)
}

func TestMessageService_MarkRead(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	// Create a message first
	msg, _ := msgRepo.Create(&models.Message{
		ConversationID: 1,
		Content:        "Test message",
		Status:         models.MessageStatusDelivered,
	})

	err := service.MarkRead(msg.ID)
	require.NoError(t, err)

	// Verify status update
	updated, _ := msgRepo.GetByID(msg.ID)
	assert.Equal(t, models.MessageStatusRead, updated.Status)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestMessageService_MarkRead_RepoError(t *testing.T) {
	msgRepo := testutils.NewMockMessageRepository()
	msgRepo.UpdateError = errors.New("update failed")
	convRepo := testutils.NewMockConversationRepository()
	userRepo := testutils.NewMockExternalUserRepository()
	emitter := testutils.NewMockEmitter()
	service := NewMessageService(msgRepo, convRepo, userRepo, emitter)

	err := service.MarkRead(1)
	assert.Error(t, err)
}
