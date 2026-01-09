package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes and returns a GORM database connection
func InitDB(dbPath string) (*gorm.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys via raw SQL
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	fmt.Println("Database connected successfully")
	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	err := db.AutoMigrate(
		&models.Organization{},
		&models.ChatChannel{},
		&models.ExternalUser{},
		&models.Conversation{},
		&models.Message{},
		&models.WebhookEvent{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	fmt.Println("Auto-migration completed successfully")
	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
