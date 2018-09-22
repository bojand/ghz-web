package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// getProject gets Project
func getProject(ps service.ProjectService, c echo.Context) (*model.Project, error) {
	idparam := c.Param("pid")
	id, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	var p *model.Project

	if getByID {
		if p, err = ps.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	} else {
		if p, err = ps.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return p, nil
}

func getTest(ts service.TestService, c echo.Context) (*model.Test, error) {
	idparam := c.Param("tid")
	id, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	var t *model.Test

	if getByID {
		if t, err = ts.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	} else {
		po := c.Get("project")
		p, ok := po.(*model.Project)

		var pid uint

		if p == nil || !ok {
			pidparam := c.Param("pid")
			pidint, err := strconv.Atoi(pidparam)

			if err != nil {
				return nil, echo.NewHTTPError(http.StatusNotFound, "No project in context or pid")
			}

			pid = uint(pidint)
		} else {
			pid = uint(p.ID)
		}

		if t, err = ts.FindByName(pid, idparam); gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return t, nil
}

func getRun(rs service.RunService, c echo.Context) (*model.Run, error) {
	idparam := c.Param("rid")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Invalid id")
	}

	var r *model.Run

	if r, err = rs.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
		return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return r, nil
}

// PrintRoutes prints routes in the server
func PrintRoutes(echoServer *echo.Echo) {
	routes := echoServer.Routes()
	for _, r := range routes {
		index := strings.Index(r.Name, "ghz api:")
		if index >= 0 {
			desc := fmt.Sprintf("%+v %+v", r.Method, r.Path)
			fmt.Println(desc)
		}
	}

}
