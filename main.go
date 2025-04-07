package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	api := echo.New()
	api.Use(middleware.Logger())

	api.GET("/", func(c echo.Context) error {
		// Get IP address - different methods available
		ip := c.RealIP()                     // Gets the real IP, respecting X-Forwarded-For header
		remoteAddr := c.Request().RemoteAddr // Gets direct connection IP address

		// Get hostname from request
		host := c.Request().Host

		// Return information in the response
		return c.String(200, fmt.Sprintf("Hello from %s!\nYour IP: %s\nRemote Address: %s",
			host, ip, remoteAddr))
	})

	if err := api.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
