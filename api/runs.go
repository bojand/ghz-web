package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// SetupRunAPI sets up the API
func SetupRunAPI(g *echo.Group) {
	api := &RunAPI{}

	g.POST("", api.create)
	g.GET("/:rid", api.get)
	g.PUT("/:rid", api.update)
	g.DELETE("/:rid", api.delete)
}

// RunAPI provides the api
type RunAPI struct {
	rs service.RunService
}

func (api *RunAPI) create(c echo.Context) error {
	r := new(model.Run)

	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	r.TestID = t.ID

	err := api.rs.Create(r)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, r)
}

func (api *RunAPI) get(c echo.Context) error {
	idparam := c.Param("rid")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	var r *model.Run

	if r, err = api.rs.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) update(c echo.Context) error {
	r := new(model.Run)

	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	uid := uint(id)
	r.ID = uid

	if err = api.rs.Update(r); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	r.TestID = t.ID

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) delete(c echo.Context) error {
	return c.String(http.StatusOK, "Delete Run")
}
