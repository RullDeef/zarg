package session

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"zarg/lib/model"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/items/armor"
	"zarg/lib/model/weapon"
	"zarg/lib/utils"
)

func (s *Session) exploreTreasureRoom(ctx context.Context, fm *model.FloorMaze) {
	s.interactor.Printf("Вы находите комнату с сокровищами.")

	if s.makePauseFor(ctx, 3*time.Second) != nil {
		return
	}

	inf := "Среди горы хлама, вы находите редкие предметы:\n"

	probMap := utils.NewPropMap()
	for i := 0; i < 6; i++ {
		probMap.Add(weapon.RandomWeapon(10, 4), 1)
	}
	for i := 0; i < 3; i++ {
		probMap.Add(armor.Random(), 1)
	}

	var items []any
	for i := 0; i < 6; i++ {
		item := probMap.Choose()
		items = append(items, item)
		probMap.Add(item, -1)

		if x, ok := item.(I.Weapon); ok {
			inf += fmt.Sprintf(" %d) %s. %s\n", i+1, x.Title(), x.Description())
		} else if x, ok := item.(*armor.ArmorItem); ok {
			inf += fmt.Sprintf(" %d) %s. %s\n", i+1, x.Name(), x.Description())
		}
	}

	inf += "Каждый может забрать не более двух предметов!"
	s.interactor.Printf(inf)

	taken := make(map[int]int)

	s.receiveWithAlert(ctx, 30*time.Second, func(umsg model.UserMessage, cancel func()) {
		opt, ok := strconv.Atoi(umsg.Message)
		p := s.players.GetByID(umsg.User.ID())
		if p == nil || ok != nil {
			return
		}
		opt -= 1

		if items[opt] == nil {
			return
		} else if taken[p.ID()] == 2 {
			s.interactor.Printf("%s уже взял 2 предмета!", p.FullName())
		} else {
			item := items[opt]
			items[opt] = nil
			taken[p.ID()] += 1

			if x, ok := item.(I.Weapon); ok {
				s.interactor.Printf("%s забирает %s!", p.FullName(), x.Title())
				p.PickWeapon(x)
			} else if x, ok := item.(*armor.ArmorItem); ok {
				s.interactor.Printf("%s надевает %s!", p.FullName(), x.Name())
				p.PickItem(x)
			}
		}
	}, 20*time.Second, "Осталось 10 секунд чтобы взять предметы!")

	s.interactor.Printf("Статы игроков:")
	s.interactor.Printf(s.players.Info())
}
