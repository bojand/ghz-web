package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/bojand/ghz-web/dao"
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

func (app *Application) testStuff() {
	// TEST STUFF

	pdao := &dao.ProjectService{DB: app.DB}

	project := &model.Project{Name: "Testproject2"}

	app.Logger.Infof("Project: %+v\n", project)

	err := pdao.Create(project)
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		app.Logger.Infof("Saved: %+v", project.ID)
	}

	t1 := &model.Test{
		Name: "test2",
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

	t2 := &model.Test{}
	err = app.DB.First(t2, "name = ?", "test3").Error
	if err != nil {
		app.Logger.Errorf("Error: %+v\n", err.Error())
	} else {
		str, _ := json.Marshal(t2)
		app.Logger.Infof("Found: %+v\n", t2)
		app.Logger.Infof("JSON: %+v\n", string(str))
	}
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
