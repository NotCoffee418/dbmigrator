# Go Database Migrator

Simple database migration system for Go.

## Usage

### Expected migration file structure

- Must contain a `-- +up` comment to indicate the SQL below should run when applying a migration.
- May contain a `-- +down` comment to indicate the SQL below should run when reverting a migration.
- Comments behind `-- +up` and `-- +down` are allowed.

```sql
-- +up  <- SQL below runs when applying a migration
CREATE TABLE demo_guestbook (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    message text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

-- +down <- SQL below runs when reverting a migration
DROP TABLE demo_guestbook;
```

### Expected project structure

Your migration files must be named in the format `0001_initial_migration.sql` where `0001` is the migration number and `initial_migration` is the name of the migration.

```md
|-- main.go
|-- migrations
|   |-- 0001_initial_migration.sql
|   |-- 0002_second_migration.sql
```

### Apply and Revert Migrations

```go
import (
    "github.com/NotCoffee418/dbmigrator"
)

//go:embed all:migrations
var migrationFS embed.FS // `embed.FS` recommended but any `fs.FS` will work

func main() {
    // Directory to migrations inside the FS
    migrationsDir := "migrations"

    // Migrations CLI (optional)
    if len(os.Args > 1) {
        // help           - Display this help message.
        // migrate up     - Apply all new database migrations.
        // migrate down   - Rollback a single database migration.`
        dbmigrator.HandleMigratorCommand(
            db *sql.DB, 
            migrationFS embed.FS,
            migrationsDir string, // Path to migrations dir in your fs 
            os.Args[1:] ...string)


    // Manage migrations programatically
    } else {
        // Apply all new migrations
        doneUp := <-dbmigrator.MigrateUpCh(
            db *sql.DB,
            migrations embed.FS,
            migrationsDir string)

        // Revert Single Migration
        doneDown := <-dbmigrator.MigrateDownCh(
            db *sql.DB,
            migrations embed.FS,
            migrationsDir string)
    }
}
```
