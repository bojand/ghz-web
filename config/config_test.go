package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Read(t *testing.T) {
	var tests = []struct {
		name     string
		in       string
		expected *Config
	}{
		{"config1.toml",
			"../test/config1.toml",
			&Config{
				Server:   ServerConfig{Port: 3000},
				Database: DBConfig{Type: "sqlite", Host: "localhost", Name: "ghz", Path: "ghz.db", SSLMode: "disable"},
				Log:      LogConfig{Level: "info"}}},
		{"config2.toml",
			"../test/config2.toml",
			&Config{
				Server:   ServerConfig{Port: 4321},
				Database: DBConfig{Type: "postgres", Host: "123.0.0.1", Name: "ghz", Path: "ghz.db", SSLMode: "disable", User: "dbuser", Port: 1234},
				Log:      LogConfig{Level: "warn", Path: "/tmp/ghz.log"}}},
		{"config3.toml",
			"../test/config3.toml",
			&Config{
				Server:   ServerConfig{Port: 3000},
				Database: DBConfig{Type: "postgres", Host: "localhost", Name: "ghz", Path: "ghz.db", SSLMode: "disable"},
				Log:      LogConfig{Level: "debug", Path: ""}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Read(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLogConfig_Validate(t *testing.T) {
	var tests = []struct {
		name     string
		in       *LogConfig
		expected string
	}{
		{"level=unknown", &LogConfig{Level: "unknown"}, "Unsupported log level: unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.in.Validate()
			assert.Equal(t, tt.expected, actual.Error())
		})
	}
}

func TestServerConfig_GetHostPort(t *testing.T) {
	var tests = []struct {
		name     string
		in       *ServerConfig
		expected string
	}{
		{"with port", &ServerConfig{Address: "127.0.0.1", Port: 2000}, "127.0.0.1:2000"},
		{"no address", &ServerConfig{Port: 1000}, ":1000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.in.GetHostPort()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
