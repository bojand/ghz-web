package api

import (
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

// Setup sets up the application API
func Setup(g *echo.Group, ps service.ProjectService, ts service.TestService) {
	projectGroup := g.Group("/projects")
	SetupProjectAPI(projectGroup, ps, ts)
}
