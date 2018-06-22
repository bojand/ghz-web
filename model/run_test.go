package model

import (
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

const (
	milli500  = 500 * time.Millisecond
	milli1000 = 1000 * time.Millisecond
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

func TestRunService_Create(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := RunService{DB: db}
	var tid, pid, rid uint

	t.Run("fail new without test", func(t *testing.T) {
		o := Run{
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}
		err := dao.Create(&o)

		assert.Error(t, err)
	})

	t.Run("new run and test and project", func(t *testing.T) {
		p := &Project{}

		o := &Test{
			Project:     p,
			Name:        "Test 111 ",
			Description: "Test Description Asdf ",
		}

		r := &Run{
			Test:    o,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		err := dao.Create(r)

		assert.NoError(t, err)
		assert.NotZero(t, p.ID)
		assert.NotEmpty(t, p.Name)
		assert.Equal(t, "", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		assert.NotZero(t, o.ID)
		assert.Equal(t, p.ID, o.ProjectID)
		assert.Equal(t, "test111", o.Name)
		assert.Equal(t, "Test Description Asdf", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		assert.NotZero(t, r.ID)
		assert.Equal(t, o.ID, r.TestID)
		assert.Equal(t, uint64(100), r.Count)
		assert.Equal(t, StatusOK, r.Status)
		assert.NotNil(t, r.CreatedAt)
		assert.NotNil(t, r.UpdatedAt)
		assert.Nil(t, r.DeletedAt)

		tid = o.ID
		pid = p.ID
		rid = r.ID

		cp := &Project{}
		err = db.First(cp, pid).Error
		assert.NoError(t, err)
		assert.Equal(t, p.Name, cp.Name)

		ct := &Test{}
		err = db.First(ct, tid).Error
		assert.NoError(t, err)
		assert.Equal(t, o.ProjectID, ct.ProjectID)
		assert.Equal(t, o.Name, ct.Name)
		assert.Equal(t, o.Description, ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Equal(t, o.Thresholds, ct.Thresholds)
		assert.Empty(t, ct.ThresholdsJSON)

		cr := &Run{}
		err = db.First(cr, rid).Error
		assert.NoError(t, err)
		assert.Equal(t, r.TestID, cr.TestID)
		assert.Equal(t, r.Status, cr.Status)
		assert.Empty(t, cr.StatusCodeDistJSON)
		assert.Empty(t, cr.ErrorDistJSON)
		assert.Equal(t, r.ErrorDistJSON, cr.ErrorDistJSON)
		assert.True(t, r.CreatedAt.Equal(cr.CreatedAt))
		assert.True(t, r.UpdatedAt.Equal(cr.CreatedAt))
	})
}
