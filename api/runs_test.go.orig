package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestRunAPI(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&model.Project{}, &model.Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := model.TestService{DB: db}
	ps := model.ProjectService{DB: db}
	rs := model.RunService{DB: db}
	projectAPI := &ProjectAPI{ps: &ps, ts: &ts}
	testAPI := &TestAPI{ts: &ts}
	runAPI := &RunAPI{rs: &rs}

	var projectID, testID uint
	var pid, pid2, tid, tid2 string
	var project, project2 *model.Project
	var test, test2 *model.Test

	t.Run("create new project", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Project Name "}`
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, projectAPI.create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			p := new(model.Project)
			err := json.Unmarshal(rec.Body.Bytes(), p)

			assert.NoError(t, err)

			assert.NotZero(t, p.ID)
			assert.Equal(t, "testprojectname", p.Name)
			assert.Equal(t, "", p.Description)

			project = p
			projectID = p.ID
			pid = strconv.FormatUint(uint64(projectID), 10)
		}
	})

	t.Run("create 2nd project", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Project Name Two","description":"Asdf"}`
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, projectAPI.create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			p := new(model.Project)
			err := json.Unmarshal(rec.Body.Bytes(), p)

			assert.NoError(t, err)

			assert.NotZero(t, p.ID)
			assert.Equal(t, "testprojectnametwo", p.Name)
			assert.Equal(t, "Asdf", p.Description)

			pid2 = strconv.FormatUint(uint64(p.ID), 10)
			project2 = p
		}
	})

	t.Run("Create new test", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Name ","description":"Test description"}`
		req := httptest.NewRequest(echo.POST, "/projects/"+pid+"/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)
		c.Set("project", project)

		if assert.NoError(t, testAPI.create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			tm := new(model.Test)
			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.NotZero(t, tm.ID)
			assert.Equal(t, "testname", tm.Name)
			assert.Equal(t, "Test description", tm.Description)

			testID = tm.ID
			tid = strconv.FormatUint(uint64(testID), 10)
			test = tm
		}
	})

	t.Run("Create new test 2", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Name 2 ","description":"Test description two"}`
		req := httptest.NewRequest(echo.POST, "/projects/"+pid+"/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)
		c.Set("project", project)

		if assert.NoError(t, testAPI.create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			tm := new(model.Test)
			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.NotZero(t, tm.ID)
			assert.Equal(t, "testname2", tm.Name)
			assert.Equal(t, "Test description two", tm.Description)

			testID = tm.ID
			tid2 = strconv.FormatUint(uint64(testID), 10)
			test2 = tm
		}
	})

	t.Run("POST Create new run", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.Error(t, runAPI.create(c)) {
			// assert.Equal(t, http.StatusNotFound, rec.Code)
			// assert.Equal(t, "Invalid id", rec.Body.String())
		}
	})

	t.Run("GET get existing run", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.Error(t, runAPI.get(c)) {
			// assert.Equal(t, http.StatusNotFound, rec.Code)
			// assert.Equal(t, "Invalid id", rec.Body.String())
		}
	})

	t.Run("PUT update run", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.Error(t, runAPI.update(c)) {
			// assert.Equal(t, http.StatusBadRequest, rec.Code)
			// assert.Equal(t, "Request body can't be empty", rec.Body.String())
		}
	})

	t.Run("DELETE run", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, runAPI.delete(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Delete Run", rec.Body.String())
		}
	})
}
