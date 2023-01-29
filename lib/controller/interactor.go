package controller

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"zarg/lib/model"
	"zarg/lib/utils"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

type VKInteractor struct {
	vk        *api.VK
	pub       *utils.Publisher
	groupID   int
	lock      sync.Mutex
	fakeUsers map[int]string
}

func NewVKInteractor(vk *api.VK, groupID int) *VKInteractor {
	return &VKInteractor{
		vk:        vk,
		pub:       utils.NewPublisher(),
		groupID:   groupID,
		lock:      sync.Mutex{},
		fakeUsers: map[int]string{},
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

	// check fake users first
	if name := i.fakeUsers[userID]; name != "" {
		return name
	}

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

	// check message is fake
	terms := strings.Split(msg, " ")
	if len(terms) > 0 && terms[len(terms)-1][0] == ':' {
		name := terms[len(terms)-1][1:]
		msg = strings.Join(terms[:len(terms)-1], " ")
		log.Printf("fake message detected from \"%s\": \"%s\"", name, msg)
		i.sendFakeMessage(name, msg)
	} else {
		i.pub.Publish(model.NewUserMessage(userID, msg))
	}
}

func (i *VKInteractor) sendFakeMessage(userName string, msg string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	userID := 1
	for id, name := range i.fakeUsers {
		if name == userName {
			userID = id
			break
		}
		if userID < id+1 {
			userID = id + 1
		}
	}
	i.fakeUsers[userID] = userName
	i.pub.Publish(model.NewUserMessage(userID, msg))
}

func (i *VKInteractor) Receive(ctx context.Context, f func(model.UserMessage)) error {
	s := utils.NewSubscriber(i.pub)
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
