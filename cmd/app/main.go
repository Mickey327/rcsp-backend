package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://client:3000", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/hello", hello)
	e.Logger.Fatal(e.Start(":8080"))
}

type Answer struct {
	Code    int
	Message string
}

func hello(c echo.Context) error {
	a := &Answer{
		http.StatusOK,
		"hello",
	}
	return c.JSON(a.Code, a)
}
