package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	milli1 = 1 * time.Millisecond
	milli2 = 2 * time.Millisecond
	milli3 = 3 * time.Millisecond
	milli4 = 4 * time.Millisecond
	milli5 = 5 * time.Millisecond
)

func TestTestModel_SetStatus(t *testing.T) {
	var tests = []struct {
		name     string
		model    *Test
		in       [4]time.Duration
		inError  bool
		expected *Test
	}{
		{"empty", &Test{}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{}},
		{"with error true but no fail on error setting", &Test{}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, true, &Test{}},
		{"with error true but and fail on error", &Test{FailOnError: true}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, true, &Test{Status: StatusFail, FailOnError: true}},
		{"no values over limit", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5},
			Threshold95th:   &ThresholdSetting{Threshold: milli5},
			Threshold99th:   &ThresholdSetting{Threshold: milli5},
		}}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{Status: StatusOK, FailOnError: false, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold95th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
		}}},
		{"no values over limit but with error", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5},
			Threshold95th:   &ThresholdSetting{Threshold: milli5},
			Threshold99th:   &ThresholdSetting{Threshold: milli5},
		}, FailOnError: true}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, true, &Test{Status: StatusFail, FailOnError: true, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold95th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
		}}},
		{"mean over limit", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: 500 * time.Nanosecond},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5},
			Threshold95th:   &ThresholdSetting{Threshold: milli5},
			Threshold99th:   &ThresholdSetting{Threshold: milli5},
		}}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{Status: StatusFail, FailOnError: false, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: 500 * time.Nanosecond, Status: StatusFail},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold95th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
		}}},
		{"median over limit", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli1},
			ThresholdMedian: &ThresholdSetting{Threshold: milli1},
			Threshold95th:   &ThresholdSetting{Threshold: milli5},
			Threshold99th:   &ThresholdSetting{Threshold: milli5},
		}}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{Status: StatusFail, FailOnError: false, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli1, Status: StatusOK},
			ThresholdMedian: &ThresholdSetting{Threshold: milli1, Status: StatusFail},
			Threshold95th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
		}}},
		{"95th over limit", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5},
			Threshold95th:   &ThresholdSetting{Threshold: milli1},
			Threshold99th:   &ThresholdSetting{Threshold: milli5},
		}}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{Status: StatusFail, FailOnError: false, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold95th:   &ThresholdSetting{Threshold: milli1, Status: StatusFail},
			Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
		}}},
		{"99th over limit", &Test{Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5},
			Threshold95th:   &ThresholdSetting{Threshold: milli5},
			Threshold99th:   &ThresholdSetting{Threshold: milli3},
		}}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, false, &Test{Status: StatusFail, FailOnError: false, Thresholds: map[Threshold]*ThresholdSetting{
			ThresholdMean:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			ThresholdMedian: &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold95th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
			Threshold99th:   &ThresholdSetting{Threshold: milli3, Status: StatusFail},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.model
			actual.SetStatus(tt.in[0], tt.in[1], tt.in[2], tt.in[3], tt.inError)
			assert.Equal(t, tt.expected, actual)
			fmt.Printf("%+v\n%+v\n\n", tt.expected, actual)
		})
	}
}
