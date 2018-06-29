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

	var projectID uint
	var pid string

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

			projectID = p.ID
			pid = strconv.FormatUint(uint64(projectID), 10)
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
}
