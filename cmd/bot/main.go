package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"zarg/lib/controller"
	"zarg/lib/model"
	"zarg/lib/service"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

var sessionBroker *service.SessionBroker

func main() {
	rand.Seed(time.Now().Unix())

	sessionBroker = service.NewSessionBroker()

	initVK()
}

func initVK() {
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
	return func(_ context.Context, obj events.MessageNewObject) {
		chatID := obj.Message.PeerID
		msg := obj.Message.Text
		log.Printf("%d: %s\n", chatID, msg)

		i := sessionBroker.Interactor(chatID, func() model.Interactor {
			return controller.NewVKInteractor(vk, chatID)
		}).(*controller.VKInteractor)

		cmd := strings.ToLower(strings.TrimSpace(msg))

		if cmd == "в поход" {
			if !sessionBroker.AddSession(chatID, func() {
				log.Printf("session for chatID=%d ended", chatID)
			}) {
				i.Printf("Нельзя начать еще один поход!")
			}
		} else if cmd == "/стоп" {
			if s := sessionBroker.Session(chatID); s != nil {
				s.Stop()
				log.Print("session canceled")
			}
		} else if cmd == "/пинг" {
			i.Printf("понг")
		} else {
			i.SendMessage(obj)
		}
	}
}
