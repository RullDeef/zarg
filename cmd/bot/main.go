package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"zarg/lib/controllers"
	"zarg/lib/model"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func main() {
	vk := api.NewVK(os.Getenv("API_TOKEN"))

	// Получаем информацию о группе
	group, err := vk.GroupsGetByID(api.Params{})
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем longpoll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// Событие нового сообщения
	lp.MessageNew(makeLongpollHandler(vk))

	// Запускаем Bots Longpoll
	log.Println("Start longpoll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func makeLongpollHandler(vk *api.VK) func(context.Context, events.MessageNewObject) {
	sessionActive := false
	var chat chan model.Replica
	var out chan string

	return func(_ context.Context, obj events.MessageNewObject) {
		log.Print(obj.Message.Text)
		if strings.ToLower(obj.Message.Text) == "в поход" {
			if sessionActive {
				sendMessage(vk, obj.Message.PeerID, "Нельзя начать еще один поход!")
			} else {
				sessionActive = true
				chat = make(chan model.Replica)
				out = make(chan string)

				go messageSender(vk, obj.Message.PeerID, out)
				go func() {
					controllers.BeginSession(chat, out)
					sessionActive = false
				}()
			}
		} else if sessionActive {
			userName := getUserName(vk, obj.Message.FromID)
			chat <- model.NewReplica(obj.Message.FromID, userName, obj.Message.Text)
		}
	}
}

func messageSender(vk *api.VK, peerID int, in chan string) {
	for {
		msg, ok := <-in
		if !ok {
			break
		}

		sendMessage(vk, peerID, msg)
	}
}

func sendMessage(vk *api.VK, peerID int, message string) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(peerID)
	b.Message(message)
	b.RandomID(0)

	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserName(vk *api.VK, userID int) string {
	users, err := vk.UsersGet(map[string]interface{}{
		"user_ids": []int{userID},
		"fields":   []string{},
	})
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s %s", users[0].FirstName, users[0].LastName)
}
