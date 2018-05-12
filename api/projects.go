package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// SetupProjectAPI sets up the API
func SetupProjectAPI(g *echo.Group) {
	api := &ProjectAPI{}

	g.POST("/", api.create)
	g.GET("/:id", api.get)
	g.PUT("/:id", api.update)
	g.PUT("/:id", api.delete)

	testsGroup := g.Group("/:id/tests")
	SetupTestAPI(testsGroup)
}

// ProjectAPI provides the api
type ProjectAPI struct {
}

func (api *ProjectAPI) create(c echo.Context) error {
	return c.String(http.StatusCreated, "Created Project")
}

func (api *ProjectAPI) get(c echo.Context) error {
	return c.String(http.StatusOK, "Get Project")
}

func (api *ProjectAPI) update(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not Implemented")
}

func (api *ProjectAPI) delete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete Project")
}
