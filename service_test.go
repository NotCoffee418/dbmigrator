package dbmigrator

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getInstalledMigrationVersion(t *testing.T) {
	db, mock := GetMockDB()

	// Mock sanity check
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// Mock for successful query
	rows := sqlmock.NewRows([]string{"version"}).AddRow(5)
	mock.ExpectQuery("SELECT version FROM migrations ORDER BY version DESC LIMIT 1").WillReturnRows(rows)

	// Channel to collect the result
	resultChan := getInstalledMigrationVersionCh(db)
	version := <-resultChan

	// Validate
	assert.Equal(t, 5, version)
}

// GetMockDB returns a mock DB for testing
func GetMockDB() (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return mockDB, mock
}
