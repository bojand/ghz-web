package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Run represents a project
type Run struct {
	Model
	TestID uint  `json:"testID" gorm:"type:integer REFERENCES tests(id)"`
	Test   *Test `json:"-"`

	Count   uint64        `json:"count"`
	Total   time.Duration `json:"total"`
	Average time.Duration `json:"average"`
	Fastest time.Duration `json:"fastest"`
	Slowest time.Duration `json:"slowest"`
	Rps     float64       `json:"rps"`

	Status Status `json:"status" validate:"oneof=ok fail"`

	ErrorDist      map[string]int `json:"errorDistribution,omitempty" gorm:"-"`
	StatusCodeDist map[string]int `json:"statusCodeDistribution,omitempty" gorm:"-"`

	ErrorDistJSON      string `json:"-" gorm:"column:error_dist"`
	StatusCodeDistJSON string `json:"-" gorm:"column:status_code_dist"`
}

// BeforeSave is called by GORM before save
func (r *Run) BeforeSave() error {
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

	return nil
}

// AfterSave is called by GORM after model is saved during create or update
func (r *Run) AfterSave() error {
	r.ErrorDistJSON = ""
	r.StatusCodeDistJSON = ""
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

	return nil
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
