package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

func TestDetailAPI(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	conf, cerr := config.Read("../test/config1.toml")
	if cerr != nil {
		assert.FailNow(t, cerr.Error())
	}

	db.AutoMigrate(&model.Project{}, &model.Test{}, &model.Run{}, &model.Detail{},
		&model.Bucket{}, &model.LatencyDistribution{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := &model.TestService{DB: db}
	ps := &model.ProjectService{DB: db}
	rs := &model.RunService{DB: db}
	ds := &model.DetailService{DB: db, Config: &conf.Database}

	var projectID, testID, runID uint
	var pid, tid, rid string

	var httpTest *baloo.Client
	var echoServer *echo.Echo

	echoServer = echo.New()
	echoServer.Use(middleware.AddTrailingSlash())
	echoServer.Use(middleware.Logger())

	defer echoServer.Close()

	const basePath = "/projects"

	var run0data, run1data []byte

	t.Run("Start API", func(t *testing.T) {
		projectGroup := echoServer.Group(basePath)
		SetupProjectAPI(projectGroup, ps)

		testsGroup := projectGroup.Group("/:pid/tests")
		SetupTestAPI(testsGroup, ts, rs)

		runsGroup := testsGroup.Group("/:tid/runs")
		SetupRunAPI(runsGroup, rs, ds)

		detailGroup := runsGroup.Group("/:rid/details")
		SetupDetailAPI(detailGroup, ds)

		apiGroup := echoServer.Group("/api")
		SetupRawAPI(apiGroup, ps, ts, rs, ds)

		routes := echoServer.Routes()
		for _, r := range routes {
			index := strings.Index(r.Name, "ghz api:")
			if index >= 0 {
				desc := fmt.Sprintf("%+v %+v", r.Method, r.Path)
				fmt.Println(desc)
			}
		}

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

	t.Run("Read data file", func(t *testing.T) {
		jsonFile, err := os.Open("../test/run0.json")
		assert.NoError(t, err)
		defer jsonFile.Close()

		run0data, err = ioutil.ReadAll(jsonFile)

		assert.NoError(t, err)
		assert.NotNil(t, run0data)
		assert.True(t, len(run0data) > 0)
	})

	t.Run("Read data file 2", func(t *testing.T) {
		jsonFile, err := os.Open("../test/run1.json")
		assert.NoError(t, err)
		defer jsonFile.Close()

		run1data, err = ioutil.ReadAll(jsonFile)

		assert.NoError(t, err)
		assert.NotNil(t, run1data)
		assert.True(t, len(run1data) > 0)
	})

	t.Run("POST create raw data", func(t *testing.T) {
		var data map[string]interface{}

		err := json.Unmarshal(run0data, &data)

		assert.NoError(t, err)

		httpTest.Post("/api/raw/").
			JSON(data).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rr := new(RawResponse)
				json.NewDecoder(res.Body).Decode(rr)

				assert.NoError(t, err)

				assert.NotNil(t, rr.Project)
				assert.NotNil(t, rr.Test)
				assert.NotNil(t, rr.Run)
				assert.NotNil(t, rr.Details)
				assert.NotZero(t, rr.Details.Success)

				assert.NotZero(t, rr.Project.ID)
				assert.NotZero(t, rr.Test.ID)
				assert.NotZero(t, rr.Run.ID)

				testID = rr.Test.ID
				tid = strconv.FormatUint(uint64(testID), 10)

				projectID = rr.Project.ID
				pid = strconv.FormatUint(uint64(projectID), 10)

				runID = rr.Run.ID
				rid = strconv.FormatUint(uint64(runID), 10)

				return nil
			}).
			Done()
	})

	t.Run("GET list of details", func(t *testing.T) {
		httpTest.Get("/projects/" + pid + "/tests/" + tid + "/runs/" + rid + "/details/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				dl := new(DetailList)
				json.NewDecoder(res.Body).Decode(dl)

				assert.NoError(t, err)
				assert.Len(t, dl.Data, 20)
				assert.Equal(t, 937, int(dl.Total))

				return nil
			}).
			Done()
	})
}
