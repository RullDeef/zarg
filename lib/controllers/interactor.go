package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"zarg/lib/model"
	"zarg/lib/pubsub"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

type VKInteractor struct {
	vk      *api.VK
	pub     *pubsub.Publisher
	groupID int
}

func NewVKInteractor(vk *api.VK, groupID int) *VKInteractor {
	return &VKInteractor{
		vk:      vk,
		pub:     pubsub.NewPublisher(),
		groupID: groupID,
	}
}

func (i *VKInteractor) Close() {
	i.pub.Close()
}

func (i *VKInteractor) Printf(format string, a ...any) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(i.groupID)
	b.Message(fmt.Sprintf(format, a...))
	b.RandomID(0)

	_, err := i.vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *VKInteractor) GetUserName(userID int) string {
	b := params.NewUsersGetBuilder()
	b.UserIDs([]string{strconv.Itoa(userID)})

	users, err := i.vk.UsersGet(b.Params)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s %s", users[0].FirstName, users[0].LastName)
}

func (i *VKInteractor) SendMessage(m events.MessageNewObject) {
	userID := m.Message.FromID
	msg := m.Message.Text
	i.pub.Publish(model.NewUserMessage(userID, msg))
}

func (i *VKInteractor) ReceiveFor(d time.Duration) chan model.UserMessage {
	m := make(chan model.UserMessage)

	go func() {
		s := pubsub.NewSubscriber(i.pub)
		defer s.Unsubscribe()
		defer close(m)

		delay := time.After(d)
	outer:
		for {
			select {
			case msg := <-s.Receive():
				m <- msg.(model.UserMessage)
			case <-delay:
				break outer
			}
		}
	}()

	return m
}
