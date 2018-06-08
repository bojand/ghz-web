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
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const dbName = "../test/project_test.db"

func TestProjectAPI(t *testing.T) {

	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&model.Project{}, &model.Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ps := model.ProjectService{DB: db}
	ts := model.TestService{DB: db}
	projectAPI := &ProjectAPI{ps: &ps, ts: &ts}

	var projectID uint

	t.Run("POST create new", func(t *testing.T) {
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

			projectID = p.ID
		}
	})

	t.Run("POST create new empty", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{}`
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
			assert.NotEmpty(t, p.Name)
			assert.Equal(t, "", p.Description)
		}
	})

	t.Run("POST create new with just description", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"description":"asdf"}`
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
			assert.NotEmpty(t, p.Name)
			assert.Equal(t, "asdf", p.Description)
		}
	})

	t.Run("POST fail with same name", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Test Project Name"}`
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := projectAPI.create(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		}
	})

	t.Run("GET id", func(t *testing.T) {
		e := echo.New()
		pid := strconv.FormatUint(uint64(projectID), 10)

		req := httptest.NewRequest(echo.GET, "/"+pid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

		if assert.NoError(t, projectAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			p := new(model.Project)

			err := json.Unmarshal(rec.Body.Bytes(), p)

			assert.NoError(t, err)

			assert.Equal(t, projectID, p.ID)
			assert.Equal(t, "testprojectname", p.Name)
			assert.Equal(t, "", p.Description)
		}
	})

	t.Run("GET name", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/testprojectname", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues("testprojectname")

		if assert.NoError(t, projectAPI.get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			p := new(model.Project)

			err := json.Unmarshal(rec.Body.Bytes(), p)

			assert.NoError(t, err)

			assert.Equal(t, projectID, p.ID)
			assert.Equal(t, "testprojectname", p.Name)
			assert.Equal(t, "", p.Description)
		}
	})

	t.Run("GET 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/tstprj", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues("tstprj")

		err := projectAPI.get(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("PUT /:id", func(t *testing.T) {
		pid := strconv.FormatUint(uint64(projectID), 10)
		e := echo.New()

		jsonStr := `{"name":" Updated Project Name ","description":"My project description!"}`
		req := httptest.NewRequest(echo.PUT, "/"+pid, strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("pid")
		c.SetParamValues(pid)

		if assert.NoError(t, projectAPI.update(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			p := new(model.Project)
			err := json.Unmarshal(rec.Body.Bytes(), p)

			assert.NoError(t, err)

			assert.Equal(t, projectID, p.ID)
			assert.Equal(t, "updatedprojectname", p.Name)
			assert.Equal(t, "My project description!", p.Description)
		}
	})

	t.Run("PUT invalid id num", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Updated Project Name ","description":"My project description!"}`
		req := httptest.NewRequest(echo.PUT, "/12345", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues("12345")

		if assert.NoError(t, projectAPI.update(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("PUT invalid id string", func(t *testing.T) {
		e := echo.New()

		jsonStr := `{"name":" Updated Project Name 2","description":"My project description!"}`
		req := httptest.NewRequest(echo.PUT, "/updatedprojectnameasdf", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues("updatedprojectnameasdf")

		err := projectAPI.update(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("DELETE /:id", func(t *testing.T) {
		e := echo.New()
		pid := strconv.FormatUint(uint64(projectID), 10)

		req := httptest.NewRequest(echo.DELETE, "/"+pid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

		err := projectAPI.delete(c)
		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotImplemented, httpErr.Code)
		}
	})

	t.Run("GET /", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if assert.NoError(t, projectAPI.listProjects(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			ps := make([]model.Project, 0)

			err := json.Unmarshal(rec.Body.Bytes(), &ps)

			assert.NoError(t, err)
			assert.Len(t, ps, 3)
			assert.NotZero(t, ps[0].ID)
			assert.NotEmpty(t, ps[0].Name)
			assert.NotZero(t, ps[1].ID)
			assert.NotEmpty(t, ps[1].Name)
			assert.NotZero(t, ps[2].ID)
			assert.NotEmpty(t, ps[2].Name)
		}
	})

	t.Run("GET /:pid/tests", func(t *testing.T) {
		e := echo.New()
		pid := strconv.FormatUint(uint64(projectID), 10)

		// create sample tests
		for i := 0; i < 25; i++ {
			t := &model.Test{
				ProjectID: projectID,
				Name:      "test" + strconv.FormatInt(int64(i), 10),
			}
			ts.Create(t)
		}

		req := httptest.NewRequest(echo.GET, "/"+pid+"/tests", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

		if assert.NoError(t, projectAPI.listTests(c)) {

			assert.Equal(t, http.StatusOK, rec.Code)

			ps := make([]model.Project, 0)

			err := json.Unmarshal(rec.Body.Bytes(), &ps)

			assert.NoError(t, err)
			assert.Len(t, ps, 20)

			assert.NotZero(t, ps[0].ID)
			assert.NotEmpty(t, ps[0].Name)
			assert.NotZero(t, ps[1].ID)
			assert.NotEmpty(t, ps[1].Name)
			assert.NotZero(t, ps[19].ID)
			assert.NotEmpty(t, ps[19].Name)
		}
	})

	t.Run("GET /:pid/tests?page=1", func(t *testing.T) {
		e := echo.New()
		pid := strconv.FormatUint(uint64(projectID), 10)

		req := httptest.NewRequest(echo.GET, "/"+pid+"/tests?page=1", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid")
		c.SetParamValues(pid)

		if assert.NoError(t, projectAPI.listTests(c)) {

			assert.Equal(t, http.StatusOK, rec.Code)

			ps := make([]model.Project, 0)

			err := json.Unmarshal(rec.Body.Bytes(), &ps)

			assert.NoError(t, err)
			assert.Len(t, ps, 5)

			assert.NotZero(t, ps[0].ID)
			assert.NotEmpty(t, ps[0].Name)
			assert.NotZero(t, ps[1].ID)
			assert.NotEmpty(t, ps[1].Name)
			assert.NotZero(t, ps[4].ID)
			assert.NotEmpty(t, ps[4].Name)
		}
	})
}
