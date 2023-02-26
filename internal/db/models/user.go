package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Model reg/auth users
type User struct {
	UUID              string `json:"-"`
	Login             string `json:"login"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

// Constructor
func NewUser() *User {
	return &User{}
}

// Encrypt password
func (u *User) EncryptPassword() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		log.Printf("error in function EncryptPassword: %s\n", err)
		return err
	}
	u.EncryptedPassword = string(pass)
	return nil
}

// Compare password
func (u *User) ComparePassword() bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(u.Password)) == nil
}

// Clear password field
func (u *User) Sanitize() {
	u.Password = ""
}
