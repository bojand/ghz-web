package main

import (
	"net/http"

	"github.com/bojand/ghz-web/config"
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
}

// Start starts the app
func (app *Application) Start() {
	app.Server = echo.New()

	app.setupLogger()

	db, err := app.setupDatabase()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer db.Close()

	app.setupServer()

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

func (app *Application) setupDatabase() (*gorm.DB, error) {
	dbType := app.Config.Database.GetDialect()
	dbConn := app.Config.Database.GetConnectionString()

	app.Logger.Infof("Connecting DB. Type: %+v Connection string: %+v", dbType, dbConn)

	db, err := gorm.Open(dbType, dbConn)
	if err != nil {
		return nil, err
	}

	if app.Config.Log.Level == "debug" {
		db.LogMode(true)
	}

	// Migrate the schema
	db.AutoMigrate()

	// db.SetLogger(gorm.Logger{e.Logger})

	return db, nil
}

func (app *Application) setupServer() {
	app.Server.Use(middleware.RequestID())
	app.Server.Use(middleware.Logger())
	app.Server.Use(middleware.Recover())

	// apiGroup := e.Group("/api")

	// userDAO := model.UserService{DB: db}

	// api.Setup(apiGroup, &userDAO)

	app.Server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
