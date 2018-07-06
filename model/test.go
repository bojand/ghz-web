package model

import (
	"encoding/json"
	"errors"
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

	// Threshold99th is the threshold for the 99th percentile
	Threshold99th = Threshold("99th")
)

// ThresholdSetting setting
type ThresholdSetting struct {
	Status    Status        `json:"status" validate:"oneof=ok fail"`
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

	m.Status = StatusFromString(aux.Status)

	return nil
}

var constants = [4]Threshold{ThresholdMean, ThresholdMedian, Threshold95th, Threshold99th}

// Test represents a test
type Test struct {
	Model
	ProjectID      uint                            `json:"projectID" gorm:"type:integer REFERENCES projects(id)"`
	Project        *Project                        `json:"-"`
	Name           string                          `json:"name" gorm:"unique_index;not null" validate:"required"`
	Description    string                          `json:"description"`
	Status         Status                          `json:"status" validate:"oneof=ok fail"`
	Thresholds     map[Threshold]*ThresholdSetting `json:"thresholds,omitempty" gorm:"-"`
	FailOnError    bool                            `json:"failOnError"`
	ThresholdsJSON string                          `json:"-" gorm:"column:thresholds"`
}

// UnmarshalJSON prases a Test value from JSON string
func (t *Test) UnmarshalJSON(data []byte) error {
	type Alias Test
	aux := &struct {
		Status string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Status = StatusFromString(aux.Status)

	return nil
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
func (t *Test) BeforeSave(scope *gorm.Scope) error {
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

	if scope != nil {
		scope.SetColumn("name", t.Name)
		scope.SetColumn("description", t.Description)
		scope.SetColumn("thresholds", t.ThresholdsJSON)
	}

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

// FindByID finds test by id
func (ts *TestService) FindByID(id uint) (*Test, error) {
	t := new(Test)
	err := ts.DB.First(t, id).Error
	if err != nil {
		t = nil
	}
	return t, err
}

// FindByName finds test by name
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
func (ts *TestService) FindByProjectID(pid, num, page uint) ([]*Test, error) {
	p := &Project{}
	p.ID = pid

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	s := make([]*Test, num)

	err := ts.DB.Offset(offset).Limit(num).Order("name desc").Model(p).Related(&s).Error

	return s, err
}

// FindByProjectIDSorted lists tests using sorting
func (ts *TestService) FindByProjectIDSorted(pid, num, page uint, sortField, order string) ([]*Test, error) {
	if (sortField != "name" && sortField != "id") || (order != "asc" && order != "desc") {
		return nil, errors.New("Invalid sort parameters")
	}

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	orderSQL := sortField + " " + order

	p := &Project{}
	p.ID = pid

	s := make([]*Test, num)

	err := ts.DB.Order(orderSQL).Offset(offset).Limit(num).Model(p).Related(&s).Error

	return s, err
}

// Count returns the total number of tests
func (ts *TestService) Count(pid uint) (uint, error) {
	count := uint(0)
	err := ts.DB.Model(&Test{}).Where("project_id = ?", pid).Count(&count).Error
	return count, err
}

// Create creates a new tests
func (ts *TestService) Create(t *Test) error {
	return ts.DB.Create(t).Error
}

// Update updates tests
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

// Delete deletes tests
func (ts *TestService) Delete(t *Test) error {
	return errors.New("Not Implemented Yet")
}
