package api

import (
	"net/http"
	"time"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// Detail is a detail object
type Detail struct {
	Model

	// Run id
	RunID uint `json:"runID" example:"321"`

	// Timestamp for the detail
	Timestamp time.Time `json:"timestamp"`

	// Latency of the call
	Latency float64 `json:"latency" validate:"required"`

	// Error details
	Error string `json:"error"`

	// Status of the call
	Status string `json:"status"`
}

// DetailListRequest request
type DetailListRequest struct {
	// The property by which to sort the results
	Sort string `json:"sort" query:"sort" validate:"omitempty,oneof=id"`

	// The sort order
	Order string `json:"order" query:"order" validate:"omitempty,oneof=asc desc"`

	// The page to view
	Page uint `json:"page" query:"page"`
}

// DetailListResponse response holds a list of details
type DetailListResponse struct {
	Listable

	// List of detail objects
	Data []*Detail `json:"data"`
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

	if err := bindAndValidate(c, dlReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := dlReq.Page
	sort := dlReq.Sort
	order := "asc"

	doSort := false

	if sort != "" {
		doSort = true

		order = dlReq.Order
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

	resData := make([]*Detail, len(data))

	for i, d := range data {
		resData[i] = new(Detail)
		resData[i].ID = d.ID
		resData[i].CreatedAt = d.CreatedAt
		resData[i].UpdatedAt = d.UpdatedAt
		resData[i].DeletedAt = d.DeletedAt
		resData[i].RunID = d.RunID
		resData[i].Timestamp = d.Timestamp
		resData[i].Latency = d.Latency
		resData[i].Error = d.Error
		resData[i].Status = d.Status
	}

	pl := &DetailListResponse{Data: resData}
	pl.Total = count

	return c.JSON(http.StatusOK, pl)
}

func (api *DetailAPI) deleteAll(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
