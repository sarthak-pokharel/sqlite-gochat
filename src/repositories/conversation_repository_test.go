package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversationRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewConversationRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "John Doe")

	t.Run("create conversation successfully", func(t *testing.T) {
		req := &models.CreateConversationRequest{
			ChannelID:      channel.ID,
			ExternalUserID: user.ID,
		}

		conv, err := repo.Create(req)
		require.NoError(t, err)
		assert.NotZero(t, conv.ID)
		assert.Equal(t, channel.ID, conv.ChannelID)
		assert.Equal(t, user.ID, conv.ExternalUserID)
		assert.Equal(t, models.ConversationStatusOpen, conv.Status)
	})
}

func TestConversationRepository_GetOrCreateByUser(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewConversationRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-456", "Jane Doe")

	t.Run("create new conversation if not exists", func(t *testing.T) {
		conv, err := repo.GetOrCreateByUser(channel.ID, user.ID)
		require.NoError(t, err)
		assert.NotZero(t, conv.ID)
	})

	t.Run("return existing conversation", func(t *testing.T) {
		conv1, _ := repo.GetOrCreateByUser(channel.ID, user.ID)
		conv2, err := repo.GetOrCreateByUser(channel.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, conv1.ID, conv2.ID)
	})
}

func TestConversationRepository_List(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewConversationRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	// Create users and conversations
	for i := 0; i < 5; i++ {
		user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-"+string(rune('a'+i)), "User "+string(rune('A'+i)))
		testutils.CreateTestConversation(t, db, channel.ID, user.ID)
	}

	t.Run("list all conversations", func(t *testing.T) {
		conversations, err := repo.List(channel.ID, nil, 10, 0)
		require.NoError(t, err)
		assert.Len(t, conversations, 5)
	})

	t.Run("filter by status", func(t *testing.T) {
		status := models.ConversationStatusOpen
		conversations, err := repo.List(channel.ID, &status, 10, 0)
		require.NoError(t, err)
		assert.Len(t, conversations, 5)
	})
}

func TestConversationRepository_Update(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewConversationRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-789", "Test User")
	conv := testutils.CreateTestConversation(t, db, channel.ID, user.ID)

	t.Run("update conversation status and priority", func(t *testing.T) {
		newStatus := models.ConversationStatusResolved
		newPriority := models.PriorityHigh
		assignee := "agent-001"

		req := &models.UpdateConversationRequest{
			Status:               &newStatus,
			Priority:             &newPriority,
			AssignedToExternalID: &assignee,
		}

		err := repo.Update(conv.ID, req)
		require.NoError(t, err)

		found, _ := repo.GetByID(conv.ID)
		assert.Equal(t, models.ConversationStatusResolved, found.Status)
		assert.Equal(t, models.PriorityHigh, found.Priority)
		assert.Equal(t, "agent-001", *found.AssignedToExternalID)
	})
}
