package session

import (
	"context"
	"strconv"
	"strings"
	"time"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player"
	"zarg/lib/model/player/squad"
	"zarg/lib/model/reorder_referee"
	"zarg/lib/model/weapon"
	"zarg/lib/model/weapon/showcase"
	"zarg/lib/service/logs"
	"zarg/lib/utils"
)

type Session struct {
	logger     *logs.Logger
	interactor I.Interactor
	players    *squad.PlayerSquad
	onDone     func()
	cancelFunc func()
	pauser     *utils.Pauser
}

func NewSession(i I.Interactor, onDone func()) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Session{
		logger:     logs.New(),
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
	defer func() {
		if err := recover(); err != nil {
			s.logger.Printf("recovered from panic: %s", err)
			s.Printf("К сожалению, произошли непредвиденные обстоятельства, и подземелье завалило. Больше героев никто не видел...")
		}
	}()

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
		s.Printf("Что ж, если вы не можете определиться с очередностью, командной работы точно не будет. Поход отменён.")
		return
	}

	if ctx.Err() != nil {
		return
	}

	// generate first floor maze
	floorMaze := floormaze.GenFloorMaze("Подземелье")
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
	s.Printf("Игровая сессия завершена.")
	s.logger.Close()
	s.onDone()
}

func (s *Session) commandProcessor(ctx context.Context) {
	s.interactor.Receive(ctx, func(umsg I.UserMessage) {
		if s.pauser.IsPaused() {
			if umsg.Message() == "/прод" {
				s.Printf("Игра продолжается!")
				s.pauser.Resume()
			}
		} else {
			if umsg.Message() == "/пауза" {
				s.pauser.Pause()
				s.Printf("Игра на паузе! чтобы продолжить напишите \"/прод\".")
			} else if umsg.Message() == "/стат" {
				s.Printf("Статы игроков:\n%s", s.players.Info())
			}
		}
	})
}

func (s *Session) gatherPlayers(ctx context.Context) bool {
	s.Printf("Начинается сбор людей и нелюдей для похода в данж!\nЧтобы участвовать, напиши \"Я\".")

	s.receiveWithAlert(ctx, 30*time.Second, func(umsg I.UserMessage, _ func()) {
		msg := strings.ToLower(strings.TrimSpace(umsg.Message()))
		if msg == "я" {
			if p := s.players.GetByID(umsg.User().ID()); p != nil {
				s.Printf("%s уже в списке!", umsg.User().FullName())
			} else {
				s.players.Add(player.NewPlayer(umsg.User()))
				s.Printf("%s участвует в походе!", umsg.User().FullName())
			}
		} else if msg == "не я" {
			if p := s.players.RemoveByID(umsg.User().ID()); p != nil {
				s.Printf("%s вычеркнут(-а) из списка.", p.FullName())
			}
		}
	}, 20*time.Second, "10 секунд до окончания сбора!")

	if s.players.Len() == 0 {
		s.Printf("Cбор окончен! В поход не идёт никто.")
		return false
	} else if s.players.Len() == 1 {
		s.Printf("Одного смельчака недостаточно, чтобы покорить данж! Поход отменён.")
		return false
	}

	res := "Сбор окончен! В поход собрались:\n"
	res += s.players.ListString()
	s.Printf(res)
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
	s.Printf(ask)

	canceled := s.receiveWithAlert(ctx, 30*time.Second, func(umsg I.UserMessage, cancel func()) {
		opt, err := strconv.Atoi(umsg.Message())
		if err != nil || opt < 1 || opt > totalWeapons {
			return
		}

		// find player that chooses
		if p := s.players.GetByID(umsg.User().ID()); p != nil {
			hadNoWeapon := !weaponShowcase.HasMadePick(p)
			ok, w, other := weaponShowcase.TryPick(p, opt-1)

			if ok {
				if hadNoWeapon {
					nChosen += 1
				}
				s.Printf("%s выбирает %s.", p.FullName(), w.Title())
			} else {
				s.Printf("%s уже выбрал %s!", w.Title(), other.FullName())
			}
		}

		if nChosen == s.players.Len() {
			weaponShowcase.ConfirmPick()
			s.Printf("Все выбрали по оружию, отправляемся в данж!")
			cancel()
			return
		}
	}, 20*time.Second, "Осталось 10 секунд!")

	if !canceled {
		weaponShowcase.ConfirmPick()
		s.Printf("Выдвигаемся! А кто не успел схватиться за оружие будет сражаться кулаками!")
		s.players.ForEach(func(p I.Player) {
			if p.Weapon() == nil {
				p.PickWeapon(weapon.FistsWeapon(5, 1))
			}
		})
	}
}

func (s *Session) determinePlayersOrder(ctx context.Context) bool {
	s.Printf("Определите очередность ходов для пошагового режима. Первый игрок должен написать \"я\", остальные - \"потом я\".")

	referee := reorder_referee.New(s.players)
	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg I.UserMessage, cancel func()) {
		msg := strings.ToLower(strings.TrimSpace(umsg.Message()))
		id := umsg.User().ID()
		if (msg == "я" && referee.VoteStarter(id)) || (msg == "потом я" && referee.VoteNext(id)) {
			if referee.Completed() {
				referee.Apply()
				s.Printf("Очередность установлена!\n%s\n", s.players.OrderingString())
				cancel()
			} else {
				s.Printf("Очередность: %s\n", referee.OrderingInfo())
			}
		}
	}, 45*time.Second, "У вас еще 15 секунд, чтобы определиться!")

	return canceled
}

func (s *Session) explore(ctx context.Context, fm *floormaze.FloorMaze) {
	for s.players.LenAlive() > 0 {
		if s.makePauseFor(ctx, 3*time.Second) != nil {
			return
		}

		room := fm.NextRoom()
		switch room.(type) {
		case floormaze.EnemyRoom:
			s.exploreEnemiesRoom(ctx, room.(*floormaze.EnemyRoom))
		case floormaze.TrapRoom:
			s.exploreTrapRoom(ctx, room.(*floormaze.TrapRoom))
		case floormaze.TreasureRoom:
			s.exploreTreasureRoom(ctx, room.(*floormaze.TreasureRoom))
		case floormaze.RestRoom:
			s.exploreRestRoom(ctx, room.(*floormaze.RestRoom))
		case floormaze.BossRoom:
			s.exploreBossRoom(ctx, room.(*floormaze.BossRoom))
			return
		}
	}
}
