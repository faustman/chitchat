package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	// Logger middleware logs the information about each HTTP request.
	e.Use(middleware.Logger())

	// Recover middleware recovers from panics anywhere in the chain.
	e.Use(middleware.Recover())

	authHandler := NewAuthHandler(os.Getenv("JWT_SECRET"))

	// Create new JWT Token
	e.POST("/auth", authHandler.Create)
	// Get auth and validate JWT Token
	e.GET("/auth", authHandler.Get, authHandler.Require)

	e.Logger.Info("Connecting to NATS..")

	stream, err := NewStream()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Run consumer hub
	consumersHub := NewConsumersHub()
	go consumersHub.run()

	channelHandler := NewChannelHandler(stream, consumersHub)
	e.GET("/channel", channelHandler.Listen, authHandler.Require)
	e.GET("/messages", channelHandler.GetMessages, authHandler.Require)
	e.GET("/users", channelHandler.GetUsers, authHandler.Require)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Welcome to ChitChat!",
		})
	})

	// Graceful Shutdown stuff from https://echo.labstack.com/cookbook/graceful-shutdown/

	// Start server
	go func() {
		if err := e.Start(":4000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	e.Logger.Info("Shutdown..")
	consumersHub.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
