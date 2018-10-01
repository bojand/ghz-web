package api

import (
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// TODO add tests

func getPageParam(c echo.Context) uint {
	pageparam := c.QueryParam("page")
	page := uint(0)
	if pageparam != "" {
		pageNum, err := strconv.Atoi(pageparam)
		if err == nil {
			if pageNum < 0 {
				pageNum = pageNum * -1
			}

			page = uint(pageNum)
		}
	}

	return page
}

func getSortAndOrder(c echo.Context) (bool, string, string) {
	doSort := false
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")
	if sort != "" {
		sort = strings.ToLower(sort)
		if sort == "id" || sort == "name" {
			doSort = true
		}

		if doSort {
			order = strings.ToLower(order)
			if order != "asc" && order != "desc" {
				order = "asc"
			}
		}
	}

	return doSort, sort, order
}
