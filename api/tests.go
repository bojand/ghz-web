package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// SetupTestAPI sets up the API
func SetupTestAPI(g *echo.Group, ts service.TestService) {
	api := &TestAPI{ts: ts}

	g.POST("", api.create)
	g.GET("/:tid", api.get)
	g.PUT("/:tid", api.update)
	g.DELETE("/:tid", api.delete)

	runsGroup := g.Group("/:tid/runs")

	SetupRunAPI(runsGroup)
}

// TestAPI provides the api
type TestAPI struct {
	ts service.TestService
}

func (api *TestAPI) create(c echo.Context) error {
	idparam := c.Param("pid")
	pid, err := strconv.Atoi(idparam)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	t := new(model.Test)

	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	t.ProjectID = uint(pid)

	err = api.ts.Create(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

func (api *TestAPI) get(c echo.Context) error {
	idparam := c.Param("tid")
	id, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	var t *model.Test

	if getByID {
		if t, err = api.ts.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	} else {
		if t, err = api.ts.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

func (api *TestAPI) update(c echo.Context) error {
	t := new(model.Test)

	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	uid := uint(id)
	t.ID = uid

	if err = api.ts.Update(t); gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

func (api *TestAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
