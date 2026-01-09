package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookEventRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewWebhookEventRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	t.Run("create webhook event successfully", func(t *testing.T) {
		event := &models.WebhookEvent{
			ChannelID: channel.ID,
			EventType: "message",
			Payload:   `{"type": "text", "content": "Hello"}`,
			Processed: false,
		}

		created, err := repo.Create(event)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
	})
}

func TestWebhookEventRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewWebhookEventRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	event := &models.WebhookEvent{
		ChannelID: channel.ID,
		EventType: "message",
		Payload:   `{"content": "test"}`,
	}
	repo.Create(event)

	t.Run("get existing event", func(t *testing.T) {
		found, err := repo.GetByID(event.ID)
		require.NoError(t, err)
		assert.Equal(t, "message", found.EventType)
	})

	t.Run("return error for non-existent event", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestWebhookEventRepository_MarkProcessed(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewWebhookEventRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	event := &models.WebhookEvent{
		ChannelID: channel.ID,
		EventType: "message",
		Payload:   `{"content": "process me"}`,
		Processed: false,
	}
	repo.Create(event)

	t.Run("mark event as processed", func(t *testing.T) {
		err := repo.MarkProcessed(event.ID)
		require.NoError(t, err)

		found, _ := repo.GetByID(event.ID)
		assert.True(t, found.Processed)
		assert.NotNil(t, found.ProcessedAt)
	})
}

func TestWebhookEventRepository_MarkFailed(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewWebhookEventRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	event := &models.WebhookEvent{
		ChannelID: channel.ID,
		EventType: "message",
		Payload:   `{"content": "fail me"}`,
		Processed: false,
	}
	repo.Create(event)

	t.Run("mark event as failed", func(t *testing.T) {
		err := repo.MarkFailed(event.ID, "processing error: invalid format")
		require.NoError(t, err)

		found, _ := repo.GetByID(event.ID)
		assert.True(t, found.Processed)
		assert.NotNil(t, found.Error)
		assert.Equal(t, "processing error: invalid format", *found.Error)
	})
}

func TestWebhookEventRepository_GetUnprocessed(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewWebhookEventRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	// Create processed and unprocessed events
	for i := 0; i < 3; i++ {
		event := &models.WebhookEvent{
			ChannelID: channel.ID,
			EventType: "message",
			Payload:   `{}`,
			Processed: false,
		}
		repo.Create(event)
	}

	processedEvent := &models.WebhookEvent{
		ChannelID: channel.ID,
		EventType: "message",
		Payload:   `{}`,
		Processed: true,
	}
	repo.Create(processedEvent)

	t.Run("get unprocessed events", func(t *testing.T) {
		events, err := repo.ListUnprocessed(channel.ID, 10)
		require.NoError(t, err)
		assert.Len(t, events, 3)

		for _, e := range events {
			assert.False(t, e.Processed)
		}
	})
}
