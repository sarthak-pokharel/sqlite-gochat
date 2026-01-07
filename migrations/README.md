# Database Migrations

This directory contains database migration files for reference and documentation purposes.

## Migration Strategy

The application uses **GORM AutoMigrate** for automatic schema management. SQL files in this directory serve as:
- Documentation of the database schema
- Reference for manual database setup
- Backup for understanding schema evolution

## Files

- `001_initial_schema.sql` - Original schema design with all constraints and checks
- `002_gorm_schema.sql` - GORM-compatible schema (simplified, matches model definitions)

## GORM AutoMigrate

On startup, the application runs `database.AutoMigrate()` which automatically:
- Creates missing tables
- Adds missing columns
- Updates column types (when safe)
- Does NOT drop columns or tables

## Manual Migration

If you need to apply migrations manually:

```bash
sqlite3 data/chat.db < migrations/002_gorm_schema.sql
```

## Schema Differences

**Original SQL (`001`) vs GORM (`002`):**
- GORM uses simpler constraints (fewer CHECK constraints)
- GORM handles timestamps automatically
- GORM manages foreign keys through model associations
- Original SQL includes more strict validation

## Adding New Migrations

When models change:
1. Update GORM model structs with proper tags
2. GORM AutoMigrate handles the changes automatically
3. Optionally document changes in a new SQL file for reference
