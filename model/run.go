package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// LatencyDistribution holds latency distribution data
type LatencyDistribution struct {
	Model
	RunID      uint          `json:"runId"`
	Percentage int           `json:"percentage"`
	Latency    time.Duration `json:"latency"`
}

// Bucket holds histogram data
type Bucket struct {
	Model
	RunID uint `json:"runId"`

	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

// Options represents run options
type Options struct {
	Call          string             `json:"call,omitempty"`
	Proto         string             `json:"proto,omitempty"`
	Host          string             `json:"host,omitempty"`
	Cert          string             `json:"cert,omitempty"`
	CName         string             `json:"cname,omitempty"`
	N             int                `json:"n,omitempty"`
	C             int                `json:"c,omitempty"`
	QPS           int                `json:"qps,omitempty"`
	Z             time.Duration      `json:"z,omitempty"`
	Timeout       int                `json:"timeout,omitempty"`
	DialTimtout   int                `json:"dialTimeout,omitempty"`
	KeepaliveTime int                `json:"keepAlice,omitempty"`
	Data          interface{}        `json:"data,omitempty"`
	Metadata      *map[string]string `json:"metadata,omitempty"`
}

// Run represents a project
type Run struct {
	Model
	Date    time.Time     `json:"date"`
	TestID  uint          `json:"testID" gorm:"type:integer REFERENCES tests(id)"`
	Test    *Test         `json:"-"`
	Count   uint64        `json:"count"`
	Total   time.Duration `json:"total"`
	Average time.Duration `json:"average"`
	Fastest time.Duration `json:"fastest"`
	Slowest time.Duration `json:"slowest"`
	Rps     float64       `json:"rps"`

	Status Status `json:"status" validate:"oneof=ok fail"`

	Options *Options `json:"options,omitempty" gorm:"-"`

	ErrorDist      map[string]int `json:"errorDistribution,omitempty" gorm:"-"`
	StatusCodeDist map[string]int `json:"statusCodeDistribution,omitempty" gorm:"-"`

	LatencyDistribution []*LatencyDistribution `json:"latencyDistribution"`
	Histogram           []*Bucket              `json:"histogram"`

	// temp conversion vars
	ErrorDistJSON      string `json:"-" gorm:"column:error_dist"`
	StatusCodeDistJSON string `json:"-" gorm:"column:status_code_dist"`
	OptionsJSON        string `json:"-" gorm:"column:options"`
}

// BeforeSave is called by GORM before save
func (r *Run) BeforeSave(scope *gorm.Scope) error {
	if r.TestID == 0 && r.Test == nil {
		return errors.New("Run must belong to a test")
	}

	r.Status = StatusOK

	errDist := []byte("")
	if r.ErrorDist != nil && len(r.ErrorDist) > 0 {
		var err error
		errDist, err = json.Marshal(r.ErrorDist)
		if err != nil {
			return err
		}

		r.Status = StatusFail
	}

	r.ErrorDistJSON = string(errDist)

	statusCodeDist := []byte("")
	if r.StatusCodeDist != nil && len(r.StatusCodeDist) > 0 {
		var err error
		statusCodeDist, err = json.Marshal(r.StatusCodeDist)
		if err != nil {
			return err
		}
	}

	r.StatusCodeDistJSON = string(statusCodeDist)

	options := []byte("")
	if r.Options != nil {
		var err error
		options, err = json.Marshal(r.Options)
		if err != nil {
			return err
		}
	}

	r.OptionsJSON = string(options)

	if scope != nil {
		scope.SetColumn("status", r.Status)
		scope.SetColumn("error_dist", r.ErrorDistJSON)
		scope.SetColumn("status_code_dist", r.StatusCodeDistJSON)
		scope.SetColumn("options", r.OptionsJSON)
	}

	return nil
}

// AfterSave is called by GORM after model is saved during create or update
func (r *Run) AfterSave() error {
	r.ErrorDistJSON = ""
	r.StatusCodeDistJSON = ""
	r.OptionsJSON = ""
	return nil
}

// AfterFind is called by GORM after a query
func (r *Run) AfterFind() error {
	errDist := strings.TrimSpace(r.ErrorDistJSON)
	if errDist != "" {
		dat := map[string]int{}
		if err := json.Unmarshal([]byte(errDist), &dat); err != nil {
			return err
		}
		r.ErrorDist = dat
	}

	r.ErrorDistJSON = ""

	statusCodeDist := strings.TrimSpace(r.StatusCodeDistJSON)
	if statusCodeDist != "" {
		dat := map[string]int{}
		if err := json.Unmarshal([]byte(statusCodeDist), &dat); err != nil {
			return err
		}
		r.StatusCodeDist = dat
	}

	r.StatusCodeDistJSON = ""

	options := strings.TrimSpace(r.OptionsJSON)
	if options != "" {
		r.Options = new(Options)
		if err := json.Unmarshal([]byte(options), r.Options); err != nil {
			return err
		}
	}

	r.OptionsJSON = ""

	return nil
}

// GetThresholdValues gets median and 95th
func (r *Run) GetThresholdValues() (time.Duration, time.Duration) {
	var median, nine5 time.Duration

	latencies := len(r.LatencyDistribution)

	if latencies > 0 {
		for _, l := range r.LatencyDistribution {
			// record median
			if l.Percentage == 50 {
				median = l.Latency
			}

			// record 95th
			if l.Percentage == 95 {
				nine5 = l.Latency
			}
		}
	}

	return median, nine5
}

// HasErrors returns whether run has any errors
func (r *Run) HasErrors() bool {
	hasErrors := false
	if r.ErrorDist != nil && len(r.ErrorDist) > 0 {
		hasErrors = true
	}

	return hasErrors
}

// RunService is our implementation
type RunService struct {
	DB *gorm.DB
}

// Create creates a new run
func (rs *RunService) Create(r *Run) error {
	return rs.DB.Create(r).Error
}

// Count returns the total number of runs
func (rs *RunService) Count(tid uint) (uint, error) {
	count := uint(0)
	err := rs.DB.Model(&Run{}).Where("test_id = ?", tid).Count(&count).Error
	return count, err
}

// FindByID finds run by id
func (rs *RunService) FindByID(id uint) (*Run, error) {
	r := new(Run)
	r.Histogram = make([]*Bucket, 10)
	r.LatencyDistribution = make([]*LatencyDistribution, 10)

	err := rs.DB.First(r, id).Related(&r.Histogram).Related(&r.LatencyDistribution).Error

	if err != nil {
		r = nil
	}

	return r, err
}

// FindLatest returns the latest created run for test
func (rs *RunService) FindLatest(tid uint) (*Run, error) {
	r := new(Run)
	r.Histogram = make([]*Bucket, 100)
	r.LatencyDistribution = make([]*LatencyDistribution, 100)

	fmt.Printf("Quering\n\n")

	err := rs.DB.Model(&Run{}).Where("test_id = ?", tid).Order("date desc").First(r).Error

	if err != nil {
		r = nil
		return r, nil
	}

	err = rs.DB.Model(r).Related(&r.Histogram).Related(&r.LatencyDistribution).Error
	if err != nil {
		r = nil
	}

	return r, err
}

// FindByTestID finds tests by project
func (rs *RunService) FindByTestID(tid, num, page uint, populate bool) ([]*Run, error) {
	t := &Test{}
	t.ID = tid

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	s := make([]*Run, num)

	err := rs.DB.Offset(offset).Limit(num).Order("id desc").Model(t).Related(&s).Error

	if populate {
		if err != nil {
			return nil, err
		}

		for _, run := range s {
			run.LatencyDistribution = make([]*LatencyDistribution, 10)
			err := rs.DB.Model(run).Related(&run.LatencyDistribution).Error
			if err != nil {
				return nil, err
			}

			run.Histogram = make([]*Bucket, 10)
			err = rs.DB.Model(run).Related(&run.Histogram).Error
			if err != nil {
				return nil, err
			}
		}
	}

	return s, err
}

// FindByTestIDSorted lists tests using sorting
func (rs *RunService) FindByTestIDSorted(tid, num, page uint, sortField, order string,
	histogram bool, latency bool) ([]*Run, error) {
	if (sortField != "id" && sortField != "count" && sortField != "total" && sortField != "average" &&
		sortField != "fastest" && sortField != "slowest" && sortField != "rps" && sortField != "date") ||
		(order != "asc" && order != "desc") {
		return nil, errors.New("Invalid sort parameters")
	}

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	orderSQL := sortField + " " + order

	t := &Test{}
	t.ID = tid

	s := make([]*Run, num)

	err := rs.DB.Order(orderSQL).Offset(offset).Limit(num).Model(t).Related(&s).Error

	if err != nil {
		return nil, err
	}

	if histogram {
		for _, run := range s {
			run.Histogram = make([]*Bucket, 10)
			err = rs.DB.Model(run).Related(&run.Histogram).Error
			if err != nil {
				return nil, err
			}
		}
	}

	if latency {
		for _, run := range s {
			run.LatencyDistribution = make([]*LatencyDistribution, 10)
			err := rs.DB.Model(run).Related(&run.LatencyDistribution).Error
			if err != nil {
				return nil, err
			}
		}
	}

	return s, err
}

// Update updates a run
func (rs *RunService) Update(r *Run) error {
	runToUpdate := &Run{}
	if err := rs.DB.First(runToUpdate, r.ID).Error; err != nil {
		return err
	}

	return rs.DB.Save(r).Error
}

// Delete deletes a run
func (rs *RunService) Delete(r *Run) error {
	return errors.New("Not Implemented Yet")
}
