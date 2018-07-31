package api

import (
	"net/http"

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

// SetupProjectAPI sets up the API
func SetupProjectAPI(g *echo.Group, ps service.ProjectService) {
	api := &ProjectAPI{ps: ps}

	g.GET("/", api.listProjects).Name = "ghz api: list projects"
	g.POST("/", api.create).Name = "ghz api: create project"

	g.Use(api.populateProject)

	g.GET("/:pid/", api.get).Name = "ghz api: get project"
	g.PUT("/:pid/", api.update).Name = "ghz api:  update project"
	g.DELETE("/:pid/", api.delete).Name = "ghz api: delete project"
}

// ProjectAPI provides the api
type ProjectAPI struct {
	ps service.ProjectService
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
	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) update(c echo.Context) error {
	p := new(model.Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	po := c.Get("project")
	ep, ok := po.(*model.Project)

	if ep == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	p.ID = ep.ID

	var err error

	if err = api.ps.Update(p); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	pl := &ProjectList{Total: count, Data: data}

	return c.JSON(http.StatusOK, pl)
}

func (api *ProjectAPI) populateProject(next echo.HandlerFunc) echo.HandlerFunc {
	return populateProject(api.ps, next)
}
