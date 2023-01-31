package session

import (
	"context"
	"strconv"
	"strings"
	"time"
	"zarg/lib/model"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player"
	"zarg/lib/model/player/squad"
	"zarg/lib/model/reorder_referee"
	"zarg/lib/model/weapon"
	"zarg/lib/model/weapon/showcase"
	"zarg/lib/utils"
)

type Session struct {
	interactor model.Interactor
	players    *squad.PlayerSquad
	onDone     func()
	cancelFunc func()
	pauser     *utils.Pauser
}

func NewSession(i model.Interactor, onDone func()) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Session{
		interactor: i,
		players:    squad.NewPlayerSquad(),
		onDone:     onDone,
		cancelFunc: cancel,
		pauser:     utils.NewPauser(),
	}

	go s.startup(ctx)
	return s
}

func (s *Session) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
}

func (s *Session) startup(ctx context.Context) {
	defer s.shutdown()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go s.commandProcessor(ctx)
	if !s.gatherPlayers(ctx) {
		return
	}

	s.pickWeapons(ctx)
	if ctx.Err() != nil {
		return
	}

	if s.players.Len() > 1 && !s.determinePlayersOrder(ctx) {
		if ctx.Err() != nil {
			return
		}
		s.interactor.Printf("Что ж, если вы не можете определиться с очередностью, командной работы точно не будет. Поход отменён.")
		return
	}

	if ctx.Err() != nil {
		return
	}

	// generate first floor maze
	floorMaze := model.NewFloorMaze("Подземелье", 1, 3)
	s.explore(ctx, floorMaze)

}

func (s *Session) shutdown() {
	if s.pauser.IsPaused() {
		s.pauser.Resume()
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
	s.interactor.Printf("Игровая сессия завершена.")
	s.onDone()
}

func (s *Session) commandProcessor(ctx context.Context) {
	s.interactor.Receive(ctx, func(umsg model.UserMessage) {
		if s.pauser.IsPaused() {
			if umsg.Message == "/прод" {
				s.interactor.Printf("Игра продолжается!")
				s.pauser.Resume()
			}
		} else {
			if umsg.Message == "/пауза" {
				s.pauser.Pause()
				s.interactor.Printf("Игра на паузе! чтобы продолжить напишите \"/прод\".")
			} else if umsg.Message == "/стат" {
				s.interactor.Printf("Статы игроков:\n%s", s.players.Info())
			}
		}
	})
}

func (s *Session) gatherPlayers(ctx context.Context) bool {
	s.interactor.Printf("Начинается сбор людей и нелюдей для похода в данж!\nЧтобы участвовать, напиши \"Я\".")

	s.receiveWithAlert(ctx, 30*time.Second, func(umsg model.UserMessage, _ func()) {
		msg := strings.ToLower(strings.TrimSpace(umsg.Message))
		if msg == "я" {
			if p := s.players.GetByID(umsg.User.ID()); p != nil {
				s.interactor.Printf("%s уже в списке!", umsg.User.FullName())
			} else {
				s.players.Add(player.NewPlayer(umsg.User))
				s.interactor.Printf("%s участвует в походе!", umsg.User.FullName())
			}
		} else if msg == "не я" {
			if p := s.players.RemoveByID(umsg.User.ID()); p != nil {
				s.interactor.Printf("%s вычеркнут(-а) из списка.", p.FullName())
			}
		}
	}, 20*time.Second, "10 секунд до окончания сбора!")

	if s.players.Len() == 0 {
		s.interactor.Printf("Cбор окончен! В поход не идёт никто.")
		return false
	} else if s.players.Len() == 1 {
		s.interactor.Printf("Одного смельчака недостаточно, чтобы покорить данж! Поход отменён.")
		return false
	}

	res := "Сбор окончен! В поход собрались:\n"
	res += s.players.ListString()
	s.interactor.Printf(res)
	return true
}

func (s *Session) pickWeapons(ctx context.Context) {
	// generate weapons on start
	const totalWeapons = 6
	nChosen := 0 // amount of players that already picked weapon

	weaponShowcase := showcase.NewWeaponShowcase(totalWeapons, func() I.Weapon {
		return weapon.RandomWeapon(6, 2)
	})

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. "
	ask += "Выберите ваше оружие среди представленных:\n"
	ask += weaponShowcase.WeaponsInfo()
	ask += "И поторопитесь, через 30 секунд выдвигаемся!"
	s.interactor.Printf(ask)

	canceled := s.receiveWithAlert(ctx, 30*time.Second, func(umsg model.UserMessage, cancel func()) {
		opt, err := strconv.Atoi(umsg.Message)
		if err != nil || opt < 1 || opt > totalWeapons {
			return
		}

		// find player that chooses
		if p := s.players.GetByID(umsg.User.ID()); p != nil {
			hadNoWeapon := !weaponShowcase.HasMadePick(p)
			ok, w, other := weaponShowcase.TryPick(p, opt-1)

			if ok {
				if hadNoWeapon {
					nChosen += 1
				}
				s.interactor.Printf("%s выбирает %s.", p.FullName(), w.Title())
			} else {
				s.interactor.Printf("%s уже выбрал %s!", w.Title(), other.FullName())
			}
		}

		if nChosen == s.players.Len() {
			weaponShowcase.ConfirmPick()
			s.interactor.Printf("Все выбрали по оружию, отправляемся в данж!")
			cancel()
			return
		}
	}, 20*time.Second, "Осталось 10 секунд!")

	if !canceled {
		weaponShowcase.ConfirmPick()
		s.interactor.Printf("Выдвигаемся! А кто не успел схватиться за оружие будет сражаться кулаками!")
		s.players.ForEach(func(p I.Player) {
			if p.Weapon() == nil {
				p.PickWeapon(weapon.FistsWeapon(5, 1))
			}
		})
	}
}

func (s *Session) determinePlayersOrder(ctx context.Context) bool {
	s.interactor.Printf("Определите очередность ходов для пошагового режима. Первый игрок должен написать \"я\", остальные - \"потом я\".")

	referee := reorder_referee.New(s.players)
	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg model.UserMessage, cancel func()) {
		msg := strings.ToLower(strings.TrimSpace(umsg.Message))
		if (msg == "я" && referee.VoteStarter(umsg.User.ID())) || (msg == "потом я" && referee.VoteNext(umsg.User.ID())) {
			if referee.Completed() {
				referee.Apply()
				s.interactor.Printf("Очередность установлена!\n%s\n", s.players.OrderingString())
				cancel()
			} else {
				s.interactor.Printf("Очередность: %s\n", referee.OrderingInfo())
			}
		}
	}, 45*time.Second, "У вас еще 15 секунд, чтобы определиться!")

	return canceled
}

func (s *Session) explore(ctx context.Context, fm *model.FloorMaze) {
	if ctx.Err() != nil {
		return
	}

	roomProbGen := utils.NewPropMap()
	// roomProbGen.Add("enemy", 4)
	roomProbGen.Add("treasure", 3)
	// roomProbGen.Add("trap", 2)
	// roomProbGen.Add("rest", 2)

	//for i := 0; i < fm.roomsCount; i++ {
	for {
		if s.makePauseFor(ctx, 3*time.Second) != nil {
			return
		}

		switch roomProbGen.Choose().(string) {
		case "enemy":
			s.exploreEnemiesRoom(ctx, fm)
		case "treasure":
			s.exploreTreasureRoom(ctx, fm)
		case "trap":
			s.exploreTrapRoom(ctx, fm)
		case "rest":
			s.exploreRestRoom(ctx, fm)
		}

		if s.players.LenAlive() == 0 {
			return
		}
	}
}
