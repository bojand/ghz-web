package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

const dbName = "../test/api_test.db"

func TestProjectAPI(t *testing.T) {
	defer os.Remove(dbName)

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Project{}, &model.Test{})
	db.Exec("PRAGMA foreign_keys = ON;")

	ps := &model.ProjectService{DB: db}
	// projectAPI := &ProjectAPI{ps: ps}

	var projectID uint
	var pid string
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

	t.Run("POST create new", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"name": " Test Project Name "}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(Project)
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

	t.Run("POST new empty project", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(Project)
				err = json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.NotEmpty(t, p.Name)
				assert.Equal(t, "", p.Description)

				return nil
			}).
			Done()
	})

	t.Run("POST create new with just description", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"description": "asdf"}).
			Expect(t).
			Status(201).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(Project)
				err = json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.NotZero(t, p.ID)
				assert.NotEmpty(t, p.Name)
				assert.Equal(t, "asdf", p.Description)

				return nil
			}).
			Done()
	})

	t.Run("POST fail with same name", func(t *testing.T) {
		httpTest.Post(basePath + "/").
			JSON(map[string]string{"name": " Test Project Name"}).
			Expect(t).
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(Project)
				err = json.NewDecoder(res.Body).Decode(p)

				fmt.Printf("\n\n%+v\n\n", p)
				return nil
			}).
			Status(400).
			Type("json").
			Done()
	})

	t.Run("GET by id", func(t *testing.T) {
		httpTest.Get(basePath + "/" + pid + "/").
			Expect(t).
			Status(200).
			Type("json").
			AssertFunc(func(res *http.Response, req *http.Request) error {
				p := new(Project)
				err = json.NewDecoder(res.Body).Decode(p)

				assert.NoError(t, err)

				assert.Equal(t, projectID, p.ID)
				assert.Equal(t, "testprojectname", p.Name)
				assert.Equal(t, "", p.Description)

				return nil
			}).
			Done()
	})

	// t.Run("GET name", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/testprojectname/").
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			p := new(Project)
	// 			err = json.NewDecoder(res.Body).Decode(p)

	// 			assert.NoError(t, err)

	// 			assert.Equal(t, projectID, p.ID)
	// 			assert.Equal(t, "testprojectname", p.Name)
	// 			assert.Equal(t, "", p.Description)

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("GET 404 on unknown", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/tstprj/").
	// 		Expect(t).
	// 		Status(404).
	// 		Type("json").
	// 		Done()
	// })

	// t.Run("PUT /:id", func(t *testing.T) {
	// 	httpTest.Put(basePath + "/" + pid + "/").
	// 		JSON(map[string]string{"name": " Updated Project Name ", "description": "My project description!"}).
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			p := new(Project)
	// 			err = json.NewDecoder(res.Body).Decode(p)

	// 			assert.NoError(t, err)

	// 			assert.Equal(t, projectID, p.ID)
	// 			assert.Equal(t, "updatedprojectname", p.Name)
	// 			assert.Equal(t, "My project description!", p.Description)

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("PUT invalid id num", func(t *testing.T) {
	// 	httpTest.Put(basePath + "/12345/").
	// 		JSON(map[string]string{"name": " Updated Project Name 2", "description": " My project description two!"}).
	// 		Expect(t).
	// 		Status(404).
	// 		Type("json").
	// 		Done()
	// })

	// t.Run("PUT invalid id string", func(t *testing.T) {
	// 	httpTest.Put(basePath + "/updatedprojectnameasdf/").
	// 		JSON(map[string]string{"name": " Updated Project Name 2", "description": " My project description two!"}).
	// 		Expect(t).
	// 		Status(404).
	// 		Type("json").
	// 		Done()
	// })

	// t.Run("DELETE /:id", func(t *testing.T) {
	// 	httpTest.Delete(basePath + "/" + pid + "/").
	// 		Expect(t).
	// 		Status(501).
	// 		Type("json").
	// 		Done()
	// })

	// t.Run("DELETE 404 on unknown id", func(t *testing.T) {
	// 	httpTest.Delete(basePath + "/5432/").
	// 		Expect(t).
	// 		Status(404).
	// 		Type("json").
	// 		Done()
	// })

	// t.Run("GET /", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/").
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			pl := new(ProjectList)
	// 			err = json.NewDecoder(res.Body).Decode(pl)

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, uint(3), pl.Total)
	// 			assert.Len(t, pl.Data, 3)
	// 			assert.NotZero(t, pl.Data[0].ID)
	// 			assert.NotEmpty(t, pl.Data[0].Name)
	// 			assert.NotZero(t, pl.Data[1].ID)
	// 			assert.NotEmpty(t, pl.Data[1].Name)
	// 			assert.NotZero(t, pl.Data[2].ID)
	// 			assert.NotEmpty(t, pl.Data[2].Name)

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("GET /?sort=ID", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/").
	// 		SetQueryParams(map[string]string{"sort": "ID"}).
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			pl := new(ProjectList)
	// 			err = json.NewDecoder(res.Body).Decode(pl)

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, uint(3), pl.Total)
	// 			assert.Len(t, pl.Data, 3)
	// 			assert.NotZero(t, pl.Data[0].ID)
	// 			assert.NotEmpty(t, pl.Data[0].Name)
	// 			assert.NotZero(t, pl.Data[1].ID)
	// 			assert.NotEmpty(t, pl.Data[1].Name)
	// 			assert.NotZero(t, pl.Data[2].ID)
	// 			assert.NotEmpty(t, pl.Data[2].Name)
	// 			assert.True(t, pl.Data[0].ID < pl.Data[1].ID)
	// 			assert.True(t, pl.Data[1].ID < pl.Data[2].ID)

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("GET /?sort=ID&order=desc", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/").
	// 		SetQueryParams(map[string]string{"sort": "ID", "order": "desc"}).
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			pl := new(ProjectList)
	// 			err = json.NewDecoder(res.Body).Decode(pl)

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, uint(3), pl.Total)
	// 			assert.Len(t, pl.Data, 3)
	// 			assert.NotZero(t, pl.Data[0].ID)
	// 			assert.NotEmpty(t, pl.Data[0].Name)
	// 			assert.NotZero(t, pl.Data[1].ID)
	// 			assert.NotEmpty(t, pl.Data[1].Name)
	// 			assert.NotZero(t, pl.Data[2].ID)
	// 			assert.NotEmpty(t, pl.Data[2].Name)
	// 			assert.True(t, pl.Data[0].ID > pl.Data[1].ID)
	// 			assert.True(t, pl.Data[1].ID > pl.Data[2].ID)

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("GET /?sort=name&order=desc", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/").
	// 		SetQueryParams(map[string]string{"sort": "name", "order": "desc"}).
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			pl := new(ProjectList)
	// 			err = json.NewDecoder(res.Body).Decode(pl)

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, uint(3), pl.Total)
	// 			assert.Len(t, pl.Data, 3)
	// 			assert.NotZero(t, pl.Data[0].ID)
	// 			assert.NotEmpty(t, pl.Data[0].Name)
	// 			assert.NotZero(t, pl.Data[1].ID)
	// 			assert.NotEmpty(t, pl.Data[1].Name)
	// 			assert.NotZero(t, pl.Data[2].ID)
	// 			assert.NotEmpty(t, pl.Data[2].Name)
	// 			assert.Equal(t, 1, strings.Compare(pl.Data[0].Name, pl.Data[1].Name))
	// 			assert.Equal(t, 1, strings.Compare(pl.Data[1].Name, pl.Data[2].Name))

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("GET /?sort=name&order=asc", func(t *testing.T) {
	// 	httpTest.Get(basePath + "/").
	// 		SetQueryParams(map[string]string{"sort": "name", "order": "asc"}).
	// 		Expect(t).
	// 		Status(200).
	// 		Type("json").
	// 		AssertFunc(func(res *http.Response, req *http.Request) error {
	// 			pl := new(ProjectList)
	// 			err = json.NewDecoder(res.Body).Decode(pl)

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, uint(3), pl.Total)
	// 			assert.Len(t, pl.Data, 3)
	// 			assert.NotZero(t, pl.Data[0].ID)
	// 			assert.NotEmpty(t, pl.Data[0].Name)
	// 			assert.NotZero(t, pl.Data[1].ID)
	// 			assert.NotEmpty(t, pl.Data[1].Name)
	// 			assert.NotZero(t, pl.Data[2].ID)
	// 			assert.NotEmpty(t, pl.Data[2].Name)
	// 			assert.Equal(t, -1, strings.Compare(pl.Data[0].Name, pl.Data[1].Name))
	// 			assert.Equal(t, -1, strings.Compare(pl.Data[1].Name, pl.Data[2].Name))

	// 			return nil
	// 		}).
	// 		Done()
	// })

	// t.Run("populateProject with unknown project ID should 404", func(t *testing.T) {
	// 	e := echo.New()

	// 	req := httptest.NewRequest(echo.GET, "/56/1", strings.NewReader(""))
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()

	// 	c := e.NewContext(req, rec)
	// 	c.SetParamNames("pid")
	// 	c.SetParamValues("56")

	// 	handler := func(c echo.Context) error {
	// 		return c.String(http.StatusOK, "test")
	// 	}

	// 	popMW := projectAPI.populateProject(handler)
	// 	err := popMW(c)

	// 	if assert.Error(t, err) {
	// 		assert.IsType(t, err, &echo.HTTPError{})
	// 		httpErr := err.(*echo.HTTPError)
	// 		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	// 	}
	// })

	// t.Run("populateProject with unknown project Name should 404", func(t *testing.T) {
	// 	e := echo.New()

	// 	req := httptest.NewRequest(echo.GET, "/Asdfdsa/1", strings.NewReader(""))
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()

	// 	c := e.NewContext(req, rec)
	// 	c.SetParamNames("pid")
	// 	c.SetParamValues("Asdfdsa")

	// 	handler := func(c echo.Context) error {
	// 		return c.String(http.StatusOK, "test")
	// 	}

	// 	popMW := projectAPI.populateProject(handler)
	// 	err := popMW(c)

	// 	if assert.Error(t, err) {
	// 		assert.IsType(t, err, &echo.HTTPError{})
	// 		httpErr := err.(*echo.HTTPError)
	// 		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	// 	}
	// })

	// t.Run("populateProject with project should work", func(t *testing.T) {
	// 	e := echo.New()

	// 	pid := strconv.FormatUint(uint64(projectID), 10)

	// 	req := httptest.NewRequest(echo.GET, "/"+pid+"/1", strings.NewReader(""))
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()

	// 	c := e.NewContext(req, rec)
	// 	c.SetParamNames("pid")
	// 	c.SetParamValues(pid)

	// 	handler := func(c echo.Context) error {
	// 		return c.String(http.StatusOK, "test")
	// 	}

	// 	popMW := projectAPI.populateProject(handler)
	// 	err := popMW(c)

	// 	if assert.NoError(t, err) {
	// 		po := c.Get("project")
	// 		assert.IsType(t, po, &model.Project{})
	// 		p := po.(*model.Project)
	// 		assert.NotZero(t, p.ID)
	// 		assert.Equal(t, projectID, p.ID)
	// 		assert.Equal(t, "updatedprojectname", p.Name)
	// 		assert.Equal(t, "My project description!", p.Description)
	// 	}
	// })
}
