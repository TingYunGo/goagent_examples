package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Handler
func hello(c echo.Context) error {
	fmt.Println("on request")
	return c.String(http.StatusOK, "Hello, World!")
}

//go:noinline
func Test(p func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route) {
	p("/test", func(c echo.Context) error {
		fmt.Println("on request /test")
		return c.String(http.StatusOK, "/test Hello, World!")

	})
}
func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	Test(e.GET)
	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
