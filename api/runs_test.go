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

func TestRunAPI(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&model.Project{}, &model.Test{}, &model.Run{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := &model.TestService{DB: db}
	ps := &model.ProjectService{DB: db}
	rs := &model.RunService{DB: db}
	// projectAPI := &ProjectAPI{ps: ps}
	// testAPI := &TestAPI{ts: ts}
	// runAPI := &RunAPI{rs: rs}

	var projectID, projectID2, testID, testID2 uint
	var pid, pid2, tid, tid2, rid string
	// var project, project2 *model.Project
	// var test, test2 *model.Test

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

	t.Run("Create 2nd project", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"name": " Test Project Name Two", "description": "Asdf"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(model.Project)
				json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.Equal(t, "testprojectnametwo", p.Name)
				assert.Equal(t, "Asdf", p.Description)

				pid2 = strconv.FormatUint(uint64(p.ID), 10)
				projectID2 = p.ID

				return nil
			}).
			Done()
	})

	t.Run("POST create new test", func(t *testing.T) {
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

	t.Run("POST create 2nd test", func(t *testing.T) {
		httpTest.Post(basePath + "/" + pid + "/tests/").
			JSON(map[string]string{"name": " Test Name 2 ", "description": "Test description two"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "testname2", tm.Name)
				assert.Equal(t, "Test description two", tm.Description)

				testID2 = tm.ID
				tid2 = strconv.FormatUint(uint64(testID), 10)

				return nil
			}).
			Done()
	})

	t.Run("POST create a run", func(t *testing.T) {
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

				return nil
			}).
			Done()
	})

	t.Run("POST 404 on unknown project", func(t *testing.T) {
		httpTest.Post(basePath+"/4343/tests/"+tid+"/runs/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":2000,"total":6000000000,"average":222000000,"fastest":60000000,"slowest":444000000,"rps":6666.66}`).
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("POST 404 on unknown test", func(t *testing.T) {
		httpTest.Post(basePath+"/"+pid+"/tests/5435/runs/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":2000,"total":6000000000,"average":222000000,"fastest":60000000,"slowest":444000000,"rps":6666.66}`).
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("PUT update a run", func(t *testing.T) {
		httpTest.Put(basePath+"/"+pid+"/tests/"+tid+"/runs/"+rid+"/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":2000,"total":6000000000,"average":222000000,"fastest":60000000,"slowest":444000000,"rps":6666.66}`).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				r := new(model.Run)
				json.NewDecoder(res.Body).Decode(r)

				assert.NoError(t, err)

				assert.NotZero(t, r.ID)
				assert.Equal(t, testID, r.TestID)
				assert.Equal(t, 2000, int(r.Count))
				assert.Equal(t, 6000*time.Millisecond, r.Total)
				assert.Equal(t, 222*time.Millisecond, r.Average)
				assert.Equal(t, 60*time.Millisecond, r.Fastest)
				assert.Equal(t, 444*time.Millisecond, r.Slowest)
				assert.Equal(t, 6666.66, r.Rps)

				rid = strconv.FormatUint(uint64(r.ID), 10)

				return nil
			}).
			Done()
	})

	t.Run("PUT 404 on unknown project", func(t *testing.T) {
		httpTest.Put(basePath+"/4343/tests/"+tid+"/runs/"+rid+"/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":2000,"total":6000000000,"average":222000000,"fastest":60000000,"slowest":444000000,"rps":6666.66}`).
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("PUT 404 on unknown test", func(t *testing.T) {
		httpTest.Put(basePath+"/"+pid+"/tests/5435/runs/"+rid+"/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"count":2000,"total":6000000000,"average":222000000,"fastest":60000000,"slowest":444000000,"rps":6666.66}`).
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("GET updated run", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				r := new(model.Run)
				json.NewDecoder(res.Body).Decode(r)

				assert.NoError(t, err)

				assert.NotZero(t, r.ID)
				assert.Equal(t, testID, r.TestID)
				assert.Equal(t, 2000, int(r.Count))
				assert.Equal(t, 6000*time.Millisecond, r.Total)
				assert.Equal(t, 222*time.Millisecond, r.Average)
				assert.Equal(t, 60*time.Millisecond, r.Fastest)
				assert.Equal(t, 444*time.Millisecond, r.Slowest)
				assert.Equal(t, 6666.66, r.Rps)

				return nil
			}).
			Done()
	})

	t.Run("GET 404 on unknown id", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/5435/runs/5454/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("DELETE /:id", func(t *testing.T) {
		httpTest.Delete(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/").
			Expect(t).
			Status(501).
			Type("json").
			Done()
	})

	t.Run("DELETE should 404 on unknown id", func(t *testing.T) {
		httpTest.Delete(basePath + "/" + pid + "/tests/" + tid + "/runs/5454/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})
}
