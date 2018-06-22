package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunModel_BeforeSave(t *testing.T) {
	var runs = []struct {
		name        string
		in          *Run
		expected    *Run
		expectError bool
	}{
		{"no test id", &Run{}, &Run{}, true},
		{"with test id", &Run{TestID: 123}, &Run{TestID: 123, Status: "ok"}, false},
		{"with error dist",
			&Run{TestID: 123, ErrorDist: map[string]int{"foo": 1, "bar": 2}},
			&Run{TestID: 123, ErrorDist: map[string]int{"foo": 1, "bar": 2}, ErrorDistJSON: "{\"bar\":2,\"foo\":1}", Status: "fail"},
			false},
		{"with status dist",
			&Run{TestID: 123, StatusCodeDist: map[string]int{"foo": 1, "bar": 2}},
			&Run{TestID: 123, StatusCodeDist: map[string]int{"foo": 1, "bar": 2}, StatusCodeDistJSON: "{\"bar\":2,\"foo\":1}", Status: "ok"},
			false},
	}

	for _, tt := range runs {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.BeforeSave()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}

func TestRunModel_AfterSave(t *testing.T) {
	var runs = []struct {
		name        string
		in          *Run
		expected    *Run
		expectError bool
	}{
		{"no test id", &Run{}, &Run{}, false},
		{"with test id", &Run{TestID: 123}, &Run{TestID: 123}, false},
		{"with error dist",
			&Run{TestID: 123, ErrorDist: map[string]int{"foo": 1, "bar": 2}, ErrorDistJSON: "{\"bar\":2,\"foo\":1}"},
			&Run{TestID: 123, ErrorDist: map[string]int{"foo": 1, "bar": 2}},
			false},
		{"with status dist",
			&Run{TestID: 123, StatusCodeDist: map[string]int{"foo": 1, "bar": 2}, StatusCodeDistJSON: "{\"bar\":2,\"foo\":1}"},
			&Run{TestID: 123, StatusCodeDist: map[string]int{"foo": 1, "bar": 2}},
			false},
	}

	for _, tt := range runs {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.AfterSave()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}

func TestRunModel_AfterFind(t *testing.T) {
	var runs = []struct {
		name        string
		in          *Run
		expected    *Run
		expectError bool
	}{
		{"no test id", &Run{}, &Run{}, false},
		{"with test id", &Run{TestID: 123}, &Run{TestID: 123}, false},
		{"with error dist",
			&Run{TestID: 123, ErrorDistJSON: "{\"bar\":2,\"foo\":1}"},
			&Run{TestID: 123, ErrorDist: map[string]int{"foo": 1, "bar": 2}},
			false},
		{"with status dist",
			&Run{TestID: 123, StatusCodeDistJSON: "{\"bar\":2,\"foo\":1}"},
			&Run{TestID: 123, StatusCodeDist: map[string]int{"foo": 1, "bar": 2}},
			false},
	}

	for _, tt := range runs {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.AfterFind()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}
