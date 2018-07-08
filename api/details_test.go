package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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

	db.AutoMigrate(&model.Project{}, &model.Test{}, &model.Run{}, &model.Detail{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := &model.TestService{DB: db}
	ps := &model.ProjectService{DB: db}
	rs := &model.RunService{DB: db}
	ds := &model.DetailService{DB: db}
	// testAPI := &TestAPI{ts: ts}

	var projectID, runID, testID uint
	var pid, rid, tid string

	var httpTest *baloo.Client
	var echoServer *echo.Echo

	echoServer = echo.New()
	echoServer.Use(middleware.AddTrailingSlash())
	echoServer.Use(middleware.Logger())

	defer echoServer.Close()

	const basePath = "/projects"

	t.Run("Start API", func(t *testing.T) {
		projectGroup := echoServer.Group(basePath)
		SetupProjectAPI(projectGroup, ps)

		testsGroup := projectGroup.Group("/:pid/tests")
		SetupTestAPI(testsGroup, ts)

		runsGroup := testsGroup.Group("/:tid/runs")
		SetupRunAPI(runsGroup, rs)

		detailGroup := runsGroup.Group("/:rid/details")
		SetupDetailAPI(detailGroup, ds)

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

	t.Run("Create new project", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"name": " Test Project Name "}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(model.Project)
				json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.Equal(t, "testprojectname", p.Name)
				assert.Equal(t, "", p.Description)

				projectID = p.ID
				pid = strconv.FormatUint(uint64(projectID), 10)

				return nil
			}).
			Done()
	})

	t.Run("Create new test", func(t *testing.T) {
		httpTest.Post(basePath + "/" + pid + "/tests/").
			JSON(map[string]string{"name": " Test Name ", "description": "Test description"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "testname", tm.Name)
				assert.Equal(t, "Test description", tm.Description)

				testID = tm.ID
				tid = strconv.FormatUint(uint64(testID), 10)
				return nil
			}).
			Done()
	})

	t.Run("Create a run", func(t *testing.T) {
		httpTest.Post(basePath+"/"+pid+"/tests/"+tid+"/runs/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":1000,"total":5000000000,"average":123000000,"fastest":50000000,"slowest":234000000,"rps":6543.21}`).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				r := new(model.Run)
				json.NewDecoder(res.Body).Decode(r)

				assert.NoError(t, err)

				assert.NotZero(t, r.ID)
				assert.Equal(t, testID, r.TestID)
				assert.Equal(t, 1000, int(r.Count))
				assert.Equal(t, 5000*time.Millisecond, r.Total)
				assert.Equal(t, 123*time.Millisecond, r.Average)
				assert.Equal(t, 50*time.Millisecond, r.Fastest)
				assert.Equal(t, 234*time.Millisecond, r.Slowest)
				assert.Equal(t, 6543.21, r.Rps)

				rid = strconv.FormatUint(uint64(r.ID), 10)
				runID = r.ID

				return nil
			}).
			Done()
	})
}
