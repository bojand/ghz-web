package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bojand/ghz-web/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var projectAPI = &ProjectAPI{}

func TestCreateProject(t *testing.T) {
	// Setup
	e := echo.New()

	jsonStr := []byte(`{"name":" Test Project Name "}`)
	req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, projectAPI.create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		p := new(model.Project)
		err := json.Unmarshal(rec.Body.Bytes(), p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "testprojectname", p.Name)
		assert.Equal(t, "", p.Description)
	}
}

func TestGetProject(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, projectAPI.get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Get Project", rec.Body.String())
	}
}

func TestUpdateProject(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, projectAPI.update(c)) {
		assert.Equal(t, http.StatusNotImplemented, rec.Code)
		assert.Equal(t, "Not Implemented", rec.Body.String())
	}
}

func TestDeleteProject(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, projectAPI.delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Delete Project", rec.Body.String())
	}
}
