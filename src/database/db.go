package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the SQLite database connection
func InitDB(dbPath string) error {
	// Create database directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	DB = db
	fmt.Println("Database connected successfully")
	return nil
}

// RunMigrations executes SQL migration files
func RunMigrations(migrationsPath string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Read migration file
	migrationFile := filepath.Join(migrationsPath, "001_initial_schema.sql")
	sqlContent, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := DB.Exec(string(sqlContent)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	fmt.Println("Migration executed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
