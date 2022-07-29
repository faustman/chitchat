package main

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const AuthExpiresInHours = 72

// Auth represents current JWT auth.
type Auth struct {
	// Current User.
	User *User `json:"user"`

	// Current channel.
	Channel string `json:"channel"`

	// Rest JWT headers.
	jwt.StandardClaims
}

// authHandler resposible for all auth stuff.
type authHandler struct {
	// SigningKey for JWT signing process.
	SigningKey []byte
	// Require midddleware parse and validate jwt token.
	Require echo.MiddlewareFunc
}

// NewAuthHandler prepare authHandler
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

// Create /auth generate token for specific user.
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
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Email is not invalid")
		}
	}

	claims := &Auth{
		NewUser(name, email),
		channel,
		jwt.StandardClaims{
				Issuer: "chitchat",
				ExpiresAt: time.Now().Add(time.Hour * AuthExpiresInHours).Unix(),
		},
	}

	// For production purpose better to use RS256 Signing Method instead
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(a.SigningKey)
	if err != err {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": t,
	})
}

// Get /auth will return valid auth object.
func (a authHandler) Get(c echo.Context) error {
	auth := ExtactAuth(c)

	return c.JSON(http.StatusOK, auth)
}

// ExtactAuth from echo context.
func ExtactAuth(c echo.Context) *Auth {
	token := c.Get("token").(*jwt.Token)

	return token.Claims.(*Auth)
}