package api

import (
	"net/http"
	"time"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// RawRequest request to the create raw api
type RawRequest struct {
	Count               uint64                    `json:"count"`
	Total               time.Duration             `json:"total"`
	Average             time.Duration             `json:"average"`
	Fastest             time.Duration             `json:"fastest"`
	Slowest             time.Duration             `json:"slowest"`
	Rps                 float64                   `json:"rps"`
	ErrorDist           map[string]int            `json:"errorDistribution,omitempty"`
	StatusCodeDist      map[string]int            `json:"statusCodeDistribution,omitempty"`
	Details             []*model.Detail           `json:"details"`
	LatencyDistribution []*RawLatencyDistribution `json:"latencyDistribution"`
	Histogram           []*RawBucket              `json:"histogram"`
}

// RawLatencyDistribution holds latency distribution data
type RawLatencyDistribution struct {
	Percentage int           `json:"percentage"`
	Latency    time.Duration `json:"latency"`
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
	Project *model.Project  `json:"project"`
	Test    *model.Test     `json:"test"`
	Run     *model.Run      `json:"run"`
	Details *DetailsCreated `json:"details"`
}

// DetailsCreated summary of how many details got created and how many failed
type DetailsCreated struct {
	Success uint `json:"success"`
	Fail    uint `json:"fail"`
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

	g.POST("/projects/:pid/tests/:tid/runs/:rid/raw/", api.createRaw).Name = "ghz api: create raw"
	g.POST("/raw/", api.createNew).Name = "ghz api: create raw 2"
}

func (api *RawAPI) createRaw(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

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

	r := new(model.Run)
	r.TestID = t.ID
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

	median, nine5, nine9 := r.GetThresholdValues()
	hasErrors := r.HasErrors()

	t.SetStatus(rr.Average, median, nine5, nine9, rr.Fastest, rr.Slowest,
		rr.Rps, hasErrors)

	err = api.rs.Create(r)
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

	if errored > 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, rres)
	}

	return echo.NewHTTPError(http.StatusCreated, rres)
}
