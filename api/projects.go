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
	Listable

	// List of detail objects
	Data []*Project `json:"data"`
}

// Project represents a project
type Project struct {
	Model

	Name        string `json:"name"`
	Description string `json:"description"`
}

func formatProject(d *model.Project) *Project {
	p := new(Project)
	p.ID = d.ID
	p.CreatedAt = d.CreatedAt
	p.UpdatedAt = d.UpdatedAt
	p.DeletedAt = d.DeletedAt
	p.Name = d.Name
	p.Description = d.Description
	return p
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

// Create a project
// @Summary Create a project
// @Description Create a project
// @ID post-create-project
// @Produce json
// @Param Project body api.Project true "Project"
// @Success 200 {object} api.Project
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects [post]
func (api *ProjectAPI) create(c echo.Context) error {
	p := new(Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pm := &model.Project{Name: p.Name, Description: p.Description}

	err := api.ps.Create(pm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, formatProject(pm))
}

// Get a project
// @Summary Get a project
// @Description Get a project
// @ID get-project
// @Produce json
// @Param pid path int true "Project ID"
// @Success 200 {object} api.Project
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects/{pid} [get]
func (api *ProjectAPI) get(c echo.Context) error {
	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	return c.JSON(http.StatusOK, formatProject(p))
}

// Update a project
// @Summary Update a project
// @Description Update a project
// @ID put-update-project
// @Produce json
// @Param Project body api.Project true "Project"
// @Param pid path int true "Project ID"
// @Success 200 {object} api.Project
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects/{pid} [put]
func (api *ProjectAPI) update(c echo.Context) error {
	p := new(Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	po := c.Get("project")
	ep, ok := po.(*model.Project)

	if ep == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	pm := new(model.Project)
	pm.ID = ep.ID
	pm.Name = p.Name
	pm.Description = p.Description

	var err error

	if err = api.ps.Update(pm); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, formatProject(pm))
}

func (api *ProjectAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// Lists the projects
// @Summary Lists the projects
// @Description Lists the projects
// @ID get-list-projects
// @Produce json
// @Param page query integer false "The page to view"
// @Param order query string false "The sort order. Default: 'asc'"
// @Param sort query sring false "The property to sort by. Default: 'id'"
// @Success 200 {object} api.ProjectList
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects [get]
func (api *ProjectAPI) listProjects(c echo.Context) error {
	listReq := new(ListRequest)

	if err := bindAndValidate(c, listReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	doSort, sort, order, page := getSortAndOrderFromListRequest(listReq)

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

	resData := make([]*Project, len(data))

	for i, d := range data {
		resData[i] = formatProject(d)
	}

	pl := &ProjectList{Data: resData}
	pl.Total = count

	return c.JSON(http.StatusOK, pl)
}

func (api *ProjectAPI) populateProject(next echo.HandlerFunc) echo.HandlerFunc {
	return populateProject(api.ps, next)
}
