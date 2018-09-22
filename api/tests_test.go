package api

import (
	"encoding/json"
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

func TestTestAPI(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&model.Project{}, &model.Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ts := &model.TestService{DB: db}
	ps := &model.ProjectService{DB: db}
	runService := &model.RunService{DB: db}

	testAPI := &TestAPI{ts: ts}

	var projectID, projectID2, testID, testID2 uint
	var pid, pid2, pid3, tid string

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
		SetupTestAPI(testsGroup, ts, runService)

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

				pid2 = strconv.FormatUint(uint64(p.ID), 10)
				projectID2 = p.ID

				return nil
			}).
			Done()
	})

	t.Run("Create 3nd project", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"name": " Test Project Name Three", "description": "Three 3"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(model.Project)
				err = json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.Equal(t, "testprojectnamethree", p.Name)
				assert.Equal(t, "Three 3", p.Description)

				pid3 = strconv.FormatUint(uint64(p.ID), 10)

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

	t.Run("POST create another test", func(t *testing.T) {
		httpTest.Post(basePath + "/" + pid + "/tests/").
			JSON(map[string]string{"name": " Test Name Another", "description": "Test description", "status": "fail"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "testnameanother", tm.Name)
				assert.Equal(t, "Test description", tm.Description)
				assert.Equal(t, model.StatusFail, tm.Status)

				return nil
			}).
			Done()
	})

	t.Run("POST create test with thresholds", func(t *testing.T) {
		httpTest.Post(basePath+"/"+pid+"/tests/").
			AddHeader("Content-Type", "application/json; charset=UTF-8").
			BodyString(`{"name":"threshold Test","description":"a description","status":"fail","thresholds":{"median":{"status":"fail","threshold":10000},"mean":{"status":"ok","threshold":20000},"95th":{"status":"ok","threshold":30000}}}`).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "thresholdtest", tm.Name)
				assert.Equal(t, "a description", tm.Description)
				assert.Equal(t, model.StatusFail, tm.Status)

				expectedTH := map[model.Threshold]*model.ThresholdSetting{
					model.ThresholdMedian: &model.ThresholdSetting{Threshold: time.Duration(10000), Status: model.StatusFail},
					model.ThresholdMean:   &model.ThresholdSetting{Threshold: time.Duration(20000), Status: model.StatusOK},
					model.Threshold95th:   &model.ThresholdSetting{Threshold: time.Duration(30000), Status: model.StatusOK}}

				assert.Equal(t, expectedTH, tm.Thresholds)

				return nil
			}).
			Done()
	})

	t.Run("POST fail with same test name", func(t *testing.T) {
		httpTest.Post(basePath + "/" + pid + "/tests/").
			JSON(map[string]string{"name": " Test Name"}).
			Expect(t).
			Status(400).
			Type("json").
			Done()
	})

	t.Run("POST pass with same test name for project 2", func(t *testing.T) {
		httpTest.Post(basePath + "/" + pid2 + "/tests/").
			JSON(map[string]string{"name": " Test Name"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "testname", tm.Name)
				assert.Equal(t, projectID2, tm.ProjectID)
				assert.Equal(t, "", tm.Description)
				assert.Equal(t, model.StatusOK, tm.Status)

				testID2 = tm.ID

				return nil
			}).
			Done()
	})

	t.Run("GET id", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.Equal(t, testID, tm.ID)
				assert.Equal(t, "testname", tm.Name)
				assert.Equal(t, "Test description", tm.Description)

				return nil
			}).
			Done()
	})

	t.Run("GET name", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/testname/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.Equal(t, testID, tm.ID)
				assert.Equal(t, "testname", tm.Name)
				assert.Equal(t, "Test description", tm.Description)

				return nil
			}).
			Done()
	})

	t.Run("GET by name for project 2 should 200", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid2 + "/tests/testname/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.Equal(t, testID2, tm.ID)
				assert.Equal(t, "testname", tm.Name)
				assert.Equal(t, "", tm.Description)

				return nil
			}).
			Done()
	})

	t.Run("GET by unknown name should 404", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/testnamebgt/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("GET by unknown ID should 404", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/5454/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})

	t.Run("PUT update existing test", func(t *testing.T) {
		httpTest.Put(basePath + "/" + pid + "/tests/" + tid + "/").
			JSON(map[string]string{"name": "updatedtestname", "description": "updated test description"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.NoError(t, err)

				assert.NotZero(t, tm.ID)
				assert.Equal(t, "updatedtestname", tm.Name)
				assert.Equal(t, "updated test description", tm.Description)

				return nil
			}).
			Done()
	})

	t.Run("GET id verify update", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/tests/" + tid + "/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tm := new(model.Test)
				err = json.NewDecoder(res.Body).Decode(tm)

				assert.NoError(t, err)

				assert.Equal(t, testID, tm.ID)
				assert.Equal(t, "updatedtestname", tm.Name)
				assert.Equal(t, "updated test description", tm.Description)

				return nil
			}).
			Done()
	})

	t.Run("populateTest with unknown ID should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/156", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "156")

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := testAPI.populateTest(handler)
		err := popMW(c)

		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("populateTest with unknown test name should 404", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/asdfdsa", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "asdfdsa")

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := testAPI.populateTest(handler)
		err := popMW(c)

		if assert.Error(t, err) {
			assert.IsType(t, err, &echo.HTTPError{})
			httpErr := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
	})

	t.Run("populateTest with valid id should work", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/"+tid, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, tid)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := testAPI.populateTest(handler)
		err := popMW(c)

		if assert.NoError(t, err) {
			to := c.Get("test")
			assert.IsType(t, to, &model.Test{})
			tm := to.(*model.Test)
			assert.NotZero(t, tm.ID)
			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "updatedtestname", tm.Name)
			assert.Equal(t, "updated test description", tm.Description)
		}
	})

	t.Run("populateTest with valid name should work", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/"+pid+"/updatedtestname", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("pid", "tid")
		c.SetParamValues(pid, "updatedtestname")

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		popMW := testAPI.populateTest(handler)
		err := popMW(c)

		if assert.NoError(t, err) {
			to := c.Get("test")
			assert.IsType(t, to, &model.Test{})
			tm := to.(*model.Test)
			assert.NotZero(t, tm.ID)
			assert.Equal(t, testID, tm.ID)
			assert.Equal(t, "updatedtestname", tm.Name)
			assert.Equal(t, "updated test description", tm.Description)
		}
	})

	t.Run("GET /:pid/tests", func(t *testing.T) {
		// create sample tests
		for i := 0; i < 25; i++ {
			t := &model.Test{
				ProjectID: projectID2,
				Name:      "test" + strconv.FormatInt(int64(i), 10),
			}
			ts.Create(t)
		}

		httpTest.Get(basePath + "/" + pid2 + "/tests/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tl := new(TestList)
				err = json.NewDecoder(res.Body).Decode(tl)

				assert.NoError(t, err)
				assert.Len(t, tl.Data, 20)

				assert.Equal(t, 26, int(tl.Total))
				assert.NotZero(t, tl.Data[0].ID)
				assert.NotEmpty(t, tl.Data[0].Name)
				assert.NotZero(t, tl.Data[1].ID)
				assert.NotEmpty(t, tl.Data[1].Name)
				assert.NotZero(t, tl.Data[19].ID)
				assert.NotEmpty(t, tl.Data[19].Name)

				return nil
			}).
			Done()
	})

	t.Run("GET /:pid/tests?page=1", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid2 + "/tests/").
			SetQueryParams(map[string]string{"page": "1"}).
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tl := new(TestList)
				err = json.NewDecoder(res.Body).Decode(tl)

				assert.NoError(t, err)
				assert.Len(t, tl.Data, 6)

				assert.Equal(t, 26, int(tl.Total))
				assert.NotZero(t, tl.Data[0].ID)
				assert.NotEmpty(t, tl.Data[0].Name)
				assert.NotZero(t, tl.Data[1].ID)
				assert.NotEmpty(t, tl.Data[1].Name)
				assert.NotZero(t, tl.Data[4].ID)
				assert.NotEmpty(t, tl.Data[4].Name)

				return nil
			}).
			Done()
	})

	t.Run("GET /:pid/tests on empty project", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid3 + "/tests/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				tl := new(TestList)
				err = json.NewDecoder(res.Body).Decode(tl)

				assert.NoError(t, err)
				assert.Len(t, tl.Data, 0)

				assert.Equal(t, 0, int(tl.Total))

				return nil
			}).
			Done()
	})

	t.Run("DELETE /:id", func(t *testing.T) {
		httpTest.Delete(basePath + "/" + pid + "/tests/" + tid + "/").
			Expect(t).
			Status(501).
			Type("json").
			Done()
	})

	t.Run("DELETE should 404 on unknown id", func(t *testing.T) {
		httpTest.Delete(basePath + "/" + pid + "/tests/5354/").
			Expect(t).
			Status(404).
			Type("json").
			Done()
	})
}
