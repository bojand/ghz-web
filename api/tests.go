package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// SetupTestApi sets up the API
func SetupTestAPI(g *echo.Group) {
	api := &TestAPI{}

	g.POST("/", api.create)
	g.GET("/:id", api.get)
	g.PUT("/:id", api.update)
	g.PUT("/:id", api.delete)

	runsGroup := g.Group("/:id/runs")
	SetupRuntAPI(runsGroup)
}

// TestAPI provides the api
type TestAPI struct {
}

func (api *TestAPI) create(c echo.Context) error {
	return c.String(http.StatusCreated, "Create Test")
}

func (api *TestAPI) get(c echo.Context) error {
	return c.String(http.StatusOK, "Get Test")
}

func (api *TestAPI) update(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not Implemented")
}

func (api *TestAPI) delete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete Test")
}
