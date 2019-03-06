package main

import (
	"fmt"
	"net/http"

	"github.com/atecce/devoted/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	db := db.NewDatabase()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"GET /:name":              "db.Get",
			"PUT /:name?value=:value": "db.Set",
			"DELETE /:name":           "db.Delete",
			"GET /count?value=:value": "db.Count",
			"POST /begin":             "db.Begin",
			"POST /rollback":          "db.Rollback",
			"POST /commit":            "db.Commit",
		})
	})

	e.GET("/:name", func(c echo.Context) error {
		if val := db.Get(c.Param("name")); val == nil {
			return c.String(http.StatusNotFound, "not found")
		} else {
			return c.String(http.StatusOK, *val)
		}
	})

	e.PUT("/:name", func(c echo.Context) error {
		db.Set(c.Param("name"), c.QueryParam("value"))
		return c.String(http.StatusOK, "OK")
	})

	e.DELETE("/:name", func(c echo.Context) error {
		db.Delete(c.Param("name"))
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/count", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%d", db.Count(c.QueryParam("value"))))
	})

	e.POST("/begin", func(c echo.Context) error {
		db.Begin()
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/rollback", func(c echo.Context) error {
		db.Rollback()
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/commit", func(c echo.Context) error {
		db.Commit()
		return c.String(http.StatusOK, "OK")
	})

	e.StartTLS(":443", "/etc/devoted/server.crt", "/etc/devoted/server.key")
}
