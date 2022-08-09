package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var jwtSecret = "secret"

type tokenResponse struct {
	Token string `json:"token"`
}

func TestCreateAuth(t *testing.T) {
	authForm := url.Values{}

	userName := "Jon Snow"
	userEmail := "jon@labstack.com"
	channel := "general"

	authForm.Add("name", userName)
	authForm.Add("email", userEmail)
	authForm.Add("channel", channel)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewAuthHandler(jwtSecret)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		tokenResponse := &tokenResponse{}
		json.Unmarshal(rec.Body.Bytes(), tokenResponse)

		if assert.NotEmpty(t, tokenResponse.Token) {
			token, err := jwt.ParseWithClaims(tokenResponse.Token, &Auth{}, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(jwtSecret), nil
			})

			assert.Equal(t, err, nil)
			assert.Equal(t, token.Valid, true)

			if claims, ok := token.Claims.(*Auth); ok && token.Valid {

				assert.Equal(t, claims.User.Name, userName)
				assert.Equal(t, claims.Channel, channel)

				assert.Contains(t, claims.User.Avatar, "https://www.gravatar.com/avatar/")
			} else {
				t.Error("Token is not valid")
			}
		}
	}
}

func TestCreateAuthWithoutEmail(t *testing.T) {
	authForm := url.Values{}

	userName := "Jon Snow"
	channel := "general"

	authForm.Add("name", userName)
	authForm.Add("channel", channel)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewAuthHandler(jwtSecret)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		tokenResponse := &tokenResponse{}
		json.Unmarshal(rec.Body.Bytes(), tokenResponse)

		if assert.NotEmpty(t, tokenResponse.Token) {
			token, err := jwt.ParseWithClaims(tokenResponse.Token, &Auth{}, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(jwtSecret), nil
			})

			assert.Equal(t, err, nil)
			assert.Equal(t, token.Valid, true)

			if claims, ok := token.Claims.(*Auth); ok && token.Valid {

				assert.Equal(t, claims.User.Name, userName)
				assert.Equal(t, claims.User.Email, "")
				assert.Equal(t, claims.Channel, channel)

				assert.Equal(t, claims.User.Avatar, "")
			} else {
				t.Error("Token is not valid")
			}
		}
	}
}

func TestCreateAuthNameValidation(t *testing.T) {
	authForm := url.Values{}

	userEmail := "jon@labstack.com"
	channel := "general"

	authForm.Add("email", userEmail)
	authForm.Add("channel", channel)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewAuthHandler(jwtSecret)

	err := h.Create(c)

	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "code=422, message=Name can't be blank")
	}
}
func TestCreateAuthChannelValidation(t *testing.T) {
	authForm := url.Values{}

	userName := "Jon Snow"
	userEmail := "jon@labstack.com"

	authForm.Add("name", userName)
	authForm.Add("email", userEmail)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewAuthHandler(jwtSecret)

	err := h.Create(c)

	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "code=422, message=Channel can't be blank")
	}
}
func TestCreateAuthEmailValidation(t *testing.T) {
	authForm := url.Values{}

	userName := "Jon Snow"
	userEmail := "not_valid_email"
	channel := "general"

	authForm.Add("name", userName)
	authForm.Add("email", userEmail)
	authForm.Add("channel", channel)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewAuthHandler(jwtSecret)

	err := h.Create(c)

	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "code=422, message=Email is not invalid")
	}
}

func TestGetAuth(t *testing.T) {
	authForm := url.Values{}

	userName := "Jon Snow"
	userEmail := "jon@labstack.com"
	channel := "general"

	authForm.Add("name", userName)
	authForm.Add("email", userEmail)
	authForm.Add("channel", channel)

	h := NewAuthHandler(jwtSecret)

	e := echo.New()

	e.GET("/auth", h.Get, h.Require)

	req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(authForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		tokenResponse := &tokenResponse{}
		json.Unmarshal(rec.Body.Bytes(), tokenResponse)

		reqGet := httptest.NewRequest(http.MethodGet, "/auth?token="+tokenResponse.Token, nil)
		resGet := httptest.NewRecorder()

		e.ServeHTTP(resGet, reqGet)

		assert.Equal(t, http.StatusOK, resGet.Code)

		authResponse := &Auth{}
		json.Unmarshal(resGet.Body.Bytes(), authResponse)

		assert.Equal(t, authResponse.User.Name, userName)
		assert.Equal(t, authResponse.Channel, channel)

		assert.Contains(t, authResponse.User.Avatar, "https://www.gravatar.com/avatar/")
	}
}

func TestGetAuthWithoutToken(t *testing.T) {
	h := NewAuthHandler(jwtSecret)

	e := echo.New()

	e.GET("/auth", h.Get, h.Require)

	reqGet := httptest.NewRequest(http.MethodGet, "/auth", nil)
	resGet := httptest.NewRecorder()

	e.ServeHTTP(resGet, reqGet)

	assert.Equal(t, http.StatusUnauthorized, resGet.Code)
	assert.Equal(t, `{"message":"Unauthorized"}`+"\n", resGet.Body.String())
}

func TestGetAuthInvalidToken(t *testing.T) {
	h := NewAuthHandler(jwtSecret)

	e := echo.New()

	e.GET("/auth", h.Get, h.Require)

	reqGet := httptest.NewRequest(http.MethodGet, "/auth?token=invalid", nil)
	resGet := httptest.NewRecorder()

	e.ServeHTTP(resGet, reqGet)

	assert.Equal(t, http.StatusUnauthorized, resGet.Code)
	assert.Equal(t, `{"message":"Unauthorized"}`+"\n", resGet.Body.String())
}
