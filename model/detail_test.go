package model

import (
	"os"
	"testing"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDetailModel_BeforeSave(t *testing.T) {
	var tests = []struct {
		name        string
		in          *Detail
		expected    *Detail
		expectError bool
	}{
		{"no run", &Detail{Latency: 123.45}, &Detail{Latency: 123.45}, true},
		{"trim error", &Detail{RunID: 1, Latency: 123.45, Error: " network error "}, &Detail{RunID: 1, Latency: 123.45, Error: "network error", Status: "OK"}, false},
		{"trim status", &Detail{RunID: 1, Latency: 123.45, Error: " network error ", Status: " OK "}, &Detail{RunID: 1, Latency: 123.45, Error: "network error", Status: "OK"}, false},
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

func TestDetailModel_UnmarshalJSON(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02T15:04:05-0700", "2018-08-08T13:00:00-0300")
	assert.NoError(t, err)

	var tests = []struct {
		name        string
		in          string
		expected    *Detail
		expectError bool
	}{
		{"RFC3339",
			`{"timestamp":"2018-08-08T13:00:00.000000000-03:00","latency":1.23,"error":"","status":"OK"}`,
			&Detail{Timestamp: expectedTime, Latency: 1.23, Error: "", Status: "OK"},
			false},
		{"layoutISO2",
			`{"timestamp":"2018-08-08T13:00:00-0300","latency":1.23,"error":"","status":"OK"}`,
			&Detail{Timestamp: expectedTime, Latency: 1.23, Error: "", Status: "OK"},
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Detail
			err := d.UnmarshalJSON([]byte(tt.in))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, &d)
		})
	}
}

func TestDetailService_Create(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, pid, rid, did uint

	t.Run("fail new without run", func(t *testing.T) {
		o := &Detail{
			Latency: 123.45,
			Status:  "OK",
		}
		err := dao.Create(o)

		assert.Error(t, err)
	})

	t.Run("fail new with non existant run given id", func(t *testing.T) {
		o := &Detail{
			RunID:   5332,
			Latency: 123.45,
			Status:  "OK",
		}
		err := dao.Create(o)

		assert.Error(t, err)
	})

	t.Run("new detail with run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 123.45,
			Status:  "OK",
		}

		err := dao.Create(d)

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
		did = d.ID

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
		assert.True(t, r.UpdatedAt.Equal(cr.UpdatedAt))

		cd := &Detail{}
		err = db.First(cd, did).Error
		assert.NoError(t, err)
		assert.Equal(t, r.ID, cd.RunID)
		assert.Equal(t, d.ID, cd.RunID)
		assert.Equal(t, d.Latency, cd.Latency)
		assert.Equal(t, d.Status, cd.Status)
		assert.Equal(t, d.Error, cd.Error)
	})
}

func TestDetailService_CreateBatch(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, pid, rid, did uint

	t.Run("new detail with run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 123.45,
			Status:  "OK",
		}

		err := dao.Create(d)

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
		did = d.ID

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
		assert.True(t, r.UpdatedAt.Equal(cr.UpdatedAt))

		cd := &Detail{}
		err = db.First(cd, did).Error
		assert.NoError(t, err)
		assert.Equal(t, r.ID, cd.RunID)
		assert.Equal(t, d.ID, cd.RunID)
		assert.Equal(t, d.Latency, cd.Latency)
		assert.Equal(t, d.Status, cd.Status)
		assert.Equal(t, d.Error, cd.Error)
	})

	t.Run("create batch of details", func(t *testing.T) {
		M := 300
		s := make([]*Detail, M)

		for n := 0; n < M; n++ {
			nd := &Detail{
				RunID:   rid,
				Latency: 100.0 + float64(n),
				Status:  "OK",
			}

			s[n] = nd
		}

		created, errored := dao.CreateBatch(rid, s)

		assert.Equal(t, M, int(created))
		assert.Equal(t, 0, int(errored))
	})

	t.Run("fail create batch of details for unknown run id", func(t *testing.T) {
		M := 200
		s := make([]*Detail, M)

		for n := 0; n < M; n++ {
			nd := &Detail{
				RunID:   rid,
				Latency: 100.0 + float64(n),
				Status:  "OK",
			}

			s[n] = nd
		}

		created, errored := dao.CreateBatch(3213, s)

		assert.Equal(t, 0, int(created))
		assert.Equal(t, M, int(errored))
	})
}

