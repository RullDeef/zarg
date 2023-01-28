package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
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
	lock    sync.Mutex
}

func NewVKInteractor(vk *api.VK, groupID int) *VKInteractor {
	return &VKInteractor{
		vk:      vk,
		pub:     pubsub.NewPublisher(),
		groupID: groupID,
		lock:    sync.Mutex{},
	}
}

func (i *VKInteractor) Close() {
	log.Printf("VKInteractor #%p closed", i)
	i.pub.Close()
}

func (i *VKInteractor) Printf(format string, a ...any) {
	i.lock.Lock()
	defer i.lock.Unlock()

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
	i.lock.Lock()
	defer i.lock.Unlock()

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

func (i *VKInteractor) Receive(ctx context.Context, f func(model.UserMessage)) error {
	s := pubsub.NewSubscriber(i.pub)
	defer s.Unsubscribe()

	for {
		select {
		case msg := <-s.Receive():
			f(msg.(model.UserMessage))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
