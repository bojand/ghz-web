package model

import (
	"fmt"
	"strings"
)

// Status represents a status of a test, whether its latest run failed the threshold settings
type Status string

// String() is the string representation of threshold
func (t Status) String() string {
	if t == StatusFail {
		return "fail"
	}

	return "ok"
}

// UnmarshalJSON prases a Threshold value from JSON string
func (t *Status) UnmarshalJSON(b []byte) error {
	*t = StatusFromString(string(b))

	return nil
}

// MarshalJSON formats a Threshold value into a JSON string
func (t Status) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}

// StatusFromString creates a Status from a string
func StatusFromString(str string) Status {
	str = strings.ToLower(str)

	t := StatusOK

	if str == "fail" {
		t = StatusFail
	}

	return t
}

const (
	// StatusOK means the latest run in test was within the threshold
	StatusOK Status = Status("ok")

	// StatusFail means the latest run in test was not within the threshold
	StatusFail = Status("fail")
)
