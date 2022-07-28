package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Logger middleware logs the information about each HTTP request.
	e.Use(middleware.Logger())

	// Recover middleware recovers from panics anywhere in the chain.
	e.Use(middleware.Recover())

	authHandler := NewAuthHandler(os.Getenv("JWT_SECRET"))

	e.POST("/auth", authHandler.Create)

	e.GET("/auth", authHandler.Get, authHandler.Require)

	fmt.Println("Connecting to NATS..")

	stream, err := NewStream()
	if err != nil {
		e.Logger.Fatal(err)
	}

	consumersHub := NewConsumersHub()
	go consumersHub.run()

	channelHandler := NewChannelHandler(stream, consumersHub)
	e.GET("/channel", channelHandler.Listen, authHandler.Require)
	e.GET("/messages", channelHandler.GetMessages, authHandler.Require)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Welcome to ChitChat!",
		})
	})

	e.Logger.Fatal(e.Start(":4000"))
}