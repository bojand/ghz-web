package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bojand/ghz-web/dao"
	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/test"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const dbName = "../test/project_test.db"

func TestProjectAPI(t *testing.T) {

	defer os.Remove(dbName)

	err := test.SetupTestDatabase(dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := dao.ProjectService{DB: db}

	var projectAPI = &ProjectAPI{dao: &dao}
	var projectID uint

	t.Run("POST", func(t *testing.T) {
		e := echo.New()

		jsonStr := []byte(`{"name":" Test Project Name "}`)
		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(jsonStr))
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
			fmt.Println("projectID: " + strconv.FormatUint(uint64(projectID), 10))
		}
	})

	t.Run("GET id", func(t *testing.T) {
		e := echo.New()
		pid := strconv.FormatUint(uint64(projectID), 10)

		fmt.Println("pid: " + pid)
		req := httptest.NewRequest(echo.GET, "/"+pid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
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
		c.SetParamNames("id")
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
		c.SetParamNames("id")
		c.SetParamValues("tstprj")

		if assert.NoError(t, projectAPI.get(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("PUT /:id", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, projectAPI.update(c)) {
			assert.Equal(t, http.StatusNotImplemented, rec.Code)
			assert.Equal(t, "Not Implemented", rec.Body.String())
		}
	})

	t.Run("DELETE /:id", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, projectAPI.delete(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Delete Project", rec.Body.String())
		}
	})
}
