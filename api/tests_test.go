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
	// projectAPI := &ProjectAPI{ps: ps}
	// testAPI := &TestAPI{ts: ts}

	var projectID, testID uint
	var pid, pid2, tid string
	var project, project2 *model.Project

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

				project = p
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
				project2 = p

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
}