func TestDetailService_Count(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, rid, rid2 uint

	t.Run("new details for run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 100.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		tid = o.ID
		rid = r.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nd := &Detail{
				Run:     r,
				Latency: 100.0 + float64(n),
				Status:  "OK",
			}

			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("new details for a different run", func(t *testing.T) {
		r := &Run{
			TestID:  tid,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		d := &Detail{
			Run:     r,
			Latency: 200.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid2 = r.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nd := &Detail{
				RunID:   r.ID,
				Latency: 200.0 + float64(n),
				Status:  "OK",
			}
			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("count run 1", func(t *testing.T) {
		count, err := dao.Count(rid)

		assert.NoError(t, err)
		assert.Equal(t, uint(10), count)
	})

	t.Run("find for run 2", func(t *testing.T) {
		count, err := dao.Count(rid2)

		assert.NoError(t, err)
		assert.Equal(t, uint(20), count)
	})

	t.Run("find for test 3 unknown", func(t *testing.T) {
		count, err := dao.Count(321)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), count)
	})
}

func TestDetailService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var rid, did uint

	t.Run("fail new without run", func(t *testing.T) {
		o := &Detail{
			Latency: 123.45,
			Status:  "OK",
		}
		err := dao.Create(o)

		assert.Error(t, err)
	})

	t.Run("fail new with non existant run given id", func(t *testing.T) {
		o := &Detail{
			RunID:   5332,
			Latency: 123.45,
			Status:  "OK",
		}
		err := dao.Create(o)

		assert.Error(t, err)
	})

	t.Run("new detail with run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 123.45,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid = r.ID
		did = d.ID
	})

	t.Run("find valid", func(t *testing.T) {
		cd, err := dao.FindByID(did)

		assert.NoError(t, err)
		assert.Equal(t, rid, cd.RunID)
		assert.Equal(t, did, cd.RunID)
		assert.Equal(t, 123.45, cd.Latency)
		assert.Equal(t, "OK", cd.Status)
		assert.Equal(t, "", cd.Error)
	})

	t.Run("find invalid", func(t *testing.T) {
		o, err := dao.FindByID(123)

		assert.Error(t, err)
		assert.Nil(t, o)
	})
}

func TestDetailService_FindByRunID(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, rid1, rid2 uint

	t.Run("new details for run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 100.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		tid = o.ID
		rid1 = r.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nd := &Detail{
				Run:     r,
				Latency: 100.0 + float64(n),
				Status:  "OK",
			}

			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("new details for a different run", func(t *testing.T) {
		r := &Run{
			TestID:  tid,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		d := &Detail{
			Run:     r,
			Latency: 200.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid2 = r.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nd := &Detail{
				RunID:   rid2,
				Latency: 200.0 + float64(n),
				Status:  "OK",
			}
			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("find for run 1", func(t *testing.T) {
		details, err := dao.FindByRunID(rid1, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, details, 10)
	})

	t.Run("find for run 2", func(t *testing.T) {
		details, err := dao.FindByRunID(rid2, 30, 0)

		assert.NoError(t, err)
		assert.Len(t, details, 20)
	})

	t.Run("find for run 2 paged", func(t *testing.T) {
		details, err := dao.FindByRunID(rid1, 5, 0)

		assert.NoError(t, err)
		assert.Len(t, details, 5)

		for i, rt := range details {
			assert.Equal(t, float64(109-i), rt.Latency)
		}
	})

	t.Run("find for run 2 paged 2", func(t *testing.T) {
		details, err := dao.FindByRunID(rid2, 5, 1)

		assert.NoError(t, err)
		assert.Len(t, details, 5)

		for i, tr := range details {
			assert.Equal(t, float64(214-i), tr.Latency)
		}
	})

	t.Run("find invalid", func(t *testing.T) {
		details, err := dao.FindByRunID(1235, 5, 0)

		assert.NoError(t, err)
		assert.Len(t, details, 0)
	})
}

func TestDetailService_FindByRunIDAll(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, rid1, rid2 uint

	t.Run("new details for run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 100.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		tid = o.ID
		rid1 = r.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nd := &Detail{
				Run:     r,
				Latency: 100.0 + float64(n),
				Status:  "OK",
			}

			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("new details for a different run", func(t *testing.T) {
		r := &Run{
			TestID:  tid,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		d := &Detail{
			Run:     r,
			Latency: 200.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid2 = r.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nd := &Detail{
				RunID:   rid2,
				Latency: 200.0 + float64(n),
				Status:  "OK",
			}
			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("find for run 1", func(t *testing.T) {
		details, err := dao.FindByRunIDAll(rid1)

		assert.NoError(t, err)
		assert.Len(t, details, 10)
	})

	t.Run("find for run 2", func(t *testing.T) {
		details, err := dao.FindByRunIDAll(rid2)

		assert.NoError(t, err)
		assert.Len(t, details, 20)
	})

	t.Run("find invalid", func(t *testing.T) {
		details, err := dao.FindByRunIDAll(1235)

		assert.NoError(t, err)
		assert.Len(t, details, 0)
	})
}

func TestDetailService_FindByRunIDSorted(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var tid, rid1, rid2 uint

	t.Run("new details for run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 100.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		tid = o.ID
		rid1 = r.ID

		// create more runs
		for n := 1; n < 10; n++ {
			nd := &Detail{
				Run:     r,
				Latency: 101.0 + float64(n),
				Status:  "OK",
			}

			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("new details for a different run", func(t *testing.T) {
		r := &Run{
			TestID:  tid,
			Count:   100,
			Total:   milli1000,
			Average: milli5,
			Fastest: milli1,
			Slowest: milli500,
		}

		d := &Detail{
			Run:     r,
			Latency: 200.0,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid2 = r.ID

		// create more runs
		for n := 1; n < 20; n++ {
			nd := &Detail{
				RunID:   rid2,
				Latency: 200.0 + float64(n),
				Status:  "OK",
			}
			err := dao.Create(nd)

			assert.NoError(t, err)
		}
	})

	t.Run("find for run 1 by id asc", func(t *testing.T) {
		details, err := dao.FindByRunIDSorted(rid1, 10, 0, "id", "asc")

		assert.NoError(t, err)
		assert.Len(t, details, 10)
		assert.Equal(t, uint(1), details[0].ID)
		assert.Equal(t, uint(10), details[9].ID)
	})

	t.Run("find for run 1 by id desc", func(t *testing.T) {
		details, err := dao.FindByRunIDSorted(rid1, 20, 0, "id", "desc")

		assert.NoError(t, err)
		assert.Len(t, details, 10)
		assert.Equal(t, uint(10), details[0].ID)
		assert.Equal(t, uint(1), details[9].ID)
	})

	t.Run("find for run 1 by id desc page 1", func(t *testing.T) {
		details, err := dao.FindByRunIDSorted(rid1, 5, 1, "id", "desc")

		assert.NoError(t, err)
		assert.Len(t, details, 5)
		assert.Equal(t, uint(5), details[0].ID)
		assert.Equal(t, uint(1), details[4].ID)
	})

	t.Run("find for run 2 by latency desc page 1", func(t *testing.T) {
		details, err := dao.FindByRunIDSorted(rid2, 5, 1, "latency", "desc")

		assert.NoError(t, err)
		assert.Len(t, details, 5)
		assert.Equal(t, 214.0, details[0].Latency)
		assert.Equal(t, 210.0, details[4].Latency)
	})

	t.Run("find for run 2 by latency asc page 1", func(t *testing.T) {
		details, err := dao.FindByRunIDSorted(rid2, 5, 1, "latency", "asc")

		assert.NoError(t, err)
		assert.Len(t, details, 5)
		assert.Equal(t, 205.0, details[0].Latency)
		assert.Equal(t, 209.0, details[4].Latency)
	})

	t.Run("error on invalid sort param", func(t *testing.T) {
		_, err := dao.FindByRunIDSorted(rid2, 0, 1, "asdf", "asc")

		assert.Error(t, err)
	})

	t.Run("error on invalid order param", func(t *testing.T) {
		_, err := dao.FindByRunIDSorted(rid2, 0, 1, "latency", "asce")

		assert.Error(t, err)
	})

	t.Run("0 for invalid run id", func(t *testing.T) {
		runs, err := dao.FindByRunIDSorted(1234, 5, 1, "latency", "desc")

		assert.NoError(t, err)
		assert.Len(t, runs, 0)
	})
}

func TestDetailService_Update(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &Test{}, &Run{}, &Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	dao := DetailService{DB: db, Config: &config.DBConfig{Type: "sqlite3"}}
	var rid, did uint

	t.Run("new detail with run, test and project", func(t *testing.T) {
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

		d := &Detail{
			Run:     r,
			Latency: 123.45,
			Status:  "OK",
		}

		err := dao.Create(d)

		assert.NoError(t, err)

		rid = r.ID
		did = d.ID
	})

	t.Run("update", func(t *testing.T) {
		d := &Detail{
			RunID:   rid,
			Latency: 234.56,
			Status:  "OK",
		}
		d.ID = did

		err := dao.Update(d)

		assert.NoError(t, err)

		cd := &Detail{}
		err = db.First(cd, did).Error
		assert.NoError(t, err)
		assert.Equal(t, rid, cd.RunID)
		assert.Equal(t, d.ID, cd.RunID)
		assert.Equal(t, d.Latency, cd.Latency)
		assert.Equal(t, d.Status, cd.Status)
		assert.Equal(t, d.Error, cd.Error)
	})

	t.Run("fail update on unknown id", func(t *testing.T) {
		d := &Detail{
			RunID:   rid,
			Latency: 234.56,
			Status:  "OK",
		}
		d.ID = 123

		err := dao.Update(d)

		assert.Error(t, err)
	})

	t.Run("fail update on unknown run id", func(t *testing.T) {
		d := &Detail{
			RunID:   1212,
			Latency: 234.56,
			Status:  "OK",
		}
		d.ID = did

		err := dao.Update(d)

		assert.Error(t, err)
	})
}
