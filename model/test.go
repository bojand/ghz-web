package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
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
	Status    TestStatus    `json:"status" validate:"oneof=ok fail"`
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
	ProjectID      uint                            `json:"projectID" gorm:"type:integer REFERENCES projects(id)"`
	Project        *Project                        `json:"-"`
	Name           string                          `json:"name" gorm:"unique_index;not null" validate:"required"`
	Description    string                          `json:"description"`
	Status         TestStatus                      `json:"status" validate:"oneof=ok fail"`
	Thresholds     map[Threshold]*ThresholdSetting `json:"thresholds,omitempty" gorm:"-"`
	FailOnError    bool                            `json:"failOnError"`
	ThresholdsJSON string                          `json:"-" gorm:"column:thresholds"`
}

// BeforeCreate is a GORM hook called when a model is created
func (t *Test) BeforeCreate() error {
	if t.Name == "" {
		t.Name = random.String(16)
	}

	return nil
}

// BeforeUpdate is a GORM hook called when a model is updated
func (t *Test) BeforeUpdate() error {
	if t.Name == "" {
		return errors.New("Test name cannot be empty")
	}

	return nil
}

// BeforeSave is called by GORM before save
func (t *Test) BeforeSave() error {
	if t.ProjectID == 0 && t.Project == nil {
		return errors.New("Test must belong to a project")
	}

	tholds := []byte("")
	if t.Thresholds != nil && len(t.Thresholds) > 0 {
		var err error
		tholds, err = json.Marshal(t.Thresholds)
		if err != nil {
			return err
		}
	}

	t.ThresholdsJSON = string(tholds)

	name := strings.Replace(t.Name, " ", "", -1)
	t.Name = strings.ToLower(name)
	t.Description = strings.TrimSpace(t.Description)

	return nil
}

// AfterSave is called by GORM after model is saved during create or update
func (t *Test) AfterSave() error {
	t.ThresholdsJSON = ""
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

	t.ThresholdsJSON = ""

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

// TestService is our implementation
type TestService struct {
	DB *gorm.DB
}

// FindByID finds project by id
func (ts *TestService) FindByID(id uint) (*Test, error) {
	t := new(Test)
	err := ts.DB.First(t, id).Error
	if err != nil {
		t = nil
	}
	return t, err
}

// FindByName finds project by name
func (ts *TestService) FindByName(name string) (*Test, error) {
	name = strings.ToLower(name)
	t := new(Test)
	err := ts.DB.First(t, "name = ?", name).Error
	if err != nil {
		t = nil
	}
	return t, err
}

// FindByProjectID finds tests by project
func (ts *TestService) FindByProjectID(pid uint, num, page int) ([]*Test, error) {
	p := &Project{}
	p.ID = pid

	s := make([]*Test, 0)

	offset := -1
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	err := ts.DB.Offset(offset).Limit(num).Order("name desc").Model(p).Related(&s).Error

	return s, err
}

// Create creates a new project
func (ts *TestService) Create(t *Test) error {
	return ts.DB.Create(t).Error
}

// Update updates  project
func (ts *TestService) Update(t *Test) error {
	testToUpdate := &Test{}
	if err := ts.DB.First(testToUpdate, t.ID).Error; err != nil {
		return err
	}

	name := strings.Replace(t.Name, " ", "", -1)
	if name == "" {
		t.Name = testToUpdate.Name
	}

	return ts.DB.Save(t).Error
}

// Delete deletes project
func (ts *TestService) Delete(t *Test) error {
	return errors.New("Not Implemented Yet")
}
