package api

import (
	"net/http"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// TestList response
type TestList struct {
	Total uint          `json:"total"`
	Data  []*model.Test `json:"data"`
}

// SetupTestAPI sets up the API
func SetupTestAPI(g *echo.Group, ts service.TestService) {
	api := &TestAPI{ts: ts}

	g.GET("/", api.listTests).Name = "ghz api: list tests"
	g.POST("/", api.create).Name = "ghz api: create test"

	g.Use(api.populateTest)

	g.GET("/:tid/", api.get).Name = "ghz api: get test"
	g.PUT("/:tid/", api.update).Name = "ghz api: update test"
	g.DELETE("/:tid/", api.delete).Name = "ghz api: delete test"
}

// TestAPI provides the api
type TestAPI struct {
	ts service.TestService
}

func (api *TestAPI) create(c echo.Context) error {
	t := new(model.Test)
	var err error
	if err = c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	t.ProjectID = p.ID

	err = api.ts.Create(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

func (api *TestAPI) get(c echo.Context) error {
	to := c.Get("test")
	t, ok := to.(*model.Test)

	if t == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No Test in context")
	}

	return c.JSON(http.StatusOK, t)
}

func (api *TestAPI) update(c echo.Context) error {
	t := new(model.Test)

	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	to := c.Get("test")
	tm, ok := to.(*model.Test)

	if tm == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No Test in context")
	}

	t.ID = tm.ID

	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	t.ProjectID = p.ID

	var err error

	if err = api.ts.Update(t); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

func (api *TestAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

func (api *TestAPI) listTests(c echo.Context) error {
	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	pid := p.ID

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

	var err error

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
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	tl := &TestList{Total: count, Data: data}

	return c.JSON(http.StatusOK, tl)
}

func (api *TestAPI) populateTest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		t, err := getTest(api.ts, c)
		if err != nil {
			return err
		}

		c.Set("test", t)

		return next(c)
	}
}
