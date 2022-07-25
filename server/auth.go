package main

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const AUTH_EXPIRES_IN_HOURS = 72

type Auth struct {
	User *User `json:"user"`
	Channel string `json:"channel"`
	jwt.StandardClaims
}

type authHandler struct {
	SigningKey []byte
	Require echo.MiddlewareFunc
}

func NewAuthHandler(jwtSecret string) *authHandler {
	SigningKey := []byte(jwtSecret)

	config := middleware.JWTConfig{
		Claims: &Auth{},
		SigningKey: SigningKey,
		ContextKey: "token",
		TokenLookup: "query:token", // getting token for url, eg. /auth?token=abc
		ErrorHandler: func(err error) error {
			// Masking all jwt errors from the clients
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "Unauthorized",
			}
		},
	}

	return &authHandler{
		SigningKey: SigningKey,
		Require: middleware.JWTWithConfig(config),
	}
}

func (a authHandler) Create(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	channel := c.FormValue("channel")

	if len(name) == 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Name can't be blank")
	}

	if len(channel) == 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Channel can't be blank")
	}

	if len(email) > 0 {
		// TODO: extact to helper
		_, err := mail.ParseAddress(email)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Email invalid")
		}
	}

	claims := &Auth{
		NewUser(name, email),
		channel,
		jwt.StandardClaims{
				Issuer: "chitchat",
				ExpiresAt: time.Now().Add(time.Hour * AUTH_EXPIRES_IN_HOURS).Unix(),
		},
	}

	// For production purpose better to use RS256 Signing Method instead
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(a.SigningKey)
	if err != err {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (a authHandler) Get(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	auth := token.Claims.(*Auth)

	return c.JSON(http.StatusOK, auth)
}