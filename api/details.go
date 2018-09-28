package api

import (
	"net/http"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// DetailListRequest request
type DetailListRequest struct {
	// The property by which to sort the results
	Sort string `json:"sort" query:"sort" validate:"oneof=id"`

	// The sort order
	Order string `json:"order" query:"order" validate:"oneof=asc desc"`

	// The page to view
	Page uint `json:"page" query:"page"`
}

// DetailListResponse response holds a list of details
type DetailListResponse struct {
	// The total number of details
	Total uint `json:"total"`

	// List of detail objects
	Data []*model.Detail `json:"data"`
}

// SetupDetailAPI sets up the API
func SetupDetailAPI(g *echo.Group, ds service.DetailService) {
	api := &DetailAPI{ds: ds}

	g.GET("/", api.listDetails).Name = "ghz api: list details"
	g.DELETE("/", api.deleteAll).Name = "ghz api: delete all details"
}

// DetailAPI provides the api
type DetailAPI struct {
	ds service.DetailService
}

// Lists the details for the run
// @Summary Lists the details for the specific run
// @Description Lists the details for the specific run.
// @ID get-list-details
// @Produce json
// @Param pid path int true "Project ID"
// @Param tid path int true "Test ID"
// @Param rid path int true "Run ID"
// @Param page query integer false "The page to view"
// @Param order query string false "The sort order. Default: 'asc'"
// @Param sort query sring false "The property to sort by. Default: 'id'"
// @Success 200 {object} api.DetailListResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects/{pid}/tests/{tid}/runs/{rid}/details [get]
func (api *DetailAPI) listDetails(c echo.Context) error {
	ro := c.Get("run")
	r, ok := ro.(*model.Run)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "No test in context")
	}

	rid := r.ID

	dlReq := new(DetailListRequest)

	if err := c.Bind(dlReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(dlReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := dlReq.Page
	sort := dlReq.Sort
	order := dlReq.Order

	doSort := false

	if sort != "" {
		doSort = true
	}

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

	pl := &DetailListResponse{Total: count, Data: data}

	return c.JSON(http.StatusOK, pl)
}

func (api *DetailAPI) deleteAll(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
