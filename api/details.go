package api

import (
	"net/http"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// DetailList response
type DetailList struct {
	Total uint            `json:"total"`
	Data  []*model.Detail `json:"data"`
}

// SetupDetailAPI sets up the API
func SetupDetailAPI(g *echo.Group, ds service.DetailService) {
	api := &DetailAPI{ds: ds}

	g.GET("/", api.listDetails).Name = "ghz api: list details"
	g.POST("/", api.create).Name = "ghz api: create details"
	g.DELETE("/", api.deleteAll).Name = "ghz api: delete all details"
}

// DetailAPI provides the api
type DetailAPI struct {
	ds service.DetailService
}

func (api *DetailAPI) listDetails(c echo.Context) error {
	ro := c.Get("run")
	r, ok := ro.(*model.Run)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	rid := r.ID

	page := getPageParam(c)

	doSort, sort, order := getSortAndOrder(c)

	limit := uint(20)

	countCh := make(chan uint, 1)
	dataCh := make(chan []*model.Detail, 1)
	errCh := make(chan error, 2)
	defer close(errCh)

	go func() {
		count, err := api.ds.Count(rid)
		errCh <- err
		countCh <- count
		close(countCh)
	}()

	go func() {
		var details []*model.Detail
		var err error
		if doSort {
			details, err = api.ds.FindByRunIDSorted(rid, limit, page, sort, order)
		} else {
			details, err = api.ds.FindByRunID(rid, limit, page)
		}
		errCh <- err
		dataCh <- details
		close(dataCh)
	}()

	count, data, err1, err2 := <-countCh, <-dataCh, <-errCh, <-errCh

	if err1 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err1.Error())
	}

	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err2.Error())
	}

	pl := &DetailList{Total: count, Data: data}

	return c.JSON(http.StatusOK, pl)
}

func (api *DetailAPI) create(c echo.Context) error {
	d := new(model.Detail)

	if err := c.Bind(d); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ro := c.Get("run")
	r, ok := ro.(*model.Run)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	d.RunID = r.ID

	err := api.ds.Create(d)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, r)
}

func (api *DetailAPI) deleteAll(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
