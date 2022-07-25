package main

import "github.com/golang-jwt/jwt"

type Auth struct {
	User User `json:"user"`
	Channel string `json:"channel"`
	jwt.StandardClaims
}
