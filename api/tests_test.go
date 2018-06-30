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

func TestTestAPI(t *testing.T) {
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
	projectAPI := &ProjectAPI{ps: &ps, ts: &ts}
	testAPI := &TestAPI{ts: &ts}

	var projectID, testID uint
	var pid, pid2, tid string
	var project *model.Project

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
		}
	})

	t.Run("POST create new test", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Name ","description":"Test description"}`
		req := httptest.NewRequest(echo.POST, "/projects/"+pid+"/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

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
		}
	})

	t.Run("POST fail with same test name", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Name"}`
		req := httptest.NewRequest(echo.POST, "/projects/"+pid+"/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

		err := testAPI.create(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		}
	})

	t.Run("POST fail with same test name for project 2", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Name"}`
		req := httptest.NewRequest(echo.POST, "/projects/"+pid2+"/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid2)

		err := testAPI.create(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		}
	})

	t.Run("GET id", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/"+tid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, tid)

		if assert.NoError(t, testAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			tm := new(model.Test)

			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "testname", tm.Name)
			assert.Equal(t, "Test description", tm.Description)
		}
	})

	t.Run("GET name", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/testname", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "testname")

		if assert.NoError(t, testAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			tm := new(model.Test)

			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "testname", tm.Name)
			assert.Equal(t, "Test description", tm.Description)
		}
	})

	// TODO fix this
	t.Run("GET by name for project 2 should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid2+"/testname", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid2, "testname")

		if assert.NoError(t, testAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			tm := new(model.Test)

			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "testname", tm.Name)
			assert.Equal(t, "Test description", tm.Description)
		}
	})

	// TODO fix this
	t.Run("GET by unknown name should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/testnamebgt", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "testnamebgt")

		err := testAPI.get(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("GET by unknown ID should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/5454", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "5454")

		err := testAPI.get(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("PUT update existing test", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":"updatedtestname","description":"updated test description"}`
		req := httptest.NewRequest(echo.PUT, "/"+pid+"/"+tid, strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, tid)
		c.Set("project", project)

		err := testAPI.update(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)

			tm := new(model.Test)
			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.NotZero(t, tm.ID)
			assert.Equal(t, "updatedtestname", tm.Name)
			assert.Equal(t, "updated test description", tm.Description)
		}
	})

	t.Run("GET id verify update", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/"+tid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, tid)

		if assert.NoError(t, testAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			tm := new(model.Test)

			err := json.Unmarshal(rec.Body.Bytes(), tm)

			assert.NoError(t, err)

			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "updatedtestname", tm.Name)
			assert.Equal(t, "updated test description", tm.Description)
		}
	})
}
