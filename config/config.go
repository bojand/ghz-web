package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Database host database configs
type Database struct {
	// The SQL dialect we support
	// One of mysql, postgres, mssql, sqlite3
	Type string

	// The database connection options
	Host     string
	Name     string
	User     string
	Password string

	// Optionally full connection string in which case the above are ignored
	Connection string

	// Path to db for sqlite
	Path string
}

// GetConnectionString returns the database connection string
func (db *Database) GetConnectionString() string {
	if strings.TrimSpace(db.Connection) != "" {
		return db.Connection
	}

	if db.Type == "sqlite" {
		return db.Path
	}

	return ""
}

// Server config
type Server struct {
	RootURL string
	Address string
	Port    uint
}

// GetHostPort returns host:port
func (s *Server) GetHostPort() string {
	return s.Address + ":" + strconv.FormatUint(uint64(s.Port), 10)
}

// Config is the application config
type Config struct {
	DB  Database `toml:"database"`
	Srv Server   `toml:"server"`
}

// Validate the config
func (c *Config) Validate() error {
	if err := requiredString(c.DB.Type); err != nil {
		return errors.Wrap(err, "Database.Type")
	}

	if c.DB.Type == "sqlite" {
		if err := requiredString(c.DB.Path); err != nil {
			return errors.Wrap(err, "Database.Oath")
		}
	}

	return nil
}

// Read the config file
func Read(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		path = "config.toml"
	}

	config := Config{
		DB: Database{
			Type: "sqlite3",
			Path: "test.db",
		},
		Srv: Server{
			Port: 3000,
		},
	}

	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, &config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func requiredString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("is required")
	}

	return nil
}
