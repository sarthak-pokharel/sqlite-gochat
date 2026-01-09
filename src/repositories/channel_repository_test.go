package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewChannelRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")

	t.Run("create channel successfully", func(t *testing.T) {
		req := &models.CreateChannelRequest{
			OrganizationID:    org.ID,
			Platform:          models.PlatformWhatsApp,
			Name:              "WhatsApp Channel",
			AccountIdentifier: "wa-123",
		}

		channel, err := repo.Create(req)
		require.NoError(t, err)
		assert.NotZero(t, channel.ID)
		assert.Equal(t, "WhatsApp Channel", channel.Name)
		assert.Equal(t, models.PlatformWhatsApp, channel.Platform)
	})
}

func TestChannelRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewChannelRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformTelegram, "Telegram Channel")

	t.Run("get existing channel", func(t *testing.T) {
		found, err := repo.GetByID(channel.ID)
		require.NoError(t, err)
		assert.Equal(t, channel.ID, found.ID)
		assert.Equal(t, "Telegram Channel", found.Name)
	})

	t.Run("return error for non-existent channel", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestChannelRepository_ListByOrganization(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewChannelRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")

	// Create multiple channels
	testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA 1")
	testutils.CreateTestChannel(t, db, org.ID, models.PlatformTelegram, "TG 1")
	testutils.CreateTestChannel(t, db, org.ID, models.PlatformFacebook, "FB 1")

	t.Run("list channels with limit", func(t *testing.T) {
		channels, err := repo.ListByOrganization(org.ID, 2, 0)
		require.NoError(t, err)
		assert.Len(t, channels, 2)
	})

	t.Run("list channels with offset", func(t *testing.T) {
		channels, err := repo.ListByOrganization(org.ID, 10, 1)
		require.NoError(t, err)
		assert.Len(t, channels, 2) // 3 total - 1 offset = 2
	})
}

func TestChannelRepository_UpdateStatus(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewChannelRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA Test")

	t.Run("update channel status", func(t *testing.T) {
		err := repo.UpdateStatus(channel.ID, models.ChannelStatusInactive)
		require.NoError(t, err)

		found, _ := repo.GetByID(channel.ID)
		assert.Equal(t, models.ChannelStatusInactive, found.Status)
	})
}

func TestChannelRepository_Delete(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewChannelRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "Delete Me")

	t.Run("soft delete channel", func(t *testing.T) {
		err := repo.Delete(channel.ID)
		require.NoError(t, err)

		found, _ := repo.GetByID(channel.ID)
		assert.NotNil(t, found)
		assert.False(t, found.IsActive)
	})
}
