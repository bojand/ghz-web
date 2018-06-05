package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// SetupRunAPI sets up the API
func SetupRunAPI(g *echo.Group) {
	api := &RunAPI{}

	g.POST("/", api.create)
	g.GET("/:rid", api.get)
	g.PUT("/:rid", api.update)
	g.DELETE("/:rid", api.delete)
}

// RunAPI provides the api
type RunAPI struct {
}

func (api *RunAPI) create(c echo.Context) error {
	return c.String(http.StatusCreated, "Create Run")
}

func (api *RunAPI) get(c echo.Context) error {
	return c.String(http.StatusOK, "Get Run")
}

func (api *RunAPI) update(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not Implemented")
}

func (api *RunAPI) delete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete Run")
}
