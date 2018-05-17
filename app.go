package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/bojand/ghz-web/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Application is the app
type Application struct {
	Config *config.Config
	Logger echo.Logger
	Server *echo.Echo
	DB     *gorm.DB
}

// Start starts the app
func (app *Application) Start() {
	app.Server = echo.New()

	app.setupLogger()

	err := app.setupDatabase()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer app.DB.Close()

	app.setupServer()

	app.testStuff()

	app.Logger.Fatal(app.Server.Start(app.Config.Server.GetHostPort()))
}

func (app *Application) setupLogger() {
	if app.Config.Log.Level == "debug" {
		app.Server.Logger.SetLevel(log.DEBUG)
	} else if app.Config.Log.Level == "info" {
		app.Server.Logger.SetLevel(log.INFO)
	} else if app.Config.Log.Level == "warn" {
		app.Server.Logger.SetLevel(log.WARN)
	} else if app.Config.Log.Level == "error" {
		app.Server.Logger.SetLevel(log.ERROR)
	} else {
		app.Server.Logger.SetLevel(log.OFF)
	}

	app.Logger = app.Server.Logger
}

func (app *Application) setupDatabase() error {
	dbType := app.Config.Database.GetDialect()
	dbConn := app.Config.Database.GetConnectionString()

	app.Logger.Infof("Connecting DB. Type: %+v Connection string: %+v", dbType, dbConn)

	db, err := gorm.Open(dbType, dbConn)
	if err != nil {
		return err
	}

	if app.Config.Log.Level == "debug" {
		db.LogMode(true)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Project{}, model.Test{})

	// db.SetLogger(gorm.Logger{e.Logger})

	app.DB = db

	return nil
}

func (app *Application) setupServer() {
	s := app.Server

	root := s.Group(app.Config.Server.RootURL)

	root.Use(middleware.RequestID())
	root.Use(middleware.Logger())
	root.Use(middleware.Recover())

	// userDAO := model.UserService{DB: db}

	// api.Setup(apiGroup, &userDAO)

	root.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}

//
// =====
//

func (app *Application) testStuff() {
	// TEST STUFF

	project := &model.Project{Name: "Testproject1"}

	tdao := &model.TestService{DB: app.DB}
	t1 := &model.Test{
		Project: project,
		Name:    "test1",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
		},
		Description: "test descroption 1",
	}

	err := tdao.Create(t1)
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t1.ID)
	}

	t2 := &model.Test{
		Project: project,
		Name:    "test2",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMean: &model.ThresholdSetting{Threshold: 1 * time.Millisecond, Status: model.StatusOK},
		},
		Description: "test descroption 2",
	}

	err = tdao.Create(t2)
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t2.ID)
	}

	tests, err := tdao.FindByProjectID(project.ID, -1, -1)

	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		str, _ := json.Marshal(tests)
		fmt.Printf("Found: %+v\n\n", tests)
		fmt.Printf("JSON: %s\n\n====\n\n", string(str))
	}

	/*pdao := &model.ProjectService{DB: app.DB}

	project := &model.Project{Name: "Testproject1"}

	app.Logger.Infof("Project: %+v\n", project)

	err := pdao.Create(project)
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", project.ID)
	}

	t1 := &model.Test{
		Project: *project,
		Name:    "test1",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
		},
		Description: "test descroption 2",
	}

	err = app.DB.Create(t1).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t1.ID)
	}

	t2 := &model.Test{
		Project: *project,
		Name:    "test2",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
		},
		Description: "test descroption 2",
	}

	err = app.DB.Create(t2).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t2.ID)
	}

	t3 := &model.Test{
		Project: *project,
		Name:    "test3",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
			model.Threshold95th:   &model.ThresholdSetting{Threshold: 3 * time.Millisecond, Status: model.StatusOK},
		},
		Description: "test descroption 3",
	}

	err = app.DB.Create(t3).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t3.ID)
	}

	t5 := &model.Test{
		ProjectID: project.ID,
		Name:      "test4",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
			model.ThresholdMean:   &model.ThresholdSetting{Threshold: 1 * time.Millisecond, Status: model.StatusOK},
		},
		Description: "test descroption 4",
	}

	err = app.DB.Create(t5).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t5.ID)
	}

	t6 := &model.Test{
		ProjectID: project.ID,
		Name:      "test5",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 2 * time.Millisecond, Status: model.StatusFail},
			model.ThresholdMean:   &model.ThresholdSetting{Threshold: 1 * time.Millisecond, Status: model.StatusOK},
			model.Threshold95th:   &model.ThresholdSetting{Threshold: 3 * time.Millisecond, Status: model.StatusOK},
		},
		Description: "test descroption 5",
	}

	err = app.DB.Create(t6).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t6.ID)
	}

	t7 := &model.Test{
		ProjectID: project.ID,
		Name:      "test6",
		Thresholds: map[model.Threshold]*model.ThresholdSetting{
			model.ThresholdMedian: &model.ThresholdSetting{Threshold: 3 * time.Millisecond, Status: model.StatusOK},
			model.ThresholdMean:   &model.ThresholdSetting{Threshold: 4 * time.Millisecond, Status: model.StatusOK},
			model.Threshold95th:   &model.ThresholdSetting{Threshold: 5 * time.Millisecond, Status: model.StatusOK},
		},
		Description: "test descroption 6",
	}

	err = app.DB.Create(t7).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t7.ID)
	}

	t7 = &model.Test{
		ProjectID:   project.ID,
		Name:        "test7",
		Description: "test descroption 7",
	}

	err = app.DB.Create(t7).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t7.ID)
	}

	t8 := &model.Test{
		ProjectID:   project.ID + 100,
		Name:        "test8",
		Description: "test descroption 8",
	}

	err = app.DB.Create(t8).Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", t8.ID)
	}

	// =====

	t4 := &model.Test{}
	err = app.DB.First(t4, "name = ?", "test2").Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		str, _ := json.Marshal(t4)
		fmt.Printf("Found: %+v\n\n", t4)
		fmt.Printf("JSON: %s\n", string(str))
	}

	tests := []model.Test{}
	err = app.DB.Limit(3).Order("name desc").Model(project).Related(&tests).Error

	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		str, _ := json.Marshal(tests)
		fmt.Printf("Found: %+v\n\n", tests)
		fmt.Printf("JSON: %s\n", string(str))
	}

	tests = []model.Test{}
	err = app.DB.Model(project).Related(&tests).Error

	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		str, _ := json.Marshal(tests)
		fmt.Printf("Found: %+v\n\n", tests)
		fmt.Printf("JSON: %s\n\n====\n\n", string(str))
	}
	*/
}
