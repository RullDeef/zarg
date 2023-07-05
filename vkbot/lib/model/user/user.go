package user

import "fmt"

type User struct {
	id        int
	firstName string
	lastName  string
}

func New(id int, firstName, lastName string) *User {
	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
	}
}

func (u User) ID() int {
	return u.id
}

func (u User) FirstName() string {
	return u.firstName
}

func (u User) LastName() string {
	return u.lastName
}

func (u User) FullName() string {
	if u.firstName == "" {
		return u.lastName
	} else if u.lastName == "" {
		return u.firstName
	} else {
		return fmt.Sprintf("%s %s", u.firstName, u.lastName)
	}
}
