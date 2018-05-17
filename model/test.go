package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Threshold represends a threshold limit we may care about
type Threshold string

const (
	// ThresholdMean is the threshold for mean / average
	ThresholdMean Threshold = Threshold("mean")

	// ThresholdMedian is the threshold for the median
	ThresholdMedian = Threshold("median")

	// Threshold95th is the threshold for the 95th percentile
	Threshold95th = Threshold("95th")

	// Threshold99th is the threshold for the 96th percentile
	Threshold99th = Threshold("99th")
)

// TestStatus represents a status of a test, whether its latest run failed the threshold settings
type TestStatus string

// String() is the string representation of threshold
func (t TestStatus) String() string {
	if t == StatusFail {
		return "fail"
	}

	return "ok"
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

// TestStatusFromString creates a TestStatus from a string
func TestStatusFromString(str string) TestStatus {
	str = strings.ToLower(str)

	t := StatusOK

	if str == "fail" {
		t = StatusFail
	}

	return t
}

const (
	// StatusOK means the latest run in test was within the threshold
	StatusOK TestStatus = TestStatus("ok")

	// StatusFail means the latest run in test was not within the threshold
	StatusFail = TestStatus("fail")
)

// ThresholdSetting setting
type ThresholdSetting struct {
	Status    TestStatus    `json:"status"`
	Threshold time.Duration `json:"threshold"`
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
	Thresholds     map[Threshold]*ThresholdSetting `json:"thresholds,omitempty" gorm:"-"`
	FailOnError    bool                            `json:"failOnError"`
	ThresholdsJSON string                          `json:"-" gorm:"column:thresholds"`
}

// BeforeSave is called by GORM before save
func (t *Test) BeforeSave() error {
	tholds := []byte("")
	if t.Thresholds != nil && len(t.Thresholds) > 0 {
		var err error
		tholds, err = json.Marshal(t.Thresholds)
		if err != nil {
			return err
		}
	}

	t.ThresholdsJSON = string(tholds)
	return nil
}

// AfterFind is called by GORM after a query
func (t *Test) AfterFind() error {
	tholds := strings.TrimSpace(t.ThresholdsJSON)
	if tholds != "" {
		dat := map[Threshold]*ThresholdSetting{}
		if err := json.Unmarshal([]byte(tholds), &dat); err != nil {
			return err
		}
		t.Thresholds = dat
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
