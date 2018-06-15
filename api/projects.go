package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// ProjectList response
type ProjectList struct {
	Total uint             `json:"total"`
	Data  []*model.Project `json:"data"`
}

// TestList response
type TestList struct {
	Total uint          `json:"total"`
	Data  []*model.Test `json:"data"`
}

// SetupProjectAPI sets up the API
func SetupProjectAPI(g *echo.Group, ps service.ProjectService, ts service.TestService) {
	api := &ProjectAPI{ps: ps, ts: ts}

	g.GET("", api.listProjects)
	g.POST("", api.create)
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
	pidi, err := strconv.Atoi(idparam)
	pid := uint(pidi)
	getByID := true
	if err != nil {
		getByID = false
	}

	if !getByID {
		var p *model.Project
		if p, err = api.ps.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		pid = p.ID
	}

	page := getPageParam(c)

	doSort, sort, order := getSortAndOrder(c)

	limit := uint(20)

	countCh := make(chan uint, 1)
	dataCh := make(chan []*model.Test, 1)
	errCh := make(chan error, 2)
	defer close(errCh)

	go func() {
		count, err := api.ts.Count(pid)
		errCh <- err
		countCh <- count
		close(countCh)
	}()

	go func() {
		var tests []*model.Test
		err = nil
		if doSort {
			tests, err = api.ts.FindByProjectIDSorted(pid, limit, page, sort, order)
		} else {
			tests, err = api.ts.FindByProjectID(pid, limit, page)
		}
		errCh <- err
		dataCh <- tests
		close(dataCh)
	}()

	count, data, err1, err2 := <-countCh, <-dataCh, <-errCh, <-errCh

	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	tl := &TestList{Total: count, Data: data}

	return c.JSON(http.StatusOK, tl)
}

func (api *ProjectAPI) listProjects(c echo.Context) error {
	page := getPageParam(c)

	doSort, sort, order := getSortAndOrder(c)

	limit := uint(20)

	countCh := make(chan uint, 1)
	dataCh := make(chan []*model.Project, 1)
	errCh := make(chan error, 2)
	defer close(errCh)

	go func() {
		count, err := api.ps.Count()
		errCh <- err
		countCh <- count
		close(countCh)
	}()

	go func() {
		var projects []*model.Project
		var err error
		if doSort {
			projects, err = api.ps.ListSorted(limit, page, sort, order)
		} else {
			projects, err = api.ps.List(limit, page)
		}
		errCh <- err
		dataCh <- projects
		close(dataCh)
	}()

	count, data, err1, err2 := <-countCh, <-dataCh, <-errCh, <-errCh

	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	pl := &ProjectList{Total: count, Data: data}

	return c.JSON(http.StatusOK, pl)
}
