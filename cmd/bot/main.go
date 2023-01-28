package main

import (
	"context"
	"log"
	"os"
	"strings"
	"zarg/lib/controllers"
	"zarg/lib/model"

	"github.com/SevereCloud/vksdk/v2/api"
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
	interactors := make(map[int]*controllers.VKInteractor)
	sessions := make(map[int]*model.Session)

	return func(_ context.Context, obj events.MessageNewObject) {
		chatID := obj.Message.PeerID
		msg := obj.Message.Text
		log.Printf("%d: %s\n", chatID, msg)

		// check for interactor for this chat
		interactor, ok := interactors[chatID]
		if !ok {
			interactor = controllers.NewVKInteractor(vk, chatID)
			interactors[chatID] = interactor
		}

		if strings.ToLower(strings.TrimSpace(msg)) == "в поход" {
			if sessions[chatID] != nil {
				interactor.Printf("Нельзя начать еще один поход!")
			} else {
				sessions[chatID] = model.NewSession(interactor, func() {
					sessions[chatID] = nil
				})
			}
		} else {
			interactor.SendMessage(obj)
		}
	}
}
