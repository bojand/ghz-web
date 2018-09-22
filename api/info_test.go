package api

import (
	"encoding/json"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

func TestInfoAPI(t *testing.T) {
	var httpTest *baloo.Client
	var echoServer *echo.Echo

	echoServer = echo.New()
	echoServer.Use(middleware.AddTrailingSlash())
	echoServer.Use(middleware.Logger())

	defer echoServer.Close()

	const basePath = "/api"

	info := &config.Info{
		Version:   "dev",
		GOVersion: runtime.Version(),
		StartTime: time.Now(),
	}

	t.Run("Start API", func(t *testing.T) {
		apiGroup := echoServer.Group(basePath)
		SetupInfoAPI(info, apiGroup)

		go func() {
			echoServer.Start("localhost:0")
		}()
	})

	t.Run("Sync to get the port", func(t *testing.T) {
		done := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			done <- true
			close(done)
		}()

		<-done
	})

	t.Run("Create http test", func(t *testing.T) {
		httpTest = baloo.New(echoServer.Listener.Addr().String())
	})

	t.Run("GET info", func(t *testing.T) {
		httpTest.Get(basePath + "/info/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				infoRes := new(InfoResponse)
				err := json.NewDecoder(res.Body).Decode(infoRes)

				assert.NoError(t, err)

				assert.Equal(t, runtime.Version(), infoRes.RuntimeVersion)
				assert.NotEmpty(t, infoRes.Uptime)
				assert.NotNil(t, infoRes.MemoryStats)

				return nil
			}).
			Done()
	})
}
