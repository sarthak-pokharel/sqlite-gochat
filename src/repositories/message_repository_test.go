package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "John")
	conv := testutils.CreateTestConversation(t, db, channel.ID, user.ID)

	t.Run("create message successfully", func(t *testing.T) {
		msg := &models.Message{
			ConversationID: conv.ID,
			SenderType:     models.SenderExternal,
			SenderID:       &user.ID,
			Content:        "Hello, World!",
			MessageType:    models.MessageTypeText,
			Direction:      models.DirectionInbound,
			Status:         models.MessageStatusReceived,
		}

		created, err := repo.Create(msg)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "Hello, World!", created.Content)
	})
}

func TestMessageRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "John")
	conv := testutils.CreateTestConversation(t, db, channel.ID, user.ID)

	msg := &models.Message{
		ConversationID: conv.ID,
		SenderType:     models.SenderExternal,
		Content:        "Test message",
		MessageType:    models.MessageTypeText,
		Direction:      models.DirectionInbound,
		Status:         models.MessageStatusReceived,
	}
	created, _ := repo.Create(msg)

	t.Run("get existing message", func(t *testing.T) {
		found, err := repo.GetByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, "Test message", found.Content)
	})

	t.Run("return error for non-existent message", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestMessageRepository_ListByConversation(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "John")
	conv := testutils.CreateTestConversation(t, db, channel.ID, user.ID)

	// Create multiple messages
	for i := 0; i < 10; i++ {
		msg := &models.Message{
			ConversationID: conv.ID,
			SenderType:     models.SenderExternal,
			Content:        "Message " + string(rune('0'+i)),
			MessageType:    models.MessageTypeText,
			Direction:      models.DirectionInbound,
			Status:         models.MessageStatusReceived,
		}
		repo.Create(msg)
	}

	t.Run("list messages with limit", func(t *testing.T) {
		messages, err := repo.ListByConversation(conv.ID, 5, 0, nil)
		require.NoError(t, err)
		assert.Len(t, messages, 5)
	})

	t.Run("list messages with offset", func(t *testing.T) {
		messages, err := repo.ListByConversation(conv.ID, 10, 3, nil)
		require.NoError(t, err)
		assert.Len(t, messages, 7) // 10 total - 3 offset
	})

	t.Run("list messages before cursor", func(t *testing.T) {
		// Get all messages first
		allMessages, _ := repo.ListByConversation(conv.ID, 10, 0, nil)
		if len(allMessages) > 5 {
			beforeID := allMessages[5].ID
			messages, err := repo.ListByConversation(conv.ID, 10, 0, &beforeID)
			require.NoError(t, err)
			// Should only get messages with ID < beforeID
			for _, msg := range messages {
				assert.Less(t, msg.ID, beforeID)
			}
		}
	})
}

func TestMessageRepository_UpdateStatus(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "John")
	conv := testutils.CreateTestConversation(t, db, channel.ID, user.ID)

	msg := &models.Message{
		ConversationID: conv.ID,
		SenderType:     models.SenderSystem,
		Content:        "Outgoing message",
		MessageType:    models.MessageTypeText,
		Direction:      models.DirectionOutbound,
		Status:         models.MessageStatusSent,
	}
	created, _ := repo.Create(msg)

	t.Run("update message status to delivered", func(t *testing.T) {
		err := repo.UpdateStatus(created.ID, models.MessageStatusDelivered)
		require.NoError(t, err)

		found, _ := repo.GetByID(created.ID)
		assert.Equal(t, models.MessageStatusDelivered, found.Status)
	})

	t.Run("update message status to read", func(t *testing.T) {
		err := repo.UpdateStatus(created.ID, models.MessageStatusRead)
		require.NoError(t, err)

		found, _ := repo.GetByID(created.ID)
		assert.Equal(t, models.MessageStatusRead, found.Status)
	})
}
