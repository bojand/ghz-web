package api

import (
	"github.com/bojand/ghz-web/service"
	"github.com/labstack/echo"
)

func populateProject(ps service.ProjectService, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		p, err := getProject(ps, c)
		if err != nil {
			return err
		}

		c.Set("project", p)

		return next(c)
	}
}

func populateTest(ts service.TestService, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		t, err := getTest(ts, c)
		if err != nil {
			return err
		}

		c.Set("test", t)

		return next(c)
	}
}
