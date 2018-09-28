package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/labstack/echo"
)

// MemoryInfo some memory stats
type MemoryInfo struct {
	// Bytes of allocated heap objects.
	Alloc uint64 `json:"allocated"`

	// Cumulative bytes allocated for heap objects.
	TotalAlloc uint64 `json:"totalAllocated"`

	// The total bytes of memory obtained from the OS.
	System uint64 `json:"system"`

	// The number of pointer lookups performed by the runtime.
	Lookups uint64 `json:"lookups"`

	// The cumulative count of heap objects allocated.
	// The number of live objects is Mallocs - Frees.
	Mallocs uint64 `json:"mallocs"`

	// The cumulative count of heap objects freed.
	Frees uint64 `json:"frees"`

	// The number of completed GC cycles.
	NumGC uint32 `json:"numGC"`
}

// InfoResponse is the info response
type InfoResponse struct {
	// Version of the application
	Version string `json:"version"`

	// Go runtime version
	RuntimeVersion string `json:"runtimeVersion"`

	// Uptime of the server
	Uptime string `json:"uptime"`

	// Memory info
	MemoryInfo *MemoryInfo `json:"memoryInfo,omitempty"`
}

// SetupInfoAPI sets up the info endpoint
func SetupInfoAPI(info *config.Info, g *echo.Group) {
	api := &InfoAPI{info: info}

	g.GET("/info/", api.getInfo).Name = "ghz api: get info"
}

// InfoAPI struct for info api
type InfoAPI struct {
	info *config.Info
}

// getInfo Gets the server info
// @Summary Gets the server info
// @Description Gets the server info
// @ID get-info
// @Produce json
// @Success 200 {object} api.InfoResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /info [get]
func (api *InfoAPI) getInfo(c echo.Context) error {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	ir := InfoResponse{
		Version:        api.info.Version,
		RuntimeVersion: api.info.GOVersion,
		Uptime:         time.Since(api.info.StartTime).String(),
		MemoryInfo: &MemoryInfo{
			Alloc:      memStats.Alloc,
			TotalAlloc: memStats.TotalAlloc,
			System:     memStats.Sys,
			Lookups:    memStats.Lookups,
			Mallocs:    memStats.Mallocs,
			Frees:      memStats.Frees,
			NumGC:      memStats.NumGC,
		},
	}

	return c.JSON(http.StatusOK, ir)
}
