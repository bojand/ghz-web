package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	db.AutoMigrate(&model.Project{}, &model.Test{}, &model.Run{}, &model.Detail{},
		&model.Bucket{}, &model.LatencyDistribution{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := &model.TestService{DB: db}
	ps := &model.ProjectService{DB: db}
	rs := &model.RunService{DB: db}
	ds := &model.DetailService{DB: db}

	runAPI := &RunAPI{rs: rs, ds: ds}

	var projectID, testID, testID2, runID uint
	var pid, tid, tid2, rid string

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
		SetupTestAPI(testsGroup, ts, rs)

		runsGroup := testsGroup.Group("/:tid/runs")
		SetupRunAPI(runsGroup, rs, ds)

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
				err = json.NewDecoder(res.Body).Decode(p)

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
				err = json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.Equal(t, "testprojectnametwo", p.Name)
				assert.Equal(t, "Asdf", p.Description)

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
				err = json.NewDecoder(res.Body).Decode(tm)

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
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "testname2", tm.Name)
				assert.Equal(t, "Test description two", tm.Description)

				testID2 = tm.ID
				tid2 = strconv.FormatUint(uint64(testID2), 10)

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
				err = json.NewDecoder(res.Body).Decode(r)

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
				err = json.NewDecoder(res.Body).Decode(r)

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
				err = json.NewDecoder(res.Body).Decode(r)

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

	t.Run("populateRun with unknown ID should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+tid+"/"+"runs"+"/156", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("tid", "rid")
		c.SetParamValues(tid, "156")

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := runAPI.populateRun(handler)
		err := popMW(c)

		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("populateRun with valid id should work", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+tid+"/"+"runs"+"/"+rid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("tid", "rid")
		c.SetParamValues(tid, rid)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := runAPI.populateRun(handler)
		err := popMW(c)

		if assert.NoError(t, err) {
			ro := c.Get("run")
			assert.IsType(t, ro, &model.Run{})
			rm := ro.(*model.Run)
			assert.NotZero(t, rm.ID)
			assert.Equal(t, runID, rm.ID)
		}
	})

	t.Run("GET /:tid/runs should get empty list for test with no runs", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rl := new(RunList)
				json.NewDecoder(res.Body).Decode(rl)

				assert.NoError(t, err)
				assert.Len(t, rl.Data, 0)
				assert.Equal(t, 0, int(rl.Total))

				return nil
			}).
			Done()
	})

	t.Run("GET /:tid/runs", func(t *testing.T) {
		// create sample tests
		for i := 0; i < 25; i++ {
			nr := &model.Run{
				TestID:  testID2,
				Date:    time.Now(),
				Count:   200 + uint64(i),
				Total:   1000 * time.Millisecond,
				Average: time.Duration(5+i) * time.Millisecond,
				Fastest: 1 * time.Millisecond,
				Slowest: 500 * time.Millisecond,
				Rps:     float64(5000 + i),
			}
			err := rs.Create(nr)

			assert.NoError(t, err)
		}

		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rl := new(RunList)
				err = json.NewDecoder(res.Body).Decode(rl)

				assert.NoError(t, err)
				assert.Len(t, rl.Data, 20)

				assert.Equal(t, 25, int(rl.Total))
				assert.NotZero(t, rl.Data[0].ID)
				assert.NotZero(t, 200, rl.Data[0].Count)
				assert.NotZero(t, rl.Data[1].ID)
				assert.NotZero(t, 201, rl.Data[1].Count)
				assert.NotZero(t, rl.Data[19].ID)
				assert.NotZero(t, 219, rl.Data[19].Count)

				return nil
			}).
			Done()
	})

	t.Run("GET /:tid/runs/latest", func(t *testing.T) {

		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/latest/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				r := new(model.Run)
				err = json.NewDecoder(res.Body).Decode(r)

				assert.NoError(t, err)

				assert.NotZero(t, r.ID)
				assert.Equal(t, testID2, r.TestID)
				assert.Equal(t, 224, int(r.Count))
				assert.Equal(t, 1000*time.Millisecond, r.Total)
				assert.Equal(t, 29*time.Millisecond, r.Average)
				assert.Equal(t, 1*time.Millisecond, r.Fastest)
				assert.Equal(t, 500*time.Millisecond, r.Slowest)
				assert.Equal(t, 5024.0, r.Rps)

				return nil
			}).
			Done()
	})

	t.Run("GET /:tid/runs?sort=average&order=desc", func(t *testing.T) {

		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/").
			SetQueryParams(map[string]string{"sort": "average", "order": "desc"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rl := new(RunList)
				err = json.NewDecoder(res.Body).Decode(rl)

				assert.NoError(t, err)
				assert.Len(t, rl.Data, 20)

				assert.Equal(t, 25, int(rl.Total))
				assert.NotZero(t, rl.Data[0].ID)
				assert.NotZero(t, time.Duration(5+19)*time.Millisecond, rl.Data[0].Average)
				assert.NotZero(t, rl.Data[1].ID)
				assert.NotZero(t, time.Duration(5+18)*time.Millisecond, rl.Data[1].Average)
				assert.NotZero(t, rl.Data[19].ID)
				assert.NotZero(t, time.Duration(5+0)*time.Millisecond, rl.Data[19].Average)

				return nil
			}).
			Done()
	})

	t.Run("GET /:tid/runs?sort=average&order=desc&page=1", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/").
			SetQueryParams(map[string]string{"sort": "average", "order": "desc", "page": "1"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rl := new(RunList)
				err = json.NewDecoder(res.Body).Decode(rl)

				assert.NoError(t, err)
				assert.Len(t, rl.Data, 5)

				assert.Equal(t, 25, int(rl.Total))
				assert.NotZero(t, rl.Data[0].ID)
				assert.NotZero(t, time.Duration(5+24)*time.Millisecond, rl.Data[0].Average)
				assert.NotZero(t, rl.Data[1].ID)
				assert.NotZero(t, time.Duration(5+23)*time.Millisecond, rl.Data[1].Average)
				assert.NotZero(t, rl.Data[4].ID)
				assert.NotZero(t, time.Duration(5+20)*time.Millisecond, rl.Data[4].Average)

				return nil
			}).
			Done()
	})

	t.Run("GET /:tid/runs?sort=average&order=desc&page=2", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid2 + "/runs/").
			SetQueryParams(map[string]string{"sort": "average", "order": "desc", "page": "2"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				rl := new(RunList)
				err = json.NewDecoder(res.Body).Decode(rl)

				assert.NoError(t, err)
				assert.Len(t, rl.Data, 0)

				assert.Equal(t, 25, int(rl.Total))

				return nil
			}).
			Done()
	})

	t.Run("GET export unknown run", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/4343212/export/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("GET export unknown run without format", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/export/").
			Expect(t).
			Status(400).
			Type("json").
			Done()
	})

	t.Run("GET export unknown run invalid format", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/export/").
			SetQueryParams(map[string]string{"format": "foo"}).
			Expect(t).
			Status(400).
			Type("json").
			Done()
	})

	t.Run("GET export run json", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/export/").
			SetQueryParams(map[string]string{"format": "json"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				exportRes := new(JSONExportRespose)
				err = json.NewDecoder(res.Body).Decode(exportRes)

				assert.NoError(t, err)
				assert.Equal(t, 2000, int(exportRes.Count))
				assert.Equal(t, 6000*time.Millisecond, exportRes.Total)
				assert.Equal(t, 222*time.Millisecond, exportRes.Average)
				assert.Equal(t, 60*time.Millisecond, exportRes.Fastest)
				assert.Equal(t, 444*time.Millisecond, exportRes.Slowest)
				assert.Equal(t, 6666.66, exportRes.Rps)

				fmt.Printf("%+v", exportRes)

				return nil
			}).
			Done()
	})

	t.Run("GET export run json", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/export/").
			SetQueryParams(map[string]string{"format": "json"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				exportRes := new(JSONExportRespose)
				err = json.NewDecoder(res.Body).Decode(exportRes)

				assert.NoError(t, err)
				assert.Equal(t, 2000, int(exportRes.Count))
				assert.Equal(t, 6000*time.Millisecond, exportRes.Total)
				assert.Equal(t, 222*time.Millisecond, exportRes.Average)
				assert.Equal(t, 60*time.Millisecond, exportRes.Fastest)
				assert.Equal(t, 444*time.Millisecond, exportRes.Slowest)
				assert.Equal(t, 6666.66, exportRes.Rps)

				fmt.Printf("%+v", exportRes)

				return nil
			}).
			Done()
	})

	t.Run("GET export run json", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/runs/" + rid + "/export/").
			SetQueryParams(map[string]string{"format": "csv"}).
			Expect(t).
			Status(200).
			Type("text/csv").
			Done()
	})
}
