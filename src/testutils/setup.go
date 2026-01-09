package testutils

import (
	"os"
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/database"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use in-memory SQLite for tests
	if err := database.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return database.DB
}

func TeardownTestDB(t *testing.T) {
	t.Helper()
	database.Close()
}

func SetupTestDBFile(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	if err := database.InitDB(tmpFile.Name()); err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		database.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to run migrations: %v", err)
	}

	cleanup := func() {
		database.Close()
		os.Remove(tmpFile.Name())
	}

	return database.DB, cleanup
}

func CreateTestOrganization(t *testing.T, db *gorm.DB, name, slug string) *models.Organization {
	t.Helper()

	org := &models.Organization{
		Name:     name,
		Slug:     slug,
		IsActive: true,
	}

	if err := db.Create(org).Error; err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	return org
}

func CreateTestChannel(t *testing.T, db *gorm.DB, orgID int64, platform models.Platform, name string) *models.ChatChannel {
	t.Helper()

	channel := &models.ChatChannel{
		OrganizationID:    orgID,
		Platform:          platform,
		Name:              name,
		AccountIdentifier: "test-account-" + name,
		Status:            models.ChannelStatusActive,
		IsActive:          true,
	}

	if err := db.Create(channel).Error; err != nil {
		t.Fatalf("Failed to create test channel: %v", err)
	}

	return channel
}

func CreateTestExternalUser(t *testing.T, db *gorm.DB, channelID int64, platformUserID, displayName string) *models.ExternalUser {
	t.Helper()

	user := &models.ExternalUser{
		ChannelID:      channelID,
		PlatformUserID: platformUserID,
		DisplayName:    &displayName,
		IsBlocked:      false,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test external user: %v", err)
	}

	return user
}

func CreateTestConversation(t *testing.T, db *gorm.DB, channelID, externalUserID int64) *models.Conversation {
	t.Helper()

	conv := &models.Conversation{
		ChannelID:      channelID,
		ExternalUserID: externalUserID,
		Status:         models.ConversationStatusOpen,
		Priority:       models.PriorityNormal,
	}

	if err := db.Create(conv).Error; err != nil {
		t.Fatalf("Failed to create test conversation: %v", err)
	}

	return conv
}
