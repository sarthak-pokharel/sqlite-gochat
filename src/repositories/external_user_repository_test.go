package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalUserRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewExternalUserRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	t.Run("create external user successfully", func(t *testing.T) {
		displayName := "John Doe"
		phone := "+1234567890"
		req := &models.CreateExternalUserRequest{
			ChannelID:      channel.ID,
			PlatformUserID: "wa-user-123",
			DisplayName:    &displayName,
			PhoneNumber:    &phone,
		}

		user, err := repo.Create(req)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})
}

func TestExternalUserRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewExternalUserRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "user-123", "Jane Doe")

	t.Run("get existing user", func(t *testing.T) {
		found, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, "Jane Doe", *found.DisplayName)
	})

	t.Run("return error for non-existent user", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestExternalUserRepository_FindOrCreate(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewExternalUserRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")

	t.Run("create new user if not exists", func(t *testing.T) {
		req := &models.CreateExternalUserRequest{
			ChannelID:      channel.ID,
			PlatformUserID: "new-platform-user-id",
		}
		user, err := repo.FindOrCreate(req)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "new-platform-user-id", user.PlatformUserID)
	})

	t.Run("return existing user", func(t *testing.T) {
		req := &models.CreateExternalUserRequest{
			ChannelID:      channel.ID,
			PlatformUserID: "existing-user",
		}
		user1, _ := repo.FindOrCreate(req)
		user2, err := repo.FindOrCreate(req)
		require.NoError(t, err)
		assert.Equal(t, user1.ID, user2.ID)
	})

	t.Run("find or create with display name", func(t *testing.T) {
		displayName := "New User"
		req := &models.CreateExternalUserRequest{
			ChannelID:      channel.ID,
			PlatformUserID: "user-with-name",
			DisplayName:    &displayName,
		}
		user, err := repo.FindOrCreate(req)
		require.NoError(t, err)
		assert.Equal(t, "New User", *user.DisplayName)
	})
}

func TestExternalUserRepository_Update(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewExternalUserRepository(db)
	org := testutils.CreateTestOrganization(t, db, "Test Org", "testorg")
	channel := testutils.CreateTestChannel(t, db, org.ID, models.PlatformWhatsApp, "WA")
	user := testutils.CreateTestExternalUser(t, db, channel.ID, "update-me", "Original Name")

	t.Run("update user display name", func(t *testing.T) {
		newName := "Updated Name"
		req := &models.UpdateExternalUserRequest{
			DisplayName: &newName,
		}

		err := repo.Update(user.ID, req)
		require.NoError(t, err)

		found, _ := repo.GetByID(user.ID)
		assert.Equal(t, "Updated Name", *found.DisplayName)
	})

	t.Run("block user via update", func(t *testing.T) {
		isBlocked := true
		req := &models.UpdateExternalUserRequest{
			IsBlocked: &isBlocked,
		}

		err := repo.Update(user.ID, req)
		require.NoError(t, err)

		found, _ := repo.GetByID(user.ID)
		assert.True(t, found.IsBlocked)
	})
}
