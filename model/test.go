package model

import (
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

// TestStatus represents a status of a test, whether its latest run failed the threshold settings
type TestStatus int

const (
	// StatusOK means the latest run in test was within the threshold
	StatusOK TestStatus = iota

	// StatusFail means the latest run in test was not within the threshold
	StatusFail
)

// Test represents a test
type Test struct {
	gorm.Model
	Name        string                          `json:"name" gorm:"unique_index"`
	Description string                          `json:"description"`
	Status      TestStatus                      `json:"status"`
	Thresholds  map[Threshold]*ThresholdSetting `json:"thresholds"`
	FailOnError bool                            `json:"failOnError"`
}

// ThresholdSetting setting
type ThresholdSetting struct {
	Status    TestStatus    `json:"status"`
	Threshold time.Duration `json:"threshold"`
}

// SetStatus sets this test's status based on the settings and the values in params
func (t *Test) SetStatus(tMean, tMedian, t95, t99 time.Duration, hasError bool) {
	// reset our status
	t.Status = StatusOK

	constants := []Threshold{ThresholdMean, ThresholdMedian, Threshold95th, Threshold99th}
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
