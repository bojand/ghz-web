package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var projectAPI = &ProjectAPI{}

func TestCreateProject(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, projectAPI.create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Created Project", rec.Body.String())
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
