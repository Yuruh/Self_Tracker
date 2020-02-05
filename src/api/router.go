package api

import (
	"github.com/Yuruh/Self_Tracker/src/spotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func messageHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	})
}

func RunHttpServer()  {
	// Echo instance
	app := echo.New()
	app.HideBanner = true

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	// Routes
	app.GET("/", hello)
	app.GET("/spotify", func(c echo.Context) error {
		return c.String(http.StatusOK, spotify.BuildAuthUri())
	})

	// Start server
	app.Logger.Fatal(app.Start(":8090"))

}

// Handler
func hello(c echo.Context) error {
	panic("lol")
	return c.String(http.StatusOK, "Hello, World!")
}