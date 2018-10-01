package api

import (
	"time"

	"github.com/bojand/ghz-web/config"
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// @title ghz-web API
// @version 1.0
// @description ghz-web REST API

// Setup sets up the application API
func Setup(
	config *config.Config,
	info *config.Info,
	g *echo.Group,
	ps service.ProjectService,
	ts service.TestService,
	rs service.RunService,
	ds service.DetailService) {

	SetupInfoAPI(info, g)

	projectGroup := g.Group("/projects")
	SetupProjectAPI(projectGroup, ps)

	testsGroup := projectGroup.Group("/:pid/tests")
	SetupTestAPI(testsGroup, ts, rs)

	runsGroup := testsGroup.Group("/:tid/runs")
	SetupRunAPI(runsGroup, rs, ds)

	detailGroup := runsGroup.Group("/:rid/details")
	SetupDetailAPI(detailGroup, ds)

	SetupRawAPI(g, ps, ts, rs, ds)
}

// Model for common api objects
type Model struct {
	// The id
	ID uint `json:"id" example:"123"`

	// The creation time
	CreatedAt time.Time `json:"createdAt"`

	// The updated time
	UpdatedAt time.Time `json:"updatedAt"`

	// The deleted time
	DeletedAt *time.Time `json:"deletedAt"`
}

// Listable is list
type Listable struct {
	// The total number of items
	Total uint `json:"total" example:"10"`
}
