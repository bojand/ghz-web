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

	t.Run("fail new with non existant test given id", func(t *testing.T) {
		o := Run{
			TestID:  5332,
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

	t.Run("new run and existing test id", func(t *testing.T) {
		r := &Run{
			TestID:  tid,
			Count:   200,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		err := dao.Create(r)

		assert.NotZero(t, r.ID)
		assert.Equal(t, tid, r.TestID)
		assert.Equal(t, uint64(200), r.Count)
		assert.Equal(t, StatusOK, r.Status)
		assert.NotNil(t, r.CreatedAt)
		assert.NotNil(t, r.UpdatedAt)
		assert.Nil(t, r.DeletedAt)

		cr := &Run{}
		err = db.First(cr, r.ID).Error
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

func TestRunService_Count(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := RunService{DB: db}
	var tid, tid2, pid uint

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

		tid = o.ID
		pid = p.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nr := &Run{
				TestID:  tid,
				Count:   100 + uint64(n),
				Total:   milli1000,
				Average: milli5,
				Fastest: milli1,
				Slowest: milli500,
			}
			err := dao.Create(nr)

			assert.NoError(t, err)
		}
	})

	t.Run("new runs and different test", func(t *testing.T) {
		o := &Test{
			ProjectID:   pid,
			Name:        "Test 222 ",
			Description: "Test Description 2 Asdf 2",
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

		tid2 = o.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nr := &Run{
				TestID:  tid2,
				Count:   200 + uint64(n),
				Total:   milli1000,
				Average: milli5,
				Fastest: milli1,
				Slowest: milli500,
			}
			err := dao.Create(nr)

			assert.NoError(t, err)
		}
	})

	t.Run("count test 1", func(t *testing.T) {
		count, err := dao.Count(tid)

		assert.NoError(t, err)
		assert.Equal(t, uint(10), count)
	})

	t.Run("find for test 2", func(t *testing.T) {
		count, err := dao.Count(tid2)

		assert.NoError(t, err)
		assert.Equal(t, uint(20), count)
	})

	t.Run("find for test 3 unknown", func(t *testing.T) {
		count, err := dao.Count(321)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), count)
	})
}

func TestRunService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := RunService{DB: db}
	var rid uint
	var tid uint
	var cr *Run

	t.Run("new run and test and project", func(t *testing.T) {
		p := &Project{}

		o := &Test{
			Project:     p,
			Name:        "Test 111 ",
			Description: "Test Description Asdf ",
		}

		cr = &Run{
			Test:    o,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		err := dao.Create(cr)

		assert.NoError(t, err)
		assert.NotZero(t, cr.ID)
		assert.NotZero(t, cr.TestID)
		assert.Equal(t, o.ID, cr.TestID)

		rid = cr.ID
		tid = cr.TestID
	})

	t.Run("find valid", func(t *testing.T) {
		o, err := dao.FindByID(rid)

		assert.NoError(t, err)
		assert.Equal(t, rid, o.ID)
		assert.Equal(t, tid, o.TestID)
		assert.Equal(t, cr.Status, o.Status)
		assert.Empty(t, o.StatusCodeDistJSON)
		assert.Empty(t, o.ErrorDistJSON)
		assert.Equal(t, cr.ErrorDistJSON, o.ErrorDistJSON)
		assert.True(t, cr.CreatedAt.Equal(o.CreatedAt))
		assert.True(t, cr.UpdatedAt.Equal(o.CreatedAt))
	})

	t.Run("find invalid", func(t *testing.T) {
		o, err := dao.FindByID(123)

		assert.Error(t, err)
		assert.Nil(t, o)
	})
}

func TestRunService_FindByTestID(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := RunService{DB: db}
	var tid1, tid2, pid uint

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

		tid1 = o.ID
		pid = p.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nr := &Run{
				TestID:  tid1,
				Count:   100 + uint64(n),
				Total:   milli1000,
				Average: milli5,
				Fastest: milli1,
				Slowest: milli500,
			}
			err := dao.Create(nr)

			assert.NoError(t, err)
		}
	})

	t.Run("new runs and different test", func(t *testing.T) {
		o := &Test{
			ProjectID:   pid,
			Name:        "Test 222 ",
			Description: "Test Description 2 Asdf 2",
		}

		r := &Run{
			Test:    o,
			Count:   210,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		err := dao.Create(r)

		assert.NoError(t, err)

		tid2 = o.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nr := &Run{
				TestID:  tid2,
				Count:   210 + uint64(n),
				Total:   milli1000,
				Average: milli5,
				Fastest: milli1,
				Slowest: milli500,
			}
			err := dao.Create(nr)

			assert.NoError(t, err)
		}
	})

	t.Run("find for test 1", func(t *testing.T) {
		runs, err := dao.FindByTestID(tid1, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, runs, 10)
	})

	t.Run("find for test 2", func(t *testing.T) {
		runs, err := dao.FindByTestID(tid2, 30, 0)

		assert.NoError(t, err)
		assert.Len(t, runs, 20)
	})

	t.Run("find for test 2 paged", func(t *testing.T) {
		runs, err := dao.FindByTestID(tid2, 5, 0)

		assert.NoError(t, err)
		assert.Len(t, runs, 5)

		for i, rt := range runs {
			assert.Equal(t, 229-i, int(rt.Count))
		}
	})

	t.Run("find for test 2 paged 2", func(t *testing.T) {
		runs, err := dao.FindByTestID(tid2, 5, 1)

		assert.NoError(t, err)
		assert.Len(t, runs, 5)

		for i, tr := range runs {
			assert.Equal(t, 224-i, int(tr.Count))
		}
	})

	t.Run("find invalid", func(t *testing.T) {
		runs, err := dao.FindByTestID(1235, 5, 0)

		assert.NoError(t, err)
		assert.Len(t, runs, 0)
	})
}
