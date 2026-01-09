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

func TestChannelService_Create(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	req := &models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Test Channel",
		AccountIdentifier: "+1234567890",
	}

	channel, err := service.Create(req)
	require.NoError(t, err)
	assert.Equal(t, "Test Channel", channel.Name)
	assert.Equal(t, models.PlatformWhatsApp, channel.Platform)
	assert.Equal(t, int64(1), channel.OrganizationID)
	assert.Equal(t, models.ChannelStatusPending, channel.Status)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "channel.created", emitter.EmittedEvents[0].EventType)
}

func TestChannelService_Create_RepoError(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	repo.CreateError = errors.New("database error")
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	_, err := service.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Test Channel",
		AccountIdentifier: "+1234567890",
	})
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestChannelService_GetByID(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	// Create a channel first
	created, _ := repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformTelegram,
		Name:              "Telegram Channel",
		AccountIdentifier: "@testbot",
	})

	channel, err := service.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Telegram Channel", channel.Name)
}

func TestChannelService_GetByID_NotFound(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	channel, err := service.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, channel)
}

func TestChannelService_ListByOrganization(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	// Create channels for different orgs
	repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "WhatsApp Channel",
		AccountIdentifier: "+1111111111",
	})
	repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformTelegram,
		Name:              "Telegram Channel",
		AccountIdentifier: "@bot1",
	})
	repo.Create(&models.CreateChannelRequest{
		OrganizationID:    2,
		Platform:          models.PlatformInstagram,
		Name:              "Instagram Channel",
		AccountIdentifier: "@insta",
	})

	channels, err := service.ListByOrganization(1, 10, 0)
	require.NoError(t, err)
	assert.Len(t, channels, 2)
}

func TestChannelService_ListByOrganization_DefaultLimit(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Test Channel",
		AccountIdentifier: "+1111111111",
	})

	// Test with invalid limit (should default to 20)
	channels, err := service.ListByOrganization(1, 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, channels)

	// Test with limit > 100 (should default to 20)
	channels, err = service.ListByOrganization(1, 200, 0)
	require.NoError(t, err)
	assert.NotNil(t, channels)
}

func TestChannelService_Update(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	// Create a channel first
	created, _ := repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Original Name",
		AccountIdentifier: "+1111111111",
	})

	newName := "Updated Name"
	err := service.Update(created.ID, &models.UpdateChannelRequest{
		Name: &newName,
	})
	require.NoError(t, err)

	// Verify the update
	channel, _ := repo.GetByID(created.ID)
	assert.Equal(t, "Updated Name", channel.Name)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "channel.updated", emitter.EmittedEvents[0].EventType)
}

func TestChannelService_Update_RepoError(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	newName := "Updated Name"
	err := service.Update(1, &models.UpdateChannelRequest{
		Name: &newName,
	})
	assert.Error(t, err)
}

func TestChannelService_UpdateStatus(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	// Create a channel first
	created, _ := repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Test Channel",
		AccountIdentifier: "+1111111111",
	})

	err := service.UpdateStatus(created.ID, models.ChannelStatusActive)
	require.NoError(t, err)

	// Verify status update
	channel, _ := repo.GetByID(created.ID)
	assert.Equal(t, models.ChannelStatusActive, channel.Status)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "channel.updated", emitter.EmittedEvents[0].EventType)
}

func TestChannelService_UpdateStatus_RepoError(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	err := service.UpdateStatus(1, models.ChannelStatusActive)
	assert.Error(t, err)
}

func TestChannelService_Delete(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	// Create a channel first
	created, _ := repo.Create(&models.CreateChannelRequest{
		OrganizationID:    1,
		Platform:          models.PlatformWhatsApp,
		Name:              "Test Channel",
		AccountIdentifier: "+1111111111",
	})

	err := service.Delete(created.ID)
	require.NoError(t, err)

	// Verify deletion
	channel, _ := repo.GetByID(created.ID)
	assert.Nil(t, channel)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "channel.deleted", emitter.EmittedEvents[0].EventType)
}

func TestChannelService_Delete_RepoError(t *testing.T) {
	repo := testutils.NewMockChannelRepository()
	repo.DeleteError = errors.New("delete failed")
	emitter := testutils.NewMockEmitter()
	service := NewChannelService(repo, emitter)

	err := service.Delete(1)
	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())
}
