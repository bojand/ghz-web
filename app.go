package main

import (
	"github.com/bojand/ghz-web/api"
	"github.com/bojand/ghz-web/config"
	"github.com/bojand/ghz-web/docs"
	"github.com/bojand/ghz-web/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"

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
	Info   *config.Info
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
	db.AutoMigrate(
		&model.Project{},
		&model.Test{},
		&model.Run{},
		&model.Detail{},
		&model.LatencyDistribution{},
		&model.Bucket{},
	)

	if app.Config.Database.GetDialect() == "sqlite3" {
		// for sqlite we need this for foreign key constraint
		db.Exec("PRAGMA foreign_keys = ON;")
	}

	app.DB = db

	return nil
}

func (app *Application) setupServer() {
	ps := model.ProjectService{DB: app.DB}
	ts := model.TestService{DB: app.DB}
	rs := model.RunService{DB: app.DB}
	ds := model.DetailService{DB: app.DB, Config: &app.Config.Database}

	docs.SwaggerInfo.Host = app.Config.Server.GetHostPort()
	docs.SwaggerInfo.BasePath = app.Config.Server.RootURL + "/api"

	s := app.Server

	s.Use(middleware.CORS())

	s.Pre(middleware.AddTrailingSlash())

	root := s.Group(app.Config.Server.RootURL)

	root.Use(middleware.RequestID())
	root.Use(middleware.Logger())
	root.Use(middleware.Recover())

	apiRoot := root.Group("/api")

	api.Setup(app.Config, app.Info, apiRoot, &ps, &ts, &rs, &ds)

	s.Static("/", "ui/dist").Name = "ghz api: static"

	// cannot work with trailing slashes
	root.GET("/docs/*", echoSwagger.WrapHandler, middleware.RemoveTrailingSlash())

	api.PrintRoutes(s)
}
