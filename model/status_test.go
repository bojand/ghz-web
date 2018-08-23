package model

import (
	"encoding/json"
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

func TestStatus_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name     string
		in       string
		expected Status
	}{
		{"ok", `"ok"`, StatusOK},
		{"OK", `"OK"`, StatusOK},
		{"fail", `"fail"`, StatusFail},
		{"FAIL", `"FAIL"`, StatusFail},
		{" FAIL ", ` "FAIL" `, StatusFail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Status
			err := json.Unmarshal([]byte(tt.in), &actual)
			assert.NoError(t, err)
			// fmt.Println(tt.expected)
			// fmt.Println(actual)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
