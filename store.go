package dbmigrator

func init() {
	// Set the default query set to MySQL
	activeQueryDef = MySQL
}

var activeQueryDef *MigrationQueryDefinition
