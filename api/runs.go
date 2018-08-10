package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alecthomas/template"
	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

const (
	csvTmpl = `
duration (ms),status,error{{ range $i, $v := . }}
{{ formatDuration .Latency 1000000.0 }},{{ .Status }},{{ .Error }}{{ end }}
`
)

var tmplFuncMap = template.FuncMap{
	"formatDuration": formatDuration,
}

func formatDuration(durationNano float64, div float64) string {
	return fmt.Sprintf("%4.2f", durationNano/div)
}

// RunList response
type RunList struct {
	Total uint         `json:"total"`
	Data  []*model.Run `json:"data"`
}

// DetailExport is detail for export
type DetailExport struct {
	Latency float64 `json:"latency" validate:"required"`
	Error   string  `json:"error"`
	Status  string  `json:"status"`
}

// LatencyExport holds latency distribution data
type LatencyExport struct {
	Percentage int           `json:"percentage"`
	Latency    time.Duration `json:"latency"`
}

// BucketExport holds histogram data
type BucketExport struct {
	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

// JSONExportRespose is the response to JSON export
type JSONExportRespose struct {
	Date    time.Time     `json:"date"`
	Count   uint64        `json:"count"`
	Total   time.Duration `json:"total"`
	Average time.Duration `json:"average"`
	Fastest time.Duration `json:"fastest"`
	Slowest time.Duration `json:"slowest"`
	Rps     float64       `json:"rps"`

	Options *model.Options `json:"options,omitempty"`

	LatencyDistribution []*LatencyExport `json:"latencyDistribution"`
	Histogram           []*BucketExport  `json:"histogram"`

	Details []*DetailExport `json:"details"`
}

// SetupRunAPI sets up the API
func SetupRunAPI(g *echo.Group, rs service.RunService, ds service.DetailService) {
	api := &RunAPI{rs: rs, ds: ds}

	g.POST("/", api.create).Name = "ghz api: create run"
	g.GET("/", api.listRuns).Name = "ghz api: list runs"
	g.GET("/latest/", api.getLatest).Name = "ghz api: list runs"

	g.Use(api.populateRun)

	g.GET("/:rid/", api.get).Name = "ghz api: get run"
	g.PUT("/:rid/", api.update).Name = "ghz api: update run"
	g.DELETE("/:rid/", api.delete).Name = "ghz api: delete run"
	g.GET("/:rid/export/", api.export).Name = "ghz api: export"
}

// RunAPI provides the api
type RunAPI struct {
	rs service.RunService
	ds service.DetailService
}

func (api *RunAPI) create(c echo.Context) error {
	r := new(model.Run)

	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No test in context")
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
		return echo.NewHTTPError(http.StatusInternalServerError, "No Run in context")
	}

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) getLatest(c echo.Context) error {
	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No test in context")
	}

	r, err := api.rs.FindLatest(t.ID)

	if gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, "No Run in context")
	}

	r.ID = rm.ID

	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No test in context")
	}

	r.TestID = t.ID

	var err error

	if err = api.rs.Update(r); gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, r)
}

func (api *RunAPI) listRuns(c echo.Context) error {
	to := c.Get("test")
	t, ok := to.(*model.Test)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No test in context")
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

	histogram := false
	latency := false

	popH := strings.ToLower(c.QueryParam("histogram"))
	if popH == "true" {
		histogram = true
	}

	popL := strings.ToLower(c.QueryParam("latency"))
	if popL == "true" {
		latency = true
	}

	popQ := strings.ToLower(c.QueryParam("populate"))
	if popQ == "true" {
		histogram = true
		latency = true
	}

	go func() {
		var runs []*model.Run
		err = nil
		if doSort {
			runs, err = api.rs.FindByTestIDSorted(tid, limit, page, sort, order, histogram, latency)
		} else {
			runs, err = api.rs.FindByTestID(tid, limit, page, true)
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

func (api *RunAPI) export(c echo.Context) error {
	ro := c.Get("run")
	rm, ok := ro.(*model.Run)

	if rm == nil || !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "No Run in context")
	}

	rid := rm.ID

	format := strings.ToLower(c.QueryParam("format"))
	if format != "csv" && format != "json" {
		return echo.NewHTTPError(http.StatusBadRequest, "Unsupported format: "+format)
	}

	details, err := api.ds.FindByRunIDAll(rid)

	if format == "csv" {
		outputTmpl := csvTmpl

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
		}

		buf := &bytes.Buffer{}
		templ := template.Must(template.New("tmpl").Funcs(tmplFuncMap).Parse(outputTmpl))
		if err := templ.Execute(buf, &details); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
		}

		return c.Blob(http.StatusOK, "text/csv", buf.Bytes())
	}

	jsonRes := JSONExportRespose{}
	jsonRes.Date = rm.Date
	jsonRes.Count = rm.Count
	jsonRes.Total = rm.Total
	jsonRes.Average = rm.Average
	jsonRes.Fastest = rm.Fastest
	jsonRes.Slowest = rm.Slowest
	jsonRes.Rps = rm.Rps

	jsonRes.Options = rm.Options

	jsonRes.LatencyDistribution = make([]*LatencyExport, len(rm.LatencyDistribution))
	for i, ld := range rm.LatencyDistribution {
		jsonRes.LatencyDistribution[i] = new(LatencyExport)
		jsonRes.LatencyDistribution[i].Percentage = ld.Percentage
		jsonRes.LatencyDistribution[i].Latency = ld.Latency
	}

	jsonRes.Histogram = make([]*BucketExport, len(rm.Histogram))
	for i, h := range rm.Histogram {
		jsonRes.Histogram[i] = new(BucketExport)
		jsonRes.Histogram[i].Mark = h.Mark
		jsonRes.Histogram[i].Count = h.Count
		jsonRes.Histogram[i].Frequency = h.Frequency
	}

	jsonRes.Details = make([]*DetailExport, len(details))
	for i, d := range details {
		jsonRes.Details[i] = new(DetailExport)
		jsonRes.Details[i].Error = d.Error
		jsonRes.Details[i].Latency = d.Latency
		jsonRes.Details[i].Status = d.Status
	}

	return c.JSONPretty(http.StatusOK, jsonRes, "  ")
}

func (api *RunAPI) populateRun(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r, err := getRun(api.rs, c)
		if err != nil {
			return err
		}

		c.Set("run", r)

		return next(c)
	}
}
