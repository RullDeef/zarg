package domain

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// // if sender is nil - message is internal (someone connected or left the chat)
// type ChatMessage struct {
// 	Sender *Profile
// 	Body   string
// 	Time   time.Time
// }

// type Chat interface {
// 	Connect(user *Profile) error
// 	Disconnect(user *Profile) error

// 	Send(context.Context, ChatMessage) error
// 	ReceiverChan(context.Context) <-chan ChatMessage
// }

// func NewUserConnectedMessage(user *Profile) ChatMessage {
// 	return ChatMessage{
// 		Sender: nil,
// 		Body:   fmt.Sprintf("user %s connected to chat", user.Nickname),
// 		Time:   time.Now(),
// 	}
// }
