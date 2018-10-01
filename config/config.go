package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

// Info represents some app level info
type Info struct {
	Version   string
	GOVersion string
	StartTime time.Time
}

// ServerConfig is server config
type ServerConfig struct {
	RootURL string
	Address string `default:"localhost"`
	Port    uint   `default:"3000"`
}

// GetHostPort returns host:port
func (s *ServerConfig) GetHostPort() string {
	return s.Address + ":" + strconv.FormatUint(uint64(s.Port), 10)
}

// LogConfig log settings
type LogConfig struct {
	Level string `default:"info"`
	Path  string
}

// Validate validates the config settings
func (lc *LogConfig) Validate() error {
	lvl := strings.ToLower(lc.Level)
	lvl = strings.TrimSpace(lvl)

	supported := lvl == "off" ||
		lvl == "error" ||
		lvl == "warn" ||
		lvl == "info" ||
		lvl == "debug"

	if !supported {
		return errors.New("Unsupported log level: " + lc.Level)
	}

	lc.Level = lvl

	return nil
}

// Config is the application config
type Config struct {
	Database DBConfig
	Server   ServerConfig
	Log      LogConfig
}

// Validate the config
func (c *Config) Validate() error {
	err := c.Database.Validate()
	if err != nil {
		return err
	}

	err = c.Log.Validate()
	if err != nil {
		return err
	}

	c.Server.RootURL = strings.TrimSpace(c.Server.RootURL)

	return nil
}

// Read the config file
func Read(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		path = "config.toml"
	}

	config := Config{}

	configor.Load(&config, path)

	err := config.Validate()

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func requiredString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("is required")
	}

	return nil
}
