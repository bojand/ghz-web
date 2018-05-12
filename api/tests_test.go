package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var testAPI = &TestAPI{}

func TestCreateTest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, testAPI.create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Create Test", rec.Body.String())
	}
}

func TestGetTest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, testAPI.get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Get Test", rec.Body.String())
	}
}

func TestUpdateTest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, testAPI.update(c)) {
		assert.Equal(t, http.StatusNotImplemented, rec.Code)
		assert.Equal(t, "Not Implemented", rec.Body.String())
	}
}

func TestDeleteTest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, testAPI.delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Delete Test", rec.Body.String())
	}
}
