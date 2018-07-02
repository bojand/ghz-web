package api

import (
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// Setup sets up the application API
func Setup(g *echo.Group,
	ps service.ProjectService,
	ts service.TestService,
	rs service.RunService) {

	projectGroup := g.Group("/projects")
	SetupProjectAPI(projectGroup, ps)

	testsGroup := projectGroup.Group("/:pid/tests")
	SetupTestAPI(testsGroup, ts)

	runsGroup := testsGroup.Group("/:tid/runs")
	SetupRunAPI(runsGroup, rs)
}
