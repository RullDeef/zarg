package model

import "time"

type UserMessage struct {
	UserID  int
	Message string
}

type Interactor interface {
	Close()

	Printf(fmt string, args ...any)

	GetUserName(userID int) string

	// gets a message from chat. Close channel on timeout
	ReceiveFor(d time.Duration) chan UserMessage

	// perform action on timeout when nobody response to game
	// SetTimeoutAction(d time.Duration, action func())
}

func NewUserMessage(userID int, msg string) UserMessage {
	return UserMessage{
		UserID:  userID,
		Message: msg,
	}
}
