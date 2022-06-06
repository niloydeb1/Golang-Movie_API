package api

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/niloydeb1/Golang-Movie_API/api/v1"
	"github.com/swaggo/echo-swagger"
	"net/http"
)

// Routes base router
func Routes(e *echo.Echo) {
	// Index Page
	e.GET("/", index)

	// Health Page
	e.GET("/health", health)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	v1.Router(e.Group("/api/v1"))
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "This is Golang Movie API Service")
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "I am live!")
}