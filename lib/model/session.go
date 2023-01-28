package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	interactor Interactor
	players    *PlayerList
	onDone     func()
}

func NewSession(i Interactor, onDone func()) *Session {
	s := &Session{
		interactor: i,
		players:    NewPlayerList(),
		onDone:     onDone,
	}

	go s.startup()
	return s
}

func (s *Session) startup() {
	defer s.shutdown()

	if s.gatherPlayers() {
		s.pickWeapons()
	}
}

func (s *Session) shutdown() {
	s.onDone()
}

func (s *Session) gatherPlayers() bool {
	s.interactor.Printf("Начинается сбор людей и нелюдей для похода в данж!\nЧтобы участвовать, напиши \"Я\".")

	time.AfterFunc(time.Duration(10*time.Second), func() {
		s.interactor.Printf("5 секунд до окончания сбора!")
	})

	for umsg := range s.interactor.ReceiveFor(time.Duration(15 * time.Second)) {
		if strings.ToLower(strings.TrimSpace(umsg.Message)) == "я" {
			// if player already joined
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
	}

	if s.players == nil {
		s.interactor.Printf("Cбор окончен! В поход не идёт никто.")
		return false
	}
	res := "Сбор окончен! В поход собрались:\n"
	s.players.Foreach(func(i int, p *Player) {
		res += fmt.Sprintf("  %d. %s\n", i+1, p.Name())
	})
	s.interactor.Printf(res)
	return true
}

func (s *Session) pickWeapons() {
	// generate weapons on start
	const totalWeapons = 6
	weapons := RandomWeapons(totalWeapons, 0, 6, 2)
	chosen := make([]bool, totalWeapons)
	nChosen := 0 // amount of players that already picked weapon

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. Выберите ваше оружие среди представленных:\n"
	for i, weapon := range weapons {
		ask += fmt.Sprintf("  %d) %s\n", i+1, weapon.Summary())
	}

	ask += "И поторопитесь, через 15 секунд выдвигаемся!"
	s.interactor.Printf(ask)

	for umsg := range s.interactor.ReceiveFor(time.Duration(15 * time.Second)) {
		if opt, err := strconv.ParseInt(umsg.Message, 10, 32); err == nil {
			if opt < 1 || opt > totalWeapons {
				continue
			}

			// find player that chooses
			if p := s.players.GetByID(umsg.UserID); p != nil {
				if !chosen[opt-1] {
					if p.Weapon == nil {
						nChosen += 1
					}
					p.Weapon = weapons[opt-1]
					s.interactor.Printf("%s выбирает %s.", p.Name(), weapons[opt-1].SummaryShort())
				} else {
					s.interactor.Printf("%s уже выбрал другой игрок!", weapons[opt-1].SummaryShort())
				}
			}

			if nChosen == s.players.Len() {
				s.interactor.Printf("Все выбрали по оружию, отправляемся в данж!")
				break
			}
		}
	}

	s.interactor.Printf("Выдвигаемся! А кто не успел схватиться за оружие будет сражаться кулаками!")
	s.players.Foreach(func(_ int, p *Player) {
		if p.Weapon == nil {
			p.Weapon = FistsWeapon(5, 1)
		}
	})
}
