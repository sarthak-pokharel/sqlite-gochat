# Migration Tool

Auto-generate database migrations from GORM models.

## Usage

### Generate a new migration
```bash
go run tools/migrate/main.go generate <migration_name>
```

Example:
```bash
go run tools/migrate/main.go generate add_user_preferences
```

This will create a timestamped migration file in `migrations/` directory based on current GORM models.

### Run pending migrations
```bash
go run tools/migrate/main.go run
```

### Check migration status
```bash
go run tools/migrate/main.go status
```

## How it works

1. Reads GORM model definitions from `src/models/`
2. Uses GORM's migrator to introspect schema
3. Generates SQL DDL statements
4. Creates timestamped migration files

## Migration Files

Generated files follow the format:
```
YYYYMMDDHHMMSS_migration_name.sql
```

Example: `20260105143000_add_user_preferences.sql`
