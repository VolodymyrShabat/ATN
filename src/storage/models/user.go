package models

import "fmt"

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Name     string `json:"name,omitempty" validate:"required"`
	City     string `json:"city,omitempty" validate:"required"`
	Age      int    `json:"age,omitempty" validate:"required"`
}

func (u *User) String() string {
	return fmt.Sprintf("User id: %d\n login: %s\n email: %s\n name: %s\n city: %s\n age: %d\n", u.Id, u.Login, u.Email, u.Name, u.City, u.Age)
}
