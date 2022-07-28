package main

import (
	"crypto/md5"
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Email string `json:"-"`
	Avatar string `json:"avatar,omitempty"`
}

func NewUser(name, email string) *User {
	user := User{
		Name: name,
		Email: email,
	}

	if len(email) > 0 {
		user.Avatar = generateGravatar(email)
	}

	return &user
}

func generateGravatar(email string) string {
	hashString := []byte(email)

	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=128", md5.Sum(hashString))
}