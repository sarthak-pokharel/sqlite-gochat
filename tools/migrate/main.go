package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/config"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/database"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/migrate/main.go <command> [args]")
		fmt.Println("\nCommands:")
		fmt.Println("  generate <name>  - Generate a new migration file")
		fmt.Println("  run             - Run pending migrations")
		fmt.Println("  status          - Show migration status")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	command := os.Args[1]

	switch command {
	case "generate":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a migration name")
		}
		generateMigration(os.Args[2])
	case "run":
		runMigrations()
	case "status":
		showStatus()
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func generateMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := filepath.Join("migrations", filename)

	// Get current schema using GORM migrator
	statements := generateDDL()

	content := fmt.Sprintf("-- Migration: %s\n-- Generated: %s\n\n%s",
		name,
		time.Now().Format(time.RFC3339),
		statements,
	)

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write migration file: %v", err)
	}

	fmt.Printf("Generated migration: %s\n", filepath)
}

func generateDDL() string {
	db := database.DB
	migrator := db.Migrator()

	ddl := ""

	// Generate CREATE TABLE statements for all models
	models := []interface{}{
		&models.Organization{},
		&models.ChatChannel{},
		&models.ExternalUser{},
		&models.Conversation{},
		&models.Message{},
		&models.WebhookEvent{},
	}

	for _, model := range models {
		tableName := db.Statement.Table
		db.Statement.Parse(model)
		tableName = db.Statement.Table

		if !migrator.HasTable(model) {
			// Generate CREATE TABLE
			ddl += fmt.Sprintf("-- Table: %s\n", tableName)
			ddl += generateCreateTable(db, model)
			ddl += "\n\n"
		}
	}

	// Generate indexes
	ddl += generateIndexes(db)

	return ddl
}

func generateCreateTable(db *gorm.DB, model interface{}) string {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)

	ddl := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", stmt.Table)

	fields := []string{}
	for _, field := range stmt.Schema.Fields {
		fieldDef := fmt.Sprintf("    %s %s", field.DBName, getSQLType(field))

		if field.PrimaryKey {
			fieldDef += " PRIMARY KEY"
			if field.AutoIncrement {
				fieldDef += " AUTOINCREMENT"
			}
		}

		if field.NotNull && !field.PrimaryKey {
			fieldDef += " NOT NULL"
		}

		if field.Unique {
			fieldDef += " UNIQUE"
		}

		if field.DefaultValue != "" {
			fieldDef += fmt.Sprintf(" DEFAULT %s", field.DefaultValue)
		}

		fields = append(fields, fieldDef)
	}

	// Add foreign keys
	for _, rel := range stmt.Schema.Relationships.Relations {
		if rel.Type == "belongs_to" || rel.Type == "has_one" {
			for _, ref := range rel.References {
				fk := fmt.Sprintf("    FOREIGN KEY (%s) REFERENCES %s(%s)",
					ref.ForeignKey.DBName,
					rel.FieldSchema.Table,
					ref.PrimaryKey.DBName,
				)
				fields = append(fields, fk)
			}
		}
	}

	ddl += joinFields(fields, ",\n")
	ddl += "\n);"

	return ddl
}

func getSQLType(field *schema.Field) string {
	switch field.DataType {
	case "string", "text":
		if field.Size > 0 {
			return fmt.Sprintf("TEXT(%d)", field.Size)
		}
		return "TEXT"
	case "int", "int64":
		return "INTEGER"
	case "bool":
		return "INTEGER"
	case "time":
		return "DATETIME"
	default:
		return "TEXT"
	}
}

func generateIndexes(db *gorm.DB) string {
	ddl := "-- Indexes\n"

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);",
		"CREATE INDEX IF NOT EXISTS idx_chat_channels_org ON chat_channels(organization_id);",
		"CREATE INDEX IF NOT EXISTS idx_chat_channels_platform ON chat_channels(platform);",
		"CREATE INDEX IF NOT EXISTS idx_external_users_channel ON external_users(channel_id);",
		"CREATE INDEX IF NOT EXISTS idx_conversations_channel ON conversations(channel_id);",
		"CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);",
		"CREATE INDEX IF NOT EXISTS idx_webhook_events_channel ON webhook_events(channel_id);",
	}

	for _, idx := range indexes {
		ddl += idx + "\n"
	}

	return ddl
}

func joinFields(fields []string, sep string) string {
	result := ""
	for i, field := range fields {
		result += field
		if i < len(fields)-1 {
			result += sep
		}
	}
	return result
}

func runMigrations() {
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("Migrations completed successfully")
}

func showStatus() {
	db := database.DB

	models := []interface{}{
		&models.Organization{},
		&models.ChatChannel{},
		&models.ExternalUser{},
		&models.Conversation{},
		&models.Message{},
		&models.WebhookEvent{},
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, model := range models {
		db.Statement.Parse(model)
		tableName := db.Statement.Table

		if db.Migrator().HasTable(model) {
			fmt.Printf("✓ %s - exists\n", tableName)
		} else {
			fmt.Printf("✗ %s - missing\n", tableName)
		}
	}
}
