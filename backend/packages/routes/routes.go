package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	// 1. Create Echo instance
	e := echo.New()

	// 2. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 3. Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API is running!")
	})

	e.POST("/login", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Login endpoint reached",
		})
	})

	// 4. Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
