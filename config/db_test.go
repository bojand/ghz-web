package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConfig_GetDialect(t *testing.T) {
	var tests = []struct {
		name     string
		in       *DBConfig
		expected string
	}{
		{"type=sqlite", &DBConfig{Type: "sqlite"}, "sqlite3"},
		{"type=sqlite3", &DBConfig{Type: "sqlite3"}, "sqlite3"},
		{"type=mysql", &DBConfig{Type: "mysql"}, "mysql"},
		{"type=mssql", &DBConfig{Type: "mssql"}, "mssql"},
		{"type=postgres", &DBConfig{Type: "postgres"}, "postgres"},
		{"type=cockroachdb", &DBConfig{Type: "cockroachdb"}, "postgres"},
		{"type=unknown", &DBConfig{Type: "unknown"}, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConf := tt.in
			actual := dbConf.GetDialect()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDBConfig_Validate(t *testing.T) {
	var tests = []struct {
		name     string
		in       *DBConfig
		expected string
	}{
		{"type=unknown", &DBConfig{Type: "unknown"}, "Unsupported database type: unknown"},
		{"required path for sqlite", &DBConfig{Type: "sqlite"}, "Database.Path: is required"},
		{"required path for sqlite3", &DBConfig{Type: "sqlite3"}, "Database.Path: is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.in.Validate()
			assert.Equal(t, tt.expected, actual.Error())
		})
	}
}

func TestDBConfig_GetConnectionString(t *testing.T) {
	var tests = []struct {
		name     string
		in       *DBConfig
		expected string
	}{
		{"type=unknown", &DBConfig{Type: "unknown", Path: "test.db"}, ""},
		{"sqlite", &DBConfig{Type: "sqlite", Path: "test.db"}, "test.db"},
		{"sqlite3", &DBConfig{Type: "sqlite3", Path: "test2.db"}, "test2.db"},
		{"postgres", &DBConfig{Type: "postgres", Host: "dbhost", Name: "ghz"}, "host=dbhost dbname=ghz"},
		{"postgres 2",
			&DBConfig{Type: "postgres", Host: "dbhost", Name: "ghz", User: "dbuser", Password: "dbpwd", SSLMode: "disabled"},
			"host=dbhost user=dbuser dbname=ghz sslmode=disabled password=dbpwd"},
		{"cockroachdb",
			&DBConfig{Type: "cockroachdb", Host: "dbhost", Name: "ghz", User: "dbuser", Password: "dbpwd", SSLMode: "disabled", Port: 3210},
			"host=dbhost port=3210 user=dbuser dbname=ghz sslmode=disabled password=dbpwd"},
		{"mysql",
			&DBConfig{Type: "mysql", Host: "dbhost", Name: "ghz"},
			"tcp(dbhost)/ghz?parseTime=true"},
		{"mysql 2",
			&DBConfig{Type: "mysql", Host: "dbhost", Port: 5555, Name: "ghz", User: "dbuser", Password: "dbpass"},
			"dbuser:dbpass@tcp(dbhost:5555)/ghz?parseTime=true"},
		{"mysql 127.0.0.1",
			&DBConfig{Type: "mysql", Host: "127.0.0.1", Port: 3306, Name: "ghz", User: "dbuser", Password: "dbpass"},
			"dbuser:dbpass@tcp(127.0.0.1:3306)/ghz?parseTime=true"},
		{"mssql",
			&DBConfig{Type: "mssql", Host: "localhost", Port: 3333, Name: "ghz", User: "dbuser", Password: "dbpass"},
			"sqlserver://dbuser:dbpass@localhost:3333?database=ghz"},
		{"mssql 2",
			&DBConfig{Type: "mssql", Host: "localhost", Name: "ghz", User: "dbuser"},
			"sqlserver://dbuser@localhost?database=ghz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.in.GetConnectionString()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
