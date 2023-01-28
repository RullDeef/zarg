package model

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	interactor Interactor
	players    *PlayerList
	onDone     func()
	cancelFunc func()
}

func NewSession(i Interactor, onDone func()) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Session{
		interactor: i,
		players:    NewPlayerList(),
		onDone:     onDone,
		cancelFunc: cancel,
	}

	go s.startup(ctx)
	return s
}

func (s *Session) Stop() {
	s.shutdown()
}

func (s *Session) startup(ctx context.Context) {
	defer s.shutdown()

	if !s.gatherPlayers(ctx) {
		return
	}

	s.pickWeapons(ctx)
	if s.players.Len() > 1 && !s.determinePlayersOrder(ctx) {
		return
	}
}

func (s *Session) shutdown() {
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
	s.interactor.Printf("Игровая сессия завершена.")
	s.onDone()
}

func (s *Session) gatherPlayers(ctx context.Context) bool {
	s.interactor.Printf("Начинается сбор людей и нелюдей для похода в данж!\nЧтобы участвовать, напиши \"Я\".")

	alarm := time.AfterFunc(time.Duration(10*time.Second), func() {
		s.interactor.Printf("5 секунд до окончания сбора!")
	})

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	s.interactor.Receive(ctx, func(umsg UserMessage) {
		if strings.ToLower(strings.TrimSpace(umsg.Message)) == "я" {
			if p := s.players.GetByID(umsg.UserID); p != nil {
				s.interactor.Printf("%s уже в списке!", p.Name())
			} else {
				userName := s.interactor.GetUserName(umsg.UserID)
				s.players.Add(NewPlayer(umsg.UserID, userName))
				s.interactor.Printf("%s участвует в походе!", userName)
			}
		} else if strings.ToLower(strings.TrimSpace(umsg.Message)) == "не я" {
			if p := s.players.RemoveByID(umsg.UserID); p != nil {
				s.interactor.Printf("%s вычеркнут(-а) из списка.", p.Name())
			}
		}
	})
	cancel()
	alarm.Stop()

	if s.players.Len() == 0 {
		s.interactor.Printf("Cбор окончен! В поход не идёт никто.")
		return false
	}
	// else if s.players.Len() == 1 {
	// 	s.interactor.Printf("Одного смельчака недостаточно, чтобы покорить данж! Поход отменён.")
	// 	return false
	// }

	res := "Сбор окончен! В поход собрались:\n"
	res += s.players.ListString()
	s.interactor.Printf(res)
	return true
}

func (s *Session) pickWeapons(ctx context.Context) {
	// generate weapons on start
	const totalWeapons = 6
	nChosen := 0 // amount of players that already picked weapon

	weaponShowcase := NewWeaponShowcase(totalWeapons, func() *Weapon {
		return RandomWeapon(0, 6, 2)
	})

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. "
	ask += "Выберите ваше оружие среди представленных:\n"
	ask += weaponShowcase.WeaponsInfo()
	ask += "И поторопитесь, через 15 секунд выдвигаемся!"
	s.interactor.Printf(ask)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	err := s.interactor.Receive(ctx, func(umsg UserMessage) {
		opt, err := strconv.Atoi(umsg.Message)
		if err != nil || opt < 1 || opt > totalWeapons {
			return
		}

		// find player that chooses
		if p := s.players.GetByID(umsg.UserID); p != nil {
			hadNoWeapon := !weaponShowcase.HasMadePick(p)
			ok, w, other := weaponShowcase.TryPick(p, opt-1)

			if ok {
				if hadNoWeapon {
					nChosen += 1
				}
				s.interactor.Printf("%s выбирает %s.", p.Name(), w.SummaryShort())
			} else {
				s.interactor.Printf("%s уже выбрал %s!", w.SummaryShort(), other.Name())
			}
		}

		if nChosen == s.players.Len() {
			weaponShowcase.ConfirmPick()
			s.interactor.Printf("Все выбрали по оружию, отправляемся в данж!")
			cancel()
			return

		}
	})

	if errors.Is(err, context.DeadlineExceeded) {
		weaponShowcase.ConfirmPick()
		s.interactor.Printf("Выдвигаемся! А кто не успел схватиться за оружие будет сражаться кулаками!")
		s.players.Foreach(func(_ int, p *Player) {
			if p.Weapon == nil {
				p.Weapon = FistsWeapon(5, 1)
			}
		})
	}
}

func (s *Session) determinePlayersOrder(ctx context.Context) bool {
	s.interactor.Printf("Определите очередность ходов для пошагового режима. Первый игрок должен написать \"я\", остальные - \"потом я\".")

	var order []int

	alarm := time.AfterFunc(time.Duration(45*time.Second), func() {
		s.interactor.Printf("Если не поторопитесь, игра завершится так и не начавшись!")
	})

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	err := s.interactor.Receive(ctx, func(umsg UserMessage) {
		if s.players.GetByID(umsg.UserID) == nil {
			return // skip non-players messages
		}

		msg := strings.ToLower(strings.TrimSpace(umsg.Message))
		if msg == "я" {
			order = append(([]int)(nil), umsg.UserID)
		} else if msg == "потом я" {
			// remove id from list
			for i, id := range order {
				if id == umsg.UserID {
					order = append(order[:i], order[i+1:]...)
					break
				}
			}

			order = append(order, umsg.UserID)
		} else {
			return
		}

		if len(order) == s.players.Len() {
			alarm.Stop()
			s.players.SetOrdering(order)
			s.interactor.Printf("Очередность установлена!\n%s\n", s.players.OrderingString())
			cancel()
			return
		}

		log.Printf("order: %v\n", order)
		s.interactor.Printf("Очередность: %s\n", s.players.PhantomOrderingString(order))
	})
	alarm.Stop()

	if errors.Is(err, context.Canceled) {
		return true
	}

	s.interactor.Printf("Что ж, если вы не можете определиться с очередностью, командной работы точно не будет. Поход отменён.")
	return false
}
