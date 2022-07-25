package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Logger middleware logs the information about each HTTP request.
	e.Use(middleware.Logger())

	// Recover middleware recovers from panics anywhere in the chain.
	e.Use(middleware.Recover())

	e.POST("/auth", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		channel := c.FormValue("channel")

		claims := &Auth{
			User{name, email},
			channel,
			jwt.StandardClaims{
					Issuer: "chitchat",
					ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// For production purpose better to use RS256 Signing Method instead
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte("secret"))
		if err != err {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	})

	config := middleware.JWTConfig{
		Claims: &Auth{},
		SigningKey: []byte("secret"),
		ContextKey: "token",
		TokenLookup: "query:token",
		ErrorHandler: func(err error) error {
			// Masking all jwt errors from the clients
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "Unauthorized",
			}
		},
	}

	requireAuth := middleware.JWTWithConfig(config)

	e.GET("/auth", func(c echo.Context) error {
		token := c.Get("token").(*jwt.Token)
		auth := token.Claims.(*Auth)

		return c.JSON(http.StatusOK, auth)
	}, requireAuth)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Welcome to ChitChat!",
		})
	})

	e.Logger.Fatal(e.Start(":4000"))
}