package main

import (
	"crypto/md5"
	"fmt"
)

type User struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"-"`
	Avatar string `json:"avatar,omitempty"`
}

// NewUser will generate new users based on name and email.
// Id generating by name + email.
// If email present, Avatar will be fullfiled with Gravatar url.
func NewUser(name, email string) *User {
	id := md5.Sum([]byte(name + email))

	user := User{
		Id:    fmt.Sprintf("%x", id),
		Name:  name,
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
