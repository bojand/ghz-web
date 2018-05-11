package main

import (
	"net/http"

	"github.com/bojand/ghz-web/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Application is the app
type Application struct {
	Config *config.Config
}

// Start starts the app
func (app *Application) Start() {
	db, err := gorm.Open(app.Config.DB.Type, app.Config.DB.GetConnectionString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate()

	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	// db.SetLogger(gorm.Logger{e.Logger})

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// apiGroup := e.Group("/api")

	// userDAO := model.UserService{DB: db}

	// api.Setup(apiGroup, &userDAO)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(app.Config.Srv.GetHostPort()))
}
