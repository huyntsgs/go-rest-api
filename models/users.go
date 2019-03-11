package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type (
	User struct {
		Id        int64     `json:"userId"`
		UserName  string    `json:"userName"`
		Password  string    `json:"password"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

// Validate checks user data is valid or not for register.
func (u *User) Validate() bool {
	if govalidator.IsNull(u.UserName) || govalidator.IsNull(u.Password) || !govalidator.IsEmail(u.Email) {
		return false
	}
	return true
}

// Validate checks user data is valid or not for login.
func (u *User) ValidateLogin() bool {
	if govalidator.IsNull(u.Password) || !govalidator.IsEmail(u.Email) {
		return false
	}
	return true
}
