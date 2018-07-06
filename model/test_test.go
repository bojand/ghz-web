package model

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
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
		}, false, &Test{Status: StatusOK}},
		{"with error true but no fail on error setting", &Test{}, [4]time.Duration{
			milli1, milli2, milli3, milli4,
		}, true, &Test{Status: StatusOK}},
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
		})
	}
}

func TestThresholdSetting_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name     string
		in       string
		expected ThresholdSetting
	}{
		{"just status", `{"status":"ok"}`, ThresholdSetting{Status: StatusOK}},
		{"status and duration", `{"status":"ok","threshold":1000000}`, ThresholdSetting{Status: StatusOK, Threshold: milli1}},
		{"status and duration 2", `{"status":"fail","threshold":2000000}`, ThresholdSetting{Status: StatusFail, Threshold: milli2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual ThresholdSetting
			err := json.Unmarshal([]byte(tt.in), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTestModel_BeforeUpdate(t *testing.T) {
	var tests = []struct {
		name        string
		in          *Test
		expected    *Test
		expectError bool
	}{
		{"no name", &Test{}, &Test{}, true},
		{"with name", &Test{Name: " Test 2 "}, &Test{Name: " Test 2 "}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.BeforeUpdate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}

func TestTestModel_AfterFind(t *testing.T) {
	var tests = []struct {
		name        string
		in          *Test
		expected    *Test
		expectError bool
	}{
		{"empty", &Test{}, &Test{}, false},
		{"thresholds",
			&Test{ProjectID: 1, Name: "test4", Description: "Test Description", Status: StatusFail,
				ThresholdsJSON: `{"95th":{"status":"ok","threshold":4000000},"99th":{"status":"fail","threshold":5000000},"mean":{"status":"ok","threshold":2000000},"median":{"status":"ok","threshold":3000000}}`},
			&Test{ProjectID: 1, Name: "test4", Description: "Test Description", Status: StatusFail,
				Thresholds: map[Threshold]*ThresholdSetting{
					Threshold95th:   &ThresholdSetting{Threshold: milli4, Status: StatusOK},
					Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusFail},
					ThresholdMedian: &ThresholdSetting{Threshold: milli3, Status: StatusOK},
					ThresholdMean:   &ThresholdSetting{Threshold: milli2, Status: StatusOK},
				}},
			false},
	}

	for _, tt := range tests {
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

func TestTestModel_BeforeSave(t *testing.T) {
	var tests = []struct {
		name        string
		in          *Test
		expected    *Test
		expectError bool
	}{
		{"no thresholds no project", &Test{Name: "test1"}, &Test{Name: "test1"}, true},
		{"no thresholds project id 0", &Test{ProjectID: 0, Name: "test2"}, &Test{Name: "test2"}, true},
		{"no thresholds project id",
			&Test{ProjectID: 1, Name: "Test3", Description: " Description for test "},
			&Test{ProjectID: 1, Name: "test3", Description: "Description for test"}, false},
		{"thresholds",
			&Test{ProjectID: 1, Name: " Test 4 ", Description: " Test Description ", Status: StatusFail,
				Thresholds: map[Threshold]*ThresholdSetting{
					Threshold95th:   &ThresholdSetting{Threshold: milli4, Status: StatusOK},
					Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusFail},
					ThresholdMedian: &ThresholdSetting{Threshold: milli3, Status: StatusOK},
					ThresholdMean:   &ThresholdSetting{Threshold: milli2, Status: StatusOK},
				}},
			&Test{ProjectID: 1, Name: "test4", Description: "Test Description", Status: StatusFail,
				Thresholds: map[Threshold]*ThresholdSetting{
					Threshold95th:   &ThresholdSetting{Threshold: milli4, Status: StatusOK},
					Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusFail},
					ThresholdMedian: &ThresholdSetting{Threshold: milli3, Status: StatusOK},
					ThresholdMean:   &ThresholdSetting{Threshold: milli2, Status: StatusOK},
				},
				ThresholdsJSON: `{"95th":{"status":"ok","threshold":4000000},"99th":{"status":"fail","threshold":5000000},"mean":{"status":"ok","threshold":2000000},"median":{"status":"ok","threshold":3000000}}`},
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.BeforeSave(nil)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, tt.in)
		})
	}
}

func TestTestService_Create(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}
	var tid uint
	var pid uint
	var proj *Project

	t.Run("fail new without project", func(t *testing.T) {
		o := Test{
			Name:        "Test111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&o)

		assert.Error(t, err)
	})

	t.Run("new test and project", func(t *testing.T) {
		p := &Project{}
		o := Test{
			Project:     p,
			Name:        "Test 111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&o)

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

		tid = o.ID
		pid = p.ID
		proj = p

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
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("new test existing project", func(t *testing.T) {
		o := Test{
			Project:     proj,
			Name:        "Test 112 ",
			Description: "Test Description 2 ",
		}
		err := dao.Create(&o)

		assert.NoError(t, err)

		assert.NotZero(t, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.Equal(t, "test112", o.Name)
		assert.Equal(t, "Test Description 2", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		cp := &Project{}
		err = db.First(cp, pid).Error
		assert.NoError(t, err)
		assert.Equal(t, proj.Name, cp.Name)

		ct := &Test{}
		err = db.First(ct, o.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, o.ProjectID, ct.ProjectID)
		assert.Equal(t, o.Name, ct.Name)
		assert.Equal(t, o.Description, ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Equal(t, o.Thresholds, ct.Thresholds)
		assert.Empty(t, ct.ThresholdsJSON)
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("new test existing project ID", func(t *testing.T) {
		o := Test{
			Name:        "Test 113 ",
			Description: "Test Description 3 ",
		}
		o.ProjectID = pid

		err := dao.Create(&o)

		assert.NoError(t, err)
		assert.NotZero(t, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.Equal(t, "test113", o.Name)
		assert.Equal(t, "Test Description 3", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		cp := &Project{}
		err = db.First(cp, pid).Error
		assert.NoError(t, err)
		assert.Equal(t, proj.Name, cp.Name)

		ct := &Test{}
		err = db.First(ct, o.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, o.ProjectID, ct.ProjectID)
		assert.Equal(t, o.Name, ct.Name)
		assert.Equal(t, o.Description, ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Equal(t, o.Thresholds, ct.Thresholds)
		assert.Empty(t, ct.ThresholdsJSON)
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("fail new test non existing project ID", func(t *testing.T) {
		o := Test{
			Name:        "Test 114 ",
			Description: "Test Description 4 ",
		}
		o.ProjectID = 4321

		err := dao.Create(&o)

		assert.Error(t, err)

		cp := &Project{}
		err = db.First(cp, 4321).Error
		assert.Error(t, err)
	})

	t.Run("fail new test with same name", func(t *testing.T) {
		o := Test{
			Project:     proj,
			Name:        "Test 112 ",
			Description: "Test Description 4 ",
		}

		err := dao.Create(&o)

		assert.Error(t, err)
	})

	t.Run("new test with no name", func(t *testing.T) {
		o := Test{
			Project:     proj,
			Description: "Test Description 4 ",
		}

		err := dao.Create(&o)

		assert.NoError(t, err)
		assert.NotZero(t, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.NotEmpty(t, o.Name)
		assert.Equal(t, "Test Description 4", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		cp := &Project{}
		err = db.First(cp, pid).Error
		assert.NoError(t, err)
		assert.Equal(t, proj.Name, cp.Name)

		ct := &Test{}
		err = db.First(ct, o.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, o.ProjectID, ct.ProjectID)
		assert.Equal(t, o.Name, ct.Name)
		assert.Equal(t, o.Description, ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Equal(t, o.Thresholds, ct.Thresholds)
		assert.Empty(t, ct.ThresholdsJSON)
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})
}

func TestTestService_Update(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}
	var tid uint
	var pid uint

	t.Run("create new test and project", func(t *testing.T) {
		p := &Project{}
		o := Test{
			Project:     p,
			Name:        "Test 111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&o)

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

		tid = o.ID
		pid = p.ID

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
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("update", func(t *testing.T) {
		o := Test{
			ProjectID:   pid,
			Name:        "Test 222 ",
			Description: "Test Description 2 ",
			Status:      StatusFail,
			Thresholds: map[Threshold]*ThresholdSetting{
				Threshold95th:   &ThresholdSetting{Threshold: milli4, Status: StatusOK},
				Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusFail},
				ThresholdMedian: &ThresholdSetting{Threshold: milli3, Status: StatusOK},
				ThresholdMean:   &ThresholdSetting{Threshold: milli2, Status: StatusOK},
			},
		}
		o.ID = tid

		err := dao.Update(&o)

		assert.NoError(t, err)
		assert.Equal(t, tid, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.Equal(t, "test222", o.Name)
		assert.Equal(t, "Test Description 2", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		tid = o.ID

		ct := &Test{}
		err = db.First(ct, tid).Error
		assert.NoError(t, err)
		assert.Equal(t, o.ProjectID, ct.ProjectID)
		assert.Equal(t, o.Name, ct.Name)
		assert.Equal(t, o.Description, ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Equal(t, o.Thresholds, ct.Thresholds)
		assert.Empty(t, ct.ThresholdsJSON)
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.UpdatedAt))
	})

	t.Run("update with invalid pid", func(t *testing.T) {
		o := Test{
			ProjectID:   1234,
			Name:        "Test 333 ",
			Description: "Test Description 3 ",
			Status:      StatusFail,
			Thresholds: map[Threshold]*ThresholdSetting{
				Threshold95th:   &ThresholdSetting{Threshold: milli4, Status: StatusOK},
				Threshold99th:   &ThresholdSetting{Threshold: milli5, Status: StatusOK},
				ThresholdMedian: &ThresholdSetting{Threshold: milli3, Status: StatusOK},
				ThresholdMean:   &ThresholdSetting{Threshold: milli2, Status: StatusFail},
			},
		}
		o.ID = tid

		err := dao.Update(&o)

		assert.Error(t, err)

		ct := &Test{}
		err = db.First(ct, tid).Error
		assert.NoError(t, err)
		assert.Equal(t, pid, ct.ProjectID)
		assert.Equal(t, "test222", ct.Name)
		assert.Equal(t, "Test Description 2", ct.Description)
		assert.Equal(t, o.Status, ct.Status)
		assert.Empty(t, ct.ThresholdsJSON)
	})
}

func TestTestService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}
	var tid uint
	var pid uint

	t.Run("create new test and project", func(t *testing.T) {
		p := &Project{}
		o := Test{
			Project:     p,
			Name:        "Test 111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&o)

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

		tid = o.ID
		pid = p.ID

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
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("find valid", func(t *testing.T) {
		o, err := dao.FindByID(tid)

		assert.NoError(t, err)
		assert.Equal(t, tid, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.Equal(t, "test111", o.Name)
		assert.Equal(t, "Test Description Asdf", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)
	})

	t.Run("find invalid", func(t *testing.T) {
		o, err := dao.FindByID(123)

		assert.Error(t, err)
		assert.Nil(t, o)
	})
}

func TestTestService_FindByName(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}
	var tid uint
	var pid uint

	t.Run("create new test and project", func(t *testing.T) {
		p := &Project{}
		o := Test{
			Project:     p,
			Name:        "Test 1234 ",
			Description: "Test Description Foo ",
		}

		err := dao.Create(&o)

		assert.NoError(t, err)
		assert.NotZero(t, p.ID)
		assert.NotEmpty(t, p.Name)
		assert.Equal(t, "", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		assert.NotZero(t, o.ID)
		assert.Equal(t, p.ID, o.ProjectID)
		assert.Equal(t, "test1234", o.Name)
		assert.Equal(t, "Test Description Foo", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)

		tid = o.ID
		pid = p.ID

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
		assert.Equal(t, o.ThresholdsJSON, ct.ThresholdsJSON)
		assert.True(t, o.CreatedAt.Equal(ct.CreatedAt))
		assert.True(t, o.UpdatedAt.Equal(ct.CreatedAt))
	})

	t.Run("find valid", func(t *testing.T) {
		o, err := dao.FindByName("test1234")

		assert.NoError(t, err)
		assert.Equal(t, tid, o.ID)
		assert.Equal(t, pid, o.ProjectID)
		assert.Equal(t, "test1234", o.Name)
		assert.Equal(t, "Test Description Foo", o.Description)
		assert.NotNil(t, o.CreatedAt)
		assert.NotNil(t, o.UpdatedAt)
		assert.Nil(t, o.DeletedAt)
	})

	t.Run("find invalid", func(t *testing.T) {
		o, err := dao.FindByName("lorem")

		assert.Error(t, err)
		assert.Nil(t, o)
	})
}

func TestTestService_FindByProject(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}

	var pid1 uint
	var pid2 uint

	t.Run("create new tests and project 1", func(t *testing.T) {
		p := &Project{}

		o := Test{
			Project:     p,
			Name:        "Test 0",
			Description: "Test Description Foo 0",
		}
		err := dao.Create(&o)

		pid1 = p.ID

		assert.NoError(t, err)

		for n := 1; n <= 7; n++ {
			nStr := strconv.FormatInt(int64(n), 10)

			o = Test{
				ProjectID:   p.ID,
				Name:        "Test " + nStr,
				Description: "Test Description Foo " + nStr,
			}
			err := dao.Create(&o)

			assert.NoError(t, err)
		}
	})

	t.Run("create 2nd tests and project", func(t *testing.T) {
		p := &Project{}

		o := Test{
			Project:     p,
			Name:        "Test P2 0 ",
			Description: "Test Description Foo P2 0",
		}
		err := dao.Create(&o)

		assert.NoError(t, err)

		pid2 = p.ID

		for n := 1; n <= 9; n++ {
			nStr := strconv.FormatInt(int64(n), 10)

			o = Test{
				ProjectID:   p.ID,
				Name:        "Test P2 " + nStr,
				Description: "Test Description Foo P2 " + nStr,
			}
			err := dao.Create(&o)

			assert.NoError(t, err)
		}
	})

	t.Run("find for project 1", func(t *testing.T) {
		tests, err := dao.FindByProjectID(pid1, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, tests, 8)
	})

	t.Run("find for project 2", func(t *testing.T) {
		tests, err := dao.FindByProjectID(pid2, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, tests, 10)
	})

	t.Run("find for project 2 paged", func(t *testing.T) {
		tests, err := dao.FindByProjectID(pid2, 3, 0)

		assert.NoError(t, err)
		assert.Len(t, tests, 3)

		for i, to := range tests {
			nStr := strconv.FormatInt(int64(9-i), 10)
			assert.Equal(t, "testp2"+nStr, to.Name)
		}
	})

	t.Run("find for project 2 paged 2", func(t *testing.T) {
		tests, err := dao.FindByProjectID(pid2, 3, 1)

		assert.NoError(t, err)
		assert.Len(t, tests, 3)

		for i, to := range tests {
			nStr := strconv.FormatInt(int64(6-i), 10)
			assert.Equal(t, "testp2"+nStr, to.Name)
		}
	})

	t.Run("find invalid", func(t *testing.T) {
		tests, err := dao.FindByProjectID(123, 5, 0)

		assert.NoError(t, err)
		assert.Len(t, tests, 0)
	})
}

func TestTestService_Count(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}

	var pid1 uint
	var pid2 uint

	t.Run("create new tests and project 1", func(t *testing.T) {
		p := &Project{}

		o := Test{
			Project:     p,
			Name:        "Test 0",
			Description: "Test Description Foo 0",
		}
		err := dao.Create(&o)

		pid1 = p.ID

		assert.NoError(t, err)

		for n := 1; n <= 7; n++ {
			nStr := strconv.FormatInt(int64(n), 10)

			o = Test{
				ProjectID:   p.ID,
				Name:        "Test " + nStr,
				Description: "Test Description Foo " + nStr,
			}
			err := dao.Create(&o)

			assert.NoError(t, err)
		}
	})

	t.Run("create 2nd tests and project", func(t *testing.T) {
		p := &Project{}

		o := Test{
			Project:     p,
			Name:        "Test P2 0 ",
			Description: "Test Description Foo P2 0",
		}
		err := dao.Create(&o)

		assert.NoError(t, err)

		pid2 = p.ID

		for n := 1; n <= 9; n++ {
			nStr := strconv.FormatInt(int64(n), 10)

			o = Test{
				ProjectID:   p.ID,
				Name:        "Test P2 " + nStr,
				Description: "Test Description Foo P2 " + nStr,
			}
			err := dao.Create(&o)

			assert.NoError(t, err)
		}
	})

	t.Run("count project 1", func(t *testing.T) {
		count, err := dao.Count(pid1)

		assert.NoError(t, err)
		assert.Equal(t, uint(8), count)
	})

	t.Run("find for project 2", func(t *testing.T) {
		count, err := dao.Count(pid2)

		assert.NoError(t, err)
		assert.Equal(t, uint(10), count)
	})

	t.Run("find for project 3 unknown", func(t *testing.T) {
		count, err := dao.Count(321)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), count)
	})
}

func TestTestService_FindByProjectSorted(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := TestService{DB: db}

	var pid1 uint

	t.Run("create new tests and project 1", func(t *testing.T) {
		p := &Project{}

		o := Test{
			Project:     p,
			Name:        "Test 0",
			Description: "Test Description Foo 0",
		}
		err := dao.Create(&o)

		pid1 = p.ID

		assert.NoError(t, err)

		for n := 1; n <= 9; n++ {
			nStr := strconv.FormatInt(int64(n), 10)

			o = Test{
				ProjectID:   p.ID,
				Name:        "Test " + nStr,
				Description: "Test Description Foo " + nStr,
			}
			err := dao.Create(&o)

			assert.NoError(t, err)
		}
	})

	t.Run("find for project 1 id asc", func(t *testing.T) {
		tests, err := dao.FindByProjectIDSorted(pid1, 10, 0, "id", "asc")

		assert.NoError(t, err)
		assert.Len(t, tests, 10)

		assert.Equal(t, uint(1), tests[0].ID)
		assert.Equal(t, uint(10), tests[9].ID)
	})

	t.Run("find for project 1 id desc", func(t *testing.T) {
		tests, err := dao.FindByProjectIDSorted(pid1, 10, 0, "id", "desc")

		assert.NoError(t, err)
		assert.Len(t, tests, 10)

		assert.Equal(t, uint(10), tests[0].ID)
		assert.Equal(t, uint(1), tests[9].ID)
	})

	t.Run("find for project 1 name asc", func(t *testing.T) {
		tests, err := dao.FindByProjectIDSorted(pid1, 10, 0, "name", "asc")

		assert.NoError(t, err)
		assert.Len(t, tests, 10)

		assert.Equal(t, "test0", tests[0].Name)
		assert.Equal(t, "test9", tests[9].Name)
	})

	t.Run("find for project 1 name desc paged", func(t *testing.T) {
		tests, err := dao.FindByProjectIDSorted(pid1, 3, 1, "name", "desc")

		assert.NoError(t, err)
		assert.Len(t, tests, 3)

		assert.Equal(t, "test6", tests[0].Name)
		assert.Equal(t, "test4", tests[2].Name)
	})

	t.Run("find invalid", func(t *testing.T) {
		tests, err := dao.FindByProjectIDSorted(123, 5, 0, "id", "asc")

		assert.NoError(t, err)
		assert.Len(t, tests, 0)
	})
}
