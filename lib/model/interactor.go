package model

import (
	"context"
)

type UserMessage struct {
	User    *User
	Message string
}

type Interactor interface {
	Printf(fmt string, args ...any)

	// gets a messages from chat.
	Receive(ctx context.Context, f func(UserMessage)) error

	// perform action on timeout when nobody response to game
	// SetTimeoutAction(d time.Duration, action func())
}

func NewUserMessage(user *User, msg string) UserMessage {
	return UserMessage{
		User:    user,
		Message: msg,
	}
}
