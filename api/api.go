package api

import (
	"github.com/labstack/echo"
)

// Setup sets up the application API
func Setup(g *echo.Group) {
	projectGroup := g.Group("/projects")
	SetupProjectAPI(projectGroup)
}
