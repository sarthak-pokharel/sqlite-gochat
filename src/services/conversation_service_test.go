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

func TestConversationService_GetByID(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create a conversation first
	created, _ := repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	conv, err := service.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, conv.ID)
	assert.Equal(t, int64(1), conv.ChannelID)
}

func TestConversationService_GetByID_NotFound(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	conv, err := service.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, conv)
}

func TestConversationService_GetByID_RepoError(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	repo.GetError = errors.New("database error")
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	_, err := service.GetByID(1)
	assert.Error(t, err)
}

func TestConversationService_ListByChannel(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create conversations
	repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})
	repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 2,
		Priority:       models.PriorityHigh,
	})
	repo.Create(&models.CreateConversationRequest{
		ChannelID:      2,
		ExternalUserID: 3,
		Priority:       models.PriorityLow,
	})

	convs, err := service.ListByChannel(1, nil, 10, 0)
	require.NoError(t, err)
	assert.Len(t, convs, 2)
}

func TestConversationService_ListByChannel_WithStatus(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create conversations (all default to Open status)
	repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	status := models.ConversationStatusOpen
	convs, err := service.ListByChannel(1, &status, 10, 0)
	require.NoError(t, err)
	assert.Len(t, convs, 1)
}

func TestConversationService_ListByChannel_DefaultLimit(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	// Test with invalid limit (should default to 20)
	convs, err := service.ListByChannel(1, nil, 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, convs)

	// Test with limit > 100 (should default to 20)
	convs, err = service.ListByChannel(1, nil, 200, 0)
	require.NoError(t, err)
	assert.NotNil(t, convs)
}

func TestConversationService_Assign(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create a conversation first
	created, _ := repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	err := service.Assign(created.ID, "agent-123")
	require.NoError(t, err)

	// Verify assignment
	conv, _ := repo.GetByID(created.ID)
	assert.NotNil(t, conv.AssignedToExternalID)
	assert.Equal(t, "agent-123", *conv.AssignedToExternalID)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "conversation.assigned", emitter.EmittedEvents[0].EventType)
	assert.Equal(t, "agent-123", emitter.EmittedEvents[0].Payload["assignee_id"])
}

func TestConversationService_Assign_RepoError(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	err := service.Assign(1, "agent-123")
	assert.Error(t, err)
}

func TestConversationService_UpdateStatus(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create a conversation first
	created, _ := repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	err := service.UpdateStatus(created.ID, models.ConversationStatusClosed)
	require.NoError(t, err)

	// Verify status update
	conv, _ := repo.GetByID(created.ID)
	assert.Equal(t, models.ConversationStatusClosed, conv.Status)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestConversationService_UpdateStatus_RepoError(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	err := service.UpdateStatus(1, models.ConversationStatusClosed)
	assert.Error(t, err)
}

func TestConversationService_UpdatePriority(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	// Create a conversation first
	created, _ := repo.Create(&models.CreateConversationRequest{
		ChannelID:      1,
		ExternalUserID: 1,
		Priority:       models.PriorityNormal,
	})

	err := service.UpdatePriority(created.ID, models.PriorityUrgent)
	require.NoError(t, err)

	// Verify priority update
	conv, _ := repo.GetByID(created.ID)
	assert.Equal(t, models.PriorityUrgent, conv.Priority)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
}

func TestConversationService_UpdatePriority_RepoError(t *testing.T) {
	repo := testutils.NewMockConversationRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewConversationService(repo, emitter)

	err := service.UpdatePriority(1, models.PriorityHigh)
	assert.Error(t, err)
}
