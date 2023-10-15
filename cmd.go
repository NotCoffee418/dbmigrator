package dbmigrator

import (
	"database/sql"
	"fmt"
	"io/fs"
)

// HandleMigratorCommand is intended to be hooked into main.go
// to display help and migrate based on args for manual migrattions.
// This function is optional but can be used to as part of a CLI interface.
//
// Param: db - database connection using sqlx
//
// Param: migrationFS - ideally embed.FS based on a `migrations`
// with migration files structured as described in the documentation.
//
// Param: migrationDir - directory to use for migrations inside the FS.
// 'migrations' is recommended.
//
// Param: args - os.Args[1:] from main.go
//
// Returns: boolean indicating if a command was actionable
// Not indicative of success or failure.
func HandleMigratorCommand(
	db *sql.DB,
	migrationFS fs.FS,
	migrationDir string,
	args ...string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "help":
		fmt.Println(GetHelpString())
		return true
	case "migrate":
		if len(args) < 2 {
			return false
		}
		switch args[1] {
		case "up":
			<-MigrateUpCh(db, migrationFS, migrationDir)
			return true
		case "down":
			<-MigrateDownCh(db, migrationFS, migrationDir)
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func GetHelpString() string {
	return `
	migrate up     - Apply all new database migrations.
	migrate down   - Rollback a single database migration.`
}
