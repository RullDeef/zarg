package model

import (
	I "zarg/lib/model/interfaces"
)

type UserMessage struct {
	user    I.User
	message string
}

func NewUserMessage(user I.User, msg string) UserMessage {
	return UserMessage{
		user:    user,
		message: msg,
	}
}

// UserMessage interface implementation
func (um UserMessage) User() I.User {
	return um.user
}

// UserMessage interface implementation
func (um UserMessage) Message() string {
	return um.message
}
