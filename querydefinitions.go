package dbmigrator

var Postgres = &MigrationQueryDefinition{
	CheckTableExists:       "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'migrations')",
	CreateMigrationsTable:  "CREATE TABLE migrations (version INT NOT NULL, installed_at TIMESTAMP NOT NULL)",
	InsertMigration:        "INSERT INTO migrations (version, installed_at) VALUES ($1, $2)",
	DeleteMigration:        "DELETE FROM migrations WHERE version = $1",
	SelectInstalledVersion: "SELECT version FROM migrations ORDER BY version DESC LIMIT 1",
}

var MySQL = &MigrationQueryDefinition{
	CheckTableExists:       "SELECT EXISTS (SELECT * FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'migrations')",
	CreateMigrationsTable:  "CREATE TABLE migrations (version INT NOT NULL, installed_at TIMESTAMP NOT NULL)",
	InsertMigration:        "INSERT INTO migrations (version, installed_at) VALUES (?, ?)",
	DeleteMigration:        "DELETE FROM migrations WHERE version = ?",
	SelectInstalledVersion: "SELECT version FROM migrations ORDER BY version DESC LIMIT 1",
}

var SQLite = &MigrationQueryDefinition{
	CheckTableExists:       "SELECT EXISTS (SELECT name FROM sqlite_master WHERE type='table' AND name='migrations')",
	CreateMigrationsTable:  "CREATE TABLE migrations (version INT NOT NULL, installed_at TIMESTAMP NOT NULL)",
	InsertMigration:        "INSERT INTO migrations (version, installed_at) VALUES (?, ?)",
	DeleteMigration:        "DELETE FROM migrations WHERE version = ?",
	SelectInstalledVersion: "SELECT version FROM migrations ORDER BY version DESC LIMIT 1",
}

var MsSql = &MigrationQueryDefinition{
	CheckTableExists:       "SELECT CASE WHEN EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'migrations') THEN 1 ELSE 0 END",
	CreateMigrationsTable:  "CREATE TABLE migrations (version INT NOT NULL, installed_at DATETIME NOT NULL)",
	InsertMigration:        "INSERT INTO migrations (version, installed_at) VALUES (@p1, @p2)",
	DeleteMigration:        "DELETE FROM migrations WHERE version = @p1",
	SelectInstalledVersion: "SELECT TOP 1 version FROM migrations ORDER BY version DESC",
}
