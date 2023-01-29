package model

import "fmt"

type User struct {
	ID        int
	FirstName string
	LastName  string
}

func NewUser(id int, firstName, lastName string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
