package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var runAPI = &RunAPI{}

func TestCreateRun(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, runAPI.create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Create Run", rec.Body.String())
	}
}

func TestGetRun(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, runAPI.get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Get Run", rec.Body.String())
	}
}

func TestUpdateRun(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, runAPI.update(c)) {
		assert.Equal(t, http.StatusNotImplemented, rec.Code)
		assert.Equal(t, "Not Implemented", rec.Body.String())
	}
}

func TestDeleteRun(t *testing.T) {
	// Setup
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
}
