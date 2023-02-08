package vk

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"zarg/lib/model"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/user"
	"zarg/lib/utils"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

type VKInteractor struct {
	vk        *api.VK
	pub       *utils.Publisher
	groupID   int
	chatUsers map[int]*user.User
	fakeUsers map[string]*user.User
	lock      sync.Mutex
}

func NewVKInteractor(vk *api.VK, groupID int) *VKInteractor {
	return &VKInteractor{
		vk:        vk,
		pub:       utils.NewPublisher(),
		groupID:   groupID,
		chatUsers: map[int]*user.User{},
		fakeUsers: map[string]*user.User{},
		lock:      sync.Mutex{},
	}
}

func (i *VKInteractor) Close() {
	log.Printf("VKInteractor #%p closed", i)
	i.pub.Close()
}

func (i *VKInteractor) Printf(format string, a ...any) {
	i.lock.Lock()
	defer i.lock.Unlock()

	msg := fmt.Sprintf(format, a...)
	b := params.NewMessagesSendBuilder()
	b.PeerID(i.groupID)
	b.Message(msg)
	b.RandomID(0)

	_, err := i.vk.MessagesSend(b.Params)
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("VKInteractor sent: \"%s\"", msg)
	}
}

func (i *VKInteractor) SendMessage(m events.MessageNewObject) {
	userID := m.Message.FromID
	msg := m.Message.Text

	// check message is fake
	re := regexp.MustCompile(`^([^:]+?)\s*:\s*([^:]+?)$`)
	if match := re.FindStringSubmatch(msg); match != nil {
		msg, name := match[1], match[2]
		log.Printf("fake message detected from \"%s\": \"%s\"", name, msg)
		i.sendFakeMessage(name, msg)
	} else {
		i.sendChatMessage(userID, msg)
	}
}

func (i *VKInteractor) sendChatMessage(id int, msg string) {
	i.lock.Lock()
	u := i.lockedGetChatUser(id)
	defer i.pub.Publish(model.NewUserMessage(u, msg))
	defer i.lock.Unlock()
}

func (i *VKInteractor) sendFakeMessage(userName string, msg string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	u := i.lockedGetFakeUser(userName)
	i.pub.Publish(model.NewUserMessage(u, msg))
}

func (i *VKInteractor) Receive(ctx context.Context, f func(I.UserMessage)) error {
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

func (i *VKInteractor) lockedGetChatUser(id int) *user.User {
	if u := i.chatUsers[id]; u != nil {
		return u
	}

	// load new chat user
	b := params.NewUsersGetBuilder()
	b.UserIDs([]string{strconv.Itoa(id)})

	users, err := i.vk.UsersGet(b.Params)
	if err != nil {
		log.Fatal(err)
	}

	u := user.New(id, users[0].FirstName, users[0].LastName)
	i.chatUsers[id] = u
	return u
}

func (i *VKInteractor) lockedGetFakeUser(name string) *user.User {
	if u := i.fakeUsers[name]; u != nil {
		return u
	}

	// create new fake user
	u := user.New(rand.Int(), name, "")
	i.fakeUsers[name] = u
	return u
}
