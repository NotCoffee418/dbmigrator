package dbmigrator

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type dbTestDef struct {
	queries *MigrationQueryDefinition
	driver  string
	connStr string
}

var dbTestDefinitions = []dbTestDef{
	{
		queries: SQLite,
		driver:  "sqlite3",
		connStr: "file:test.db?cache=shared&mode=memory",
	},
	{
		queries: PostgreSQL,
		driver:  "postgres",
		connStr: "host=localhost port=10000 user=test password=test dbname=test sslmode=disable",
	},
	{
		queries: SQLServer,
		driver:  "sqlserver",
		connStr: "sqlserver://sa:Test$123@localhost:10002?database=master",
	},
	{
		queries: MySQL,
		driver:  "mysql",
		connStr: "test:test@tcp(localhost:10001)/test",
	},
}

func TestAllQueriesOnAllDatabases(t *testing.T) {
	DockerComposeDown()
	DockerComposeUp()
	defer DockerComposeDown()
	fmt.Println("This test may fail without -tags=long due to mysql startup time")
	//time.Sleep(30 * time.Second)
	for _, def := range dbTestDefinitions {
		t.Run(fmt.Sprintf("Testing %s", def.driver), func(t *testing.T) {
			var db *sql.DB = nil
			var err error = nil

			// Wait for container to be ready
			for {
				if db == nil {
					db, err = sql.Open(def.driver, def.connStr)
					if err != nil {
						time.Sleep(1 * time.Second)
						continue
					}
				}

				err = db.Ping()
				if err == nil {
					break
				}

				time.Sleep(1 * time.Second)
			}
			defer db.Close()

			// CheckTableExists
			var exists bool
			err = db.QueryRow(def.queries.CheckTableExists).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check table existence: %s\n", err)
			}
			if exists {
				t.Fatalf("Table already exists")
			}

			// CreateMigrationsTable
			_, err = db.Exec(def.queries.CreateMigrationsTable)
			if err != nil {
				t.Fatalf("Failed to create table: %s\n", err)
			}

			// Validate table creation
			err = db.QueryRow(def.queries.CheckTableExists).Scan(&exists)
			if err != nil || !exists {
				t.Fatalf("Table creation failed or table doesn't exist")
			}

			// InsertMigration
			now := time.Now()
			_, err = db.Exec(def.queries.InsertMigration, 100, now)
			if err != nil {
				t.Fatalf("Failed to insert migration: %s\n", err)
			}

			// Validate migration insertion
			var version int
			err = db.QueryRow(def.queries.SelectInstalledVersion).Scan(&version)
			if err != nil || version != 100 {
				t.Fatalf("Migration insertion failed or version mismatch")
			}

			// DeleteMigration
			_, err = db.Exec(def.queries.DeleteMigration, 100)
			if err != nil {
				t.Fatalf("Failed to delete migration: %s\n", err)
			}

			// Validate migration deletion
			err = db.QueryRow(def.queries.SelectInstalledVersion).Scan(&version)
			if err != sql.ErrNoRows {
				t.Fatalf("Migration deletion failed or version still exists")
			}
		})
	}
}

// RunDockerCompose runs the docker-compose file and waits for all services to be ready
func DockerComposeUp() {
	if IsComposeUp() {
		fmt.Println("Docker compose is already up. Restarting it")
		DockerComposeDown()
		time.Sleep(3 * time.Second)
		return
	}

	// Run docker-compose up
	cmd := exec.Command("docker", "compose", "-f", "docker-compose.integration-tests.yaml", "up", "-d")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		panic(err)
	}
}

func IsComposeUp() bool {
	cmd := exec.Command("docker", "compose", "-f", "docker-compose.integration-tests.yaml", "ps", "--services", "--filter", "status=running")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}

	runningServices := strings.Split(strings.TrimSpace(out.String()), "\n")
	requiredServices := []string{
		"dbmigrator-postgres",
		"dbmigrator-mysql",
		"dbmigrator-mssql"}

	// Check if all required services are running
	for _, required := range requiredServices {
		found := false
		for _, running := range runningServices {
			if required == running {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// ShutDownDockerCompose stops and removes all services defined in the docker-compose file
func DockerComposeDown() {
	cmd := exec.Command("docker", "compose", "-f", "docker-compose.integration-tests.yaml", "down")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to shut down services: %s\n", err)
	} else {
		fmt.Println("Successfully shut down all services")
	}
}
