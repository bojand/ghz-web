package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// SetupProjectAPI sets up the API
func SetupProjectAPI(g *echo.Group, ps service.ProjectService, ts service.TestService) {
	api := &ProjectAPI{ps: ps, ts: ts}

	g.POST("/", api.create)
	g.GET("/:pid", api.get)
	g.PUT("/:pid", api.update)
	g.DELETE("/:pid", api.delete)

	testsGroup := g.Group("/:pid/tests")

	testsGroup.GET("", api.listTests)

	SetupTestAPI(testsGroup, ts)
}

// ProjectAPI provides the api
type ProjectAPI struct {
	ps service.ProjectService
	ts service.TestService
}

func (api *ProjectAPI) create(c echo.Context) error {
	p := new(model.Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := api.ps.Create(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, p)
}

func (api *ProjectAPI) get(c echo.Context) error {
	idparam := c.Param("pid")
	id, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	var p *model.Project

	if getByID {
		if p, err = api.ps.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	} else {
		if p, err = api.ps.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) update(c echo.Context) error {
	p := new(model.Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	uid := uint(id)
	p.ID = uid

	if err = api.ps.Update(p); gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

func (api *ProjectAPI) listTests(c echo.Context) error {
	idparam := c.Param("pid")
	pid, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	if !getByID {
		var p *model.Project
		if p, err = api.ps.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		pid = int(p.ID)
	}

	pageparam := c.QueryParam("page")
	page := 0
	if pageparam != "" {
		p, err := strconv.Atoi(pageparam)
		if err != nil {
			page = p
		}
	}

	limit := 20

	c.Logger().Infof("LIST TESTS. PID: %+v LIMIT: %+v PAGE: %+v", pid, limit, page)

	tests, err := api.ts.FindByProjectID(uint(pid), limit, page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return c.JSON(http.StatusOK, tests)
}
