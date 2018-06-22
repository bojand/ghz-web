package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_String(t *testing.T) {
	var tests = []struct {
		name     string
		in       Status
		expected string
	}{
		{"ok", StatusOK, "ok"},
		{"fail", StatusFail, "fail"},
		{"unknown", Status("foo"), "ok"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.in.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStatus_StatusFromString(t *testing.T) {
	var tests = []struct {
		name     string
		in       string
		expected Status
	}{
		{"ok", "ok", StatusOK},
		{"fail", "fail", StatusFail},
		{"unknown", "foo", StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := StatusFromString(tt.in)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
