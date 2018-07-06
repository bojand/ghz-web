package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetailModel_BeforeSave(t *testing.T) {
	var tests = []struct {
		name        string
		in          *Detail
		expected    *Detail
		expectError bool
	}{
		{"no run", &Detail{Latency: 123.45}, &Detail{Latency: 123.45}, true},
		{"trim error", &Detail{RunID: 1, Latency: 123.45, Error: " network error "}, &Detail{RunID: 1, Latency: 123.45, Error: "network error", Status: "OK"}, false},
		{"trim status", &Detail{RunID: 1, Latency: 123.45, Error: " network error ", Status: " OK "}, &Detail{RunID: 1, Latency: 123.45, Error: "network error", Status: "OK"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.BeforeSave(nil)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}
