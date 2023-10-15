package dbmigrator

import "time"

type migrationFileInfo struct {
	version  int
	file     string
	contents *migrationContents // not always populated
}

type migrationContents struct {
	up   string
	down string
}

type MigrationsTable struct {
	Version     int       `db:"version"`
	InstalledAt time.Time `db:"installed_at"`
}

type MigrationState struct {
	AvailableVersion int
	InstalledVersion int
	Migrations       []migrationFileInfo
}

// MigrationQueries describes the queries used by the migrator.
// These can be overridden if you want to use a different DB or table name.
type MigrationQueryDefinition struct {
	CheckTableExists       string // Expect booly result
	CreateMigrationsTable  string
	InsertMigration        string
	DeleteMigration        string
	SelectInstalledVersion string
}
