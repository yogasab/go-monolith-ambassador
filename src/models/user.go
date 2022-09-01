package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	Password     string
	IsAmbassador bool
}

func (u *User) HashPassword(plainPassword string) string {
	hashPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	return string(hashPasswordBytes)
}
