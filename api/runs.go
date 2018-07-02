package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// RunList response
type RunList struct {
	Total uint         `json:"total"`
	Data  []*model.Run `json:"data"`
}

// SetupRunAPI sets up the API
func SetupRunAPI(g *echo.Group, rs service.RunService) {
	api := &RunAPI{rs: rs}

	g.GET("/", api.listRuns).Name = "ghz api: list runs"
	g.POST("/", api.create).Name = "ghz api: create run"

	g.Use(api.populateRun)

	g.GET("/:rid/", api.get).Name = "ghz api: get run"
	g.PUT("/:rid/", api.update).Name = "ghz api: update run"
	g.DELETE("/:rid/", api.delete).Name = "ghz api: delete run"
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
	ro := c.Get("run")
	r, ok := ro.(*model.Run)

	if r == nil || !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No Run in context")
	}

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) update(c echo.Context) error {
	r := new(model.Run)

	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ro := c.Get("run")
	rm, ok := ro.(*model.Run)

	if rm == nil || !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No Run in context")
	}

	r.ID = rm.ID

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	r.TestID = t.ID

	var err error

	if err = api.rs.Update(r); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) listRuns(c echo.Context) error {
	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	tid := t.ID

	page := getPageParam(c)

	doSort, sort, order := getSortAndOrder(c)

	limit := uint(20)

	countCh := make(chan uint, 1)
	dataCh := make(chan []*model.Run, 1)
	errCh := make(chan error, 2)
	defer close(errCh)

	go func() {
		count, err := api.rs.Count(tid)
		errCh <- err
		countCh <- count
		close(countCh)
	}()

	var err error

	go func() {
		var runs []*model.Run
		err = nil
		if doSort {
			runs, err = api.rs.FindByTestIDSorted(tid, limit, page, sort, order)
		} else {
			runs, err = api.rs.FindByTestID(tid, limit, page)
		}
		errCh <- err
		dataCh <- runs
		close(dataCh)
	}()

	count, data, err1, err2 := <-countCh, <-dataCh, <-errCh, <-errCh

	if err1 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	rl := &RunList{Total: count, Data: data}

	return c.JSON(http.StatusOK, rl)
}

func (api *RunAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

func (api *RunAPI) getRun(c echo.Context) (*model.Run, error) {
	idparam := c.Param("rid")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	var r *model.Run

	if r, err = api.rs.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
		return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return r, nil
}

func (api *RunAPI) populateRun(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r, err := api.getRun(c)
		if err != nil {
			return err
		}

		c.Set("run", r)

		return next(c)
	}
}
