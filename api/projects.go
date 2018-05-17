package api

import (
	"net/http"
	"strconv"

	"github.com/bojand/ghz-web/model"
	"github.com/bojand/ghz-web/service"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// SetupProjectAPI sets up the API
func SetupProjectAPI(g *echo.Group) {
	api := &ProjectAPI{}

	g.POST("/", api.create)
	g.GET("/:id", api.get)
	g.PUT("/:id", api.update)
	g.PUT("/:id", api.delete)

	testsGroup := g.Group("/:id/tests")
	SetupTestAPI(testsGroup)
}

// ProjectAPI provides the api
type ProjectAPI struct {
	dao service.ProjectService
}

func (api *ProjectAPI) create(c echo.Context) error {
	p := new(model.Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := api.dao.Create(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, p)
}

func (api *ProjectAPI) get(c echo.Context) error {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	getByID := true
	if err != nil {
		getByID = false
	}

	c.Logger().Info("Getting project: " + string(id))

	var p *model.Project

	if getByID {
		if p, err = api.dao.FindByID(uint(id)); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	} else {
		if p, err = api.dao.FindByName(idparam); gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Bad Request: "+err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) update(c echo.Context) error {
	p := new(model.Project)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	uid := uint(id)
	p.ID = uid

	if err = api.dao.Update(p); gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

func (api *ProjectAPI) delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
