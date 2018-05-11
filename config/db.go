package config

import (
	"fmt"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// DBConfig database configs
type DBConfig struct {
	// The SQL dialect we support
	// One of mysql, postgres, mssql, sqlite3
	Type string `default:"sqlite3"`

	// The database connection options
	Host     string `default:"localhost"`
	Port     uint
	Name     string `default:"ghz"`
	User     string
	Password string
	SSLMode  string `default:"disable"`

	// Optionally full connection string in which case the above are ignored
	Connection string

	// Path to db for sqlite
	Path string `default:"ghz.db"`
}

// GetConnectionString returns the database connection string
func (db *DBConfig) GetConnectionString() string {
	if strings.TrimSpace(db.Connection) != "" {
		return db.Connection
	}

	if db.Type == "sqlite" || db.Type == "sqlite3" {
		return db.Path
	}

	connstr := ""

	if db.Type == "postgres" || db.Type == "cockroachdb" {
		connstr = fmt.Sprintf("host=%+v", db.Host)

		if db.Port != 0 {
			connstr = fmt.Sprintf("%+v port=%+v", connstr, db.Port)
		}

		if db.User != "" {
			connstr = fmt.Sprintf("%+v user=%+v", connstr, db.User)
		}

		connstr = fmt.Sprintf("%+v dbname=%+v", connstr, db.Name)

		if db.SSLMode != "" {
			connstr = fmt.Sprintf("%+v sslmode=%+v", connstr, db.SSLMode)
		}

		if db.Password != "" {
			connstr = fmt.Sprintf("%+v password=%+v", connstr, db.Password)
		}
	}

	if db.Type == "mysql" {
		addr := db.Host

		if db.Port != 0 {
			addr = fmt.Sprintf("%+v:%+v", addr, db.Port)
		}

		mysqlConfig := mysql.Config{
			User:   db.User,
			Passwd: db.Password,
			Addr:   addr,
			DBName: db.Name,

			// Our custom defaults
			Net:       "tcp",
			ParseTime: true,
		}
		return mysqlConfig.FormatDSN()
	}

	if db.Type == "mssql" {
		connstr = "sqlserver://"

		if db.User != "" {
			connstr = connstr + db.User
		}

		if db.Password != "" {
			connstr = fmt.Sprintf("%+v:%+v", connstr, db.Password)
		}

		if db.User != "" || db.Password != "" {
			connstr = connstr + "@"
		}

		if db.Host != "" {
			connstr = connstr + db.Host

			if db.Port != 0 {
				connstr = fmt.Sprintf("%+v:%+v", connstr, db.Port)
			}
		}

		connstr = fmt.Sprintf("%+v?database=%+v", connstr, db.Name)
	}

	return connstr
}

// GetDialect gets compatible GORM dialect
func (db *DBConfig) GetDialect() string {
	if db.Type == "sqlite" || db.Type == "sqlite3" {
		return "sqlite3"
	}

	if db.Type == "cockroachdb" {
		return "postgres"
	}

	return db.Type
}

// Validate the database config
func (db *DBConfig) Validate() error {
	dbtype := strings.ToLower(db.Type)

	supported := dbtype == "sqlite" ||
		dbtype == "sqlite3" ||
		dbtype == "postgres" ||
		dbtype == "mysql" ||
		dbtype == "mssql" ||
		dbtype == "cockroachdb"

	if !supported {
		return errors.New("Unsupported database type: " + db.Type)
	}

	db.Type = dbtype

	if db.Type == "sqlite" || db.Type == "sqlite3" {
		if err := requiredString(db.Path); err != nil {
			return errors.Wrap(err, "Database.Path")
		}
	}

	return nil
}
