package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	IsAmbassador bool   `json:"-"`
}

func (u *User) HashPassword(plainPassword string) string {
	hashPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	return string(hashPasswordBytes)
}

func (u *User) ComparePassword(hashedPassword string, plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return false
	}
	return true
}
