package session

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/items/armor"
)

func (s *Session) exploreTreasureRoom(ctx context.Context, room *floormaze.TreasureRoom) {
	s.Printf("Вы находите комнату с сокровищами.")
	if s.makePauseFor(ctx, 3*time.Second) != nil {
		return
	}

	inf := "Среди горы хлама, вы находите редкие предметы:\n"
	for i, item := range room.Items {
		if x, ok := item.(I.Weapon); ok {
			inf += fmt.Sprintf(" %d) %s. %s\n", i+1, x.Title(), x.Description())
		} else if x, ok := item.(*armor.ArmorItem); ok {
			inf += fmt.Sprintf(" %d) %s. %s\n", i+1, x.Name(), x.Description())
		} else if x, ok := item.(I.Consumable); ok {
			inf += fmt.Sprintf(" %d) %s [x%d] (%s)\n", i+1, x.Name(), x.UsesLeft(), x.Description())
		} else if x, ok := item.(I.Usable); ok {
			inf += fmt.Sprintf(" %d) %s (%s)\n", i+1, x.Name(), x.Description())
		} else {
			inf += fmt.Sprintf(" %d) %s\n", i+1, item.Name())
		}
	}

	inf += "Каждый может забрать не более двух предметов!"
	s.Printf(inf)

	taken := make(map[int]int)

	s.receiveWithAlert(ctx, 60*time.Second, func(umsg I.UserMessage, cancel func()) {
		opt, ok := strconv.Atoi(umsg.Message())
		p := s.players.GetByID(umsg.User().ID())
		if p == nil || ok != nil {
			return
		}
		opt -= 1

		if room.Items[opt] == nil {
			return
		} else if taken[p.ID()] == 2 {
			s.Printf("%s уже взял 2 предмета!", p.FullName())
		} else {
			item := room.Items[opt]
			room.Items[opt] = nil
			taken[p.ID()] += 1

			if x, ok := item.(I.Weapon); ok {
				s.Printf("%s забирает %s!", p.FullName(), x.Title())
				p.PickWeapon(x)
			} else if x, ok := item.(*armor.ArmorItem); ok {
				s.Printf("%s надевает %s!", p.FullName(), x.Name())
				// drop other armor if has
				p.ForEachItem(func(item I.Pickable) {
					if x, ok := item.(*armor.ArmorItem); ok {
						p.DropItem(x)
					}
				})
				p.PickItem(x)
			} else {
				s.Printf("%s берёт %s!", p.FullName(), item.Name())
				p.PickItem(item)
			}

			if len(taken) == len(room.Items) {
				s.Printf("Все предметы разобрали! Продолжаем дальше!")
				cancel()
			}
		}
	}, 50*time.Second, "Осталось 10 секунд чтобы взять предметы!")

	s.Printf("Статы игроков:\n%s", s.players.Info())
}
