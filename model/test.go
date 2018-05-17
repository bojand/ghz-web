package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Threshold represends a threshold limit we may care about
type Threshold int

const (
	// ThresholdMean is the threshold for mean / average
	ThresholdMean Threshold = iota

	// ThresholdMedian is the threshold for the median
	ThresholdMedian

	// Threshold95th is the threshold for the 95th percentile
	Threshold95th

	// Threshold99th is the threshold for the 96th percentile
	Threshold99th
)

// String() is the string representation of threshold
func (t Threshold) String() string {
	tholds := [...]string{"mean", "median", "95th", "99th"}

	if t < ThresholdMean || t > Threshold99th {
		return ""
	}

	return tholds[t]
}

// ThresholdFromString creates a Threshold from a string
func ThresholdFromString(str string) Threshold {
	str = strings.ToLower(str)
	if str == "1" || str == "median" {
		return ThresholdMedian
	} else if str == "2" || str == "95th" {
		return Threshold95th
	} else if str == "3" || str == "99th" {
		return Threshold99th
	}

	return ThresholdMean
}

// UnmarshalJSON prases a Threshold value from JSON string
func (t *Threshold) UnmarshalJSON(b []byte) error {
	*t = ThresholdFromString(string(b))

	return nil
}

// MarshalJSON formats a Threshold value into a JSON string
func (t Threshold) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}

// TestStatus represents a status of a test, whether its latest run failed the threshold settings
type TestStatus int

const (
	// StatusOK means the latest run in test was within the threshold
	StatusOK TestStatus = iota

	// StatusFail means the latest run in test was not within the threshold
	StatusFail
)

// String() is the string representation of threshold
func (t TestStatus) String() string {
	if t == StatusFail {
		return "fail"
	}

	return "ok"
}

// TestStatusFromString creates a TestStatus from a string
func TestStatusFromString(str string) TestStatus {
	str = strings.ToLower(str)

	t := StatusOK

	if str == "1" || str == "fail" {
		t = StatusFail
	}

	return t
}

// UnmarshalJSON prases a Threshold value from JSON string
func (t *TestStatus) UnmarshalJSON(b []byte) error {
	*t = TestStatusFromString(string(b))

	return nil
}

// MarshalJSON formats a Threshold value into a JSON string
func (t TestStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}

// ThresholdSetting setting
type ThresholdSetting struct {
	Status    TestStatus    `json:"status"`
	Threshold time.Duration `json:"threshold"`
}

// MarshalJSON formats a ThresholdSetting value into a JSON string
func (m ThresholdSetting) MarshalJSON() ([]byte, error) {
	type Alias ThresholdSetting
	aux := &struct {
		Status string `json:"status"`
		*Alias
	}{
		Status: m.Status.String(),
		Alias:  (*Alias)(&m),
	}

	return json.Marshal(aux)
}

// UnmarshalJSON prases a ThresholdSetting value from JSON string
func (m *ThresholdSetting) UnmarshalJSON(data []byte) error {
	type Alias ThresholdSetting
	aux := &struct {
		Status string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.Status = TestStatusFromString(aux.Status)

	return nil
}

var constants = [4]Threshold{ThresholdMean, ThresholdMedian, Threshold95th, Threshold99th}

// Test represents a test
type Test struct {
	gorm.Model
	ProjectID      uint                            `json:"projectID"`
	Project        Project                         `json:"-"`
	Name           string                          `json:"name" gorm:"unique_index"`
	Description    string                          `json:"description"`
	Status         TestStatus                      `json:"status"`
	Thresholds     map[Threshold]*ThresholdSetting `json:"thresholds" gorm:"-"`
	FailOnError    bool                            `json:"failOnError"`
	ThresholdsJSON string                          `json:"-" gorm:"column:thresholds"`
}

// MarshalJSON formats a Test value into a JSON string
func (t Test) MarshalJSON() ([]byte, error) {
	type Alias Test
	aux := &struct {
		Thresholds map[string]*ThresholdSetting `json:"thresholds"`
		*Alias
	}{
		Thresholds: thresholdsToStringMap(t.Thresholds),
		Alias:      (*Alias)(&t),
	}

	return json.Marshal(aux)
}

// UnmarshalJSON prases a Test value from JSON string
func (t *Test) UnmarshalJSON(data []byte) error {
	type Alias Test
	aux := &struct {
		Thresholds map[string]*ThresholdSetting `json:"thresholds"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Thresholds = thresholdsFromStringMap(aux.Thresholds)

	return nil
}

// thresholdsToStringMap converts a threshold map string map
func thresholdsToStringMap(tholds map[Threshold]*ThresholdSetting) map[string]*ThresholdSetting {
	m := make(map[string]*ThresholdSetting)

	for k, v := range tholds {
		m[k.String()] = v
	}

	return m
}

// thresholdsToStringMap converts a threshold map string map
func thresholdsFromStringMap(tholds map[string]*ThresholdSetting) map[Threshold]*ThresholdSetting {
	m := make(map[Threshold]*ThresholdSetting)

	for k, v := range tholds {
		kt := ThresholdFromString(k)
		m[kt] = v
	}

	return m
}

// BeforeSave is called by GORM before save
func (t *Test) BeforeSave() error {
	tholds, err := json.Marshal(thresholdsToStringMap(t.Thresholds))
	if err != nil {
		return err
	}

	t.ThresholdsJSON = string(tholds)
	return nil
}

// AfterFind is called by GORM after a query
func (t *Test) AfterFind() error {
	tholds := strings.TrimSpace(t.ThresholdsJSON)
	if tholds != "" {
		var dat map[string]*ThresholdSetting
		if err := json.Unmarshal([]byte(tholds), &dat); err != nil {
			return err
		}
		t.Thresholds = thresholdsFromStringMap(dat)
	}
	return nil
}

// SetStatus sets this test's status based on the settings and the values in params
func (t *Test) SetStatus(tMean, tMedian, t95, t99 time.Duration, hasError bool) {
	// reset our status
	t.Status = StatusOK

	compareVal := []time.Duration{tMean, tMedian, t95, t99}

	for i, thc := range constants {
		if t.Thresholds[thc] != nil {

			// reset each threshold status
			t.Thresholds[thc].Status = StatusOK

			val := compareVal[i]

			if t.Thresholds[thc].Threshold > 0 && val > 0 && val > t.Thresholds[thc].Threshold {
				t.Thresholds[thc].Status = StatusFail
				t.Status = StatusFail
			}
		}
	}

	if t.FailOnError && hasError {
		t.Status = StatusFail
	}
}
