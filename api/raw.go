package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// RawRequest request to the create raw api
type RawRequest struct {
	// Date of the test
	Date time.Time `json:"date"`

	// Options for the test
	Options *model.Options `json:"options,omitempty"`

	// Count is the number for calls
	Count uint64 `json:"count"`

	// Total duration of the test
	Total time.Duration `json:"total"`

	// Acerage duration of a call
	Average time.Duration `json:"average"`

	// Fastest call duration
	Fastest time.Duration `json:"fastest"`

	// Slowest call duration
	Slowest time.Duration `json:"slowest"`

	// Rps is the requests per second metric
	Rps float64 `json:"rps"`

	// ErrorDist is the error distribution
	ErrorDist map[string]int `json:"errorDistribution,omitempty"`

	// Status code distribution
	StatusCodeDist map[string]int `json:"statusCodeDistribution,omitempty"`

	// Details of all the calls
	Details []*model.Detail `json:"details"`

	// Latency distribution
	LatencyDistribution []*RawLatencyDistribution `json:"latencyDistribution"`

	// Histogram is the latency histrogram
	Histogram []*RawBucket `json:"histogram"`
}

const layoutISO string = "2006-01-02T15:04:05.666Z"
const layoutISO2 string = "2006-01-02T15:04:05-0700"

// UnmarshalJSON for RawRequest
func (rr *RawRequest) UnmarshalJSON(data []byte) error {
	type Alias RawRequest
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(rr),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	rr.Date, err = time.Parse(time.RFC3339, aux.Date)
	if err != nil {
		rr.Date, err = time.Parse(layoutISO, aux.Date)
	}
	if err != nil {
		rr.Date, err = time.Parse(layoutISO2, aux.Date)
	}

	return err
}

// MarshalJSON for RawRequest
func (rr RawRequest) MarshalJSON() ([]byte, error) {
	type Alias RawRequest
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  rr.Date.Format(time.RFC3339),
		Alias: (*Alias)(&rr),
	})
}

// RawLatencyDistribution holds latency distribution data
type RawLatencyDistribution struct {
	// Percentage of the dictribution
	Percentage int `json:"percentage"`

	// Latency of this distributoin
	Latency time.Duration `json:"latency"`
}

// RawBucket holds histogram data
type RawBucket struct {
	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

// RawResponse is the response to the raw endpoint
type RawResponse struct {
	// Created project
	Project *model.Project `json:"project"`

	// Created test
	Test *model.Test `json:"test"`

	// Created run
	Run *model.Run `json:"run"`

	// The summary of created details
	Details *DetailsCreated `json:"details"`
}

// DetailsCreated summary of how many details got created and how many failed
type DetailsCreated struct {
	// Number of successfully created detail objects
	Success uint `json:"success"`

	// Number of failed detail objects
	Fail uint `json:"fail"`
}

// RawAPI provides the api
type RawAPI struct {
	ps service.ProjectService
	ts service.TestService
	rs service.RunService
	ds service.DetailService
}

// SetupRawAPI sets up the API
func SetupRawAPI(g *echo.Group,
	ps service.ProjectService,
	ts service.TestService,
	rs service.RunService,
	ds service.DetailService) {

	api := &RawAPI{ps: ps, ts: ts, rs: rs, ds: ds}

	g.POST("/raw/", api.createNew).Name = "ghz api: create raw 2"

	g.Use(api.populateProject)
	g.Use(api.populateTest)

	g.POST("/projects/:pid/tests/:tid/raw/", api.createRaw).Name = "ghz api: create raw"
}

// Create raw result api
// @Summary Created raw results for given project and test
// @Description Created raw results for given project and test
// @ID post-create-raw
// @Accept  json
// @Produce json
// @Param RawRequest body api.RawRequest true "Raw request"
// @Success 200 {object} api.RawResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects/{pid}/tests/{tid}/raw/ [post]
func (api *RawAPI) createRaw(c echo.Context) error {
	po := c.Get("project")
	p, ok := po.(*model.Project)

	if p == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No project in context")
	}

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No test in context")
	}

	rr := new(RawRequest)

	if err := c.Bind(rr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return api.createBatch(c, rr, p, t)
}

// Create new raw result api
// @Summary Created new raw results
// @Description Created new raw results. Automatically creates prject, test and run.
// @ID post-create-raw-new
// @Accept  json
// @Produce json
// @Param RawRequest body api.RawRequest true "Raw request"
// @Success 200 {object} api.RawResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /projects/{pid}/tests/{tid}/raw/ [post]
func (api *RawAPI) createNew(c echo.Context) error {
	rr := new(RawRequest)

	if err := c.Bind(rr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p := new(model.Project)

	err := api.ps.Create(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	t := new(model.Test)
	t.ProjectID = p.ID

	err = api.ts.Create(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return api.createBatch(c, rr, p, t)
}

func (api *RawAPI) createBatch(c echo.Context, rr *RawRequest, p *model.Project, t *model.Test) error {
	r := new(model.Run)
	r.TestID = t.ID
	r.Date = rr.Date
	r.Options = rr.Options
	r.Count = rr.Count
	r.Total = rr.Total
	r.Average = rr.Average
	r.Fastest = rr.Fastest
	r.Slowest = rr.Slowest
	r.Rps = rr.Rps
	r.ErrorDist = rr.ErrorDist

	r.StatusCodeDist = rr.StatusCodeDist

	latencies := len(rr.LatencyDistribution)

	if latencies > 0 {
		r.LatencyDistribution = make([]*model.LatencyDistribution, latencies)
		for i, l := range rr.LatencyDistribution {
			r.LatencyDistribution[i] = new(model.LatencyDistribution)
			r.LatencyDistribution[i].Latency = l.Latency
			r.LatencyDistribution[i].Percentage = l.Percentage
		}
	}

	buckets := len(rr.Histogram)
	if buckets > 0 {
		r.Histogram = make([]*model.Bucket, buckets)
		for i, b := range rr.Histogram {
			r.Histogram[i] = new(model.Bucket)
			r.Histogram[i].Mark = b.Mark
			r.Histogram[i].Count = b.Count
			r.Histogram[i].Frequency = b.Frequency
		}
	}

	median, nine5 := r.GetThresholdValues()
	hasErrors := r.HasErrors()

	t.SetStatus(rr.Average, median, nine5, rr.Fastest, rr.Slowest,
		rr.Rps, hasErrors)

	err := api.rs.Create(r)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	created, errored := api.ds.CreateBatch(r.ID, rr.Details)

	rres := &RawResponse{
		Project: p,
		Test:    t,
		Run:     r,
		Details: &DetailsCreated{
			Success: created,
			Fail:    errored,
		},
	}

	if errored != uint(0) {
		return echo.NewHTTPError(http.StatusInternalServerError, rres)
	}

	return c.JSON(http.StatusCreated, rres)
}

func (api *RawAPI) populateProject(next echo.HandlerFunc) echo.HandlerFunc {
	return populateProject(api.ps, next)
}

func (api *RawAPI) populateTest(next echo.HandlerFunc) echo.HandlerFunc {
	return populateTest(api.ts, next)
}
