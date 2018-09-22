package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/labstack/echo"
)

// InfoResponse is the info response
type InfoResponse struct {
	Version        string            `json:"version"`
	RuntimeVersion string            `json:"runtimeVersion"`
	Uptime         string            `json:"uptime"`
	MemoryStats    *runtime.MemStats `json:"memoryStats"`
}

// SetupInfoAPI sets up the info endpoint
func SetupInfoAPI(info *config.Info, g *echo.Group) {
	g.GET("/info/", func(c echo.Context) error {
		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)

		ir := InfoResponse{
			Version:        info.Version,
			RuntimeVersion: info.GOVersion,
			Uptime:         time.Since(info.StartTime).String(),
			MemoryStats:    memStats,
		}
		return c.JSON(http.StatusOK, ir)
	}).Name = "ghz api: get info"
}
