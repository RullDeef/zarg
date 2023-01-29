package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"zarg/lib/utils"
)

type Session struct {
	interactor Interactor
	players    *PlayerList
	onDone     func()
	cancelFunc func()
	pauser     *utils.Pauser
}

func NewSession(i Interactor, onDone func()) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Session{
		interactor: i,
		players:    NewPlayerList(),
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
	if s.players.Len() > 1 && !s.determinePlayersOrder(ctx) {
		return
	}

	// generate first floor maze
	floorMaze := NewFloorMaze("Подземелье", 1, 3)
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
	s.interactor.Receive(ctx, func(umsg UserMessage) {
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

	alarm := time.AfterFunc(time.Duration(20*time.Second), func() {
		s.interactor.Printf("10 секунд до окончания сбора!")
	})

	ctx, cancel := s.timeoutFor(ctx, 30*time.Second)
	s.receivePauseAware(ctx, func(umsg UserMessage) {
		log.Printf("receivePauseAware: %s", umsg.Message)

		if strings.ToLower(strings.TrimSpace(umsg.Message)) == "я" {
			if p := s.players.GetByID(umsg.User.ID); p != nil {
				s.interactor.Printf("%s уже в списке!", p.user.FullName())
			} else {
				userName := umsg.User.FullName()
				s.players.Add(NewPlayer(umsg.User))
				s.interactor.Printf("%s участвует в походе!", userName)
			}
		} else if strings.ToLower(strings.TrimSpace(umsg.Message)) == "не я" {
			if p := s.players.RemoveByID(umsg.User.ID); p != nil {
				s.interactor.Printf("%s вычеркнут(-а) из списка.", p.user.FullName())
			}
		}
	})
	cancel()
	alarm.Stop()

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

	weaponShowcase := NewWeaponShowcase(totalWeapons, func() *Weapon {
		return RandomWeapon(0, 6, 2)
	})

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. "
	ask += "Выберите ваше оружие среди представленных:\n"
	ask += weaponShowcase.WeaponsInfo()
	ask += "И поторопитесь, через 30 секунд выдвигаемся!"
	s.interactor.Printf(ask)

	ctx, cancel := s.timeoutFor(ctx, 30*time.Second)
	err := s.receivePauseAware(ctx, func(umsg UserMessage) {
		opt, err := strconv.Atoi(umsg.Message)
		if err != nil || opt < 1 || opt > totalWeapons {
			return
		}

		// find player that chooses
		if p := s.players.GetByID(umsg.User.ID); p != nil {
			hadNoWeapon := !weaponShowcase.HasMadePick(p)
			ok, w, other := weaponShowcase.TryPick(p, opt-1)

			if ok {
				if hadNoWeapon {
					nChosen += 1
				}
				s.interactor.Printf("%s выбирает %s.", p.user.FullName(), w.SummaryShort())
			} else {
				s.interactor.Printf("%s уже выбрал %s!", w.SummaryShort(), other.user.FullName())
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

	var order []*User

	alarm := time.AfterFunc(time.Duration(45*time.Second), func() {
		s.interactor.Printf("Если не поторопитесь, игра завершится так и не начавшись!")
	})

	ctx, cancel := s.timeoutFor(ctx, time.Minute)
	err := s.receivePauseAware(ctx, func(umsg UserMessage) {
		if s.players.GetByID(umsg.User.ID) == nil {
			return // skip non-players messages
		}

		msg := strings.ToLower(strings.TrimSpace(umsg.Message))
		if msg == "я" {
			order = append(([]*User)(nil), umsg.User)
		} else if msg == "потом я" {
			// remove id from list
			for i, u := range order {
				if u == umsg.User {
					order = append(order[:i], order[i+1:]...)
					break
				}
			}

			order = append(order, umsg.User)
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

func (s *Session) explore(ctx context.Context, fm *FloorMaze) {
	if ctx.Err() != nil {
		return
	}

	roomProbGen := utils.NewPropMap()
	roomProbGen.Add("enemy", 4)
	// roomProbGen.Add("treasure", 3)
	// roomProbGen.Add("trap", 2)
	// roomProbGen.Add("rest", 2)

	for i := 0; i < fm.roomsCount; i++ {
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

		break
	}
}

func (s *Session) exploreEnemiesRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Вы не одни... На вас напали!")

	enemySquad := NewEnemySquad(2+rand.Intn(2), func() *Enemy {
		tier := fm.tierMin + rand.Intn(fm.tierMax-fm.tierMin+1)
		return RandomEnemy(tier, 9, 2)
	})

	// show battle overall info
	battleInfo := fmt.Sprintf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), enemySquad.CompactInfo())
	s.interactor.Printf(battleInfo)

	for enemySquad.LenAlive() > 0 && s.players.LenAlive() > 0 {
		if s.makePauseFor(ctx, time.Second) != nil {
			return
		}

		// make turn generator
		turnGen := utils.NewPropMap()

		turnGen.Add("players", s.players.LenAlive())
		turnGen.Add("enemies", enemySquad.LenAlive())

		switch turnGen.Choose().(string) {
		case "players":
			// show battle overall info
			battleInfo := fmt.Sprintf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), enemySquad.CompactInfo())
			s.interactor.Printf(battleInfo)
			p := s.players.ChooseNext()
			s.makePlayerAction(ctx, p, enemySquad)
		case "enemies":
			e := enemySquad.ChooseNext()
			s.makeEnemyAction(ctx, e, enemySquad)
		default:
			log.Fatal("must never happen!")
		}
	}

	if enemySquad.LenAlive() == 0 {
		s.interactor.Printf("Битва завершена. Все враги повержены!")
	} else {
		s.interactor.Printf("Битва завершена. Все игроки мертвы!")
	}
}

func (s *Session) makePlayerAction(ctx context.Context, p *Player, es *EnemySquad) {
	s.interactor.Printf("Ходит %s.", p.User().FullName())

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}

	optsInfo := "Варианты действия:\n"
	opts := map[int]func(){}
	i := 1
	es.ForeachAlive(func(_ int, e *Enemy) {
		h1, h2 := e.Health, e.Health-p.Weapon.Attack
		if h2 > 0 {
			optsInfo += fmt.Sprintf("%d) Атаковать %s (HP: %d->%d)\n", i, e.Name, h1, h2)
			opts[i] = func() {
				s.interactor.Printf("%s атакует %s", p.User().FullName(), e.Name)
				e.MakeDamage(p.Weapon.Attack)
			}
		} else {
			optsInfo += fmt.Sprintf("%d) Убить %s\n", i, e.Name)
			opts[i] = func() {
				s.interactor.Printf("%s убивает %s", p.User().FullName(), e.Name)
				e.MakeDamage(p.Weapon.Attack)
			}
		}
		i += 1
	})

	s.interactor.Printf(optsInfo)

	ctx, cancel := s.timeoutFor(ctx, time.Minute)
	err := s.receivePauseAware(ctx, func(umsg UserMessage) {
		if umsg.User != p.User() {
			return
		}

		opt, err := strconv.Atoi(umsg.Message)
		if err != nil {
			return
		}

		action := opts[opt]
		if action != nil {
			action()
			cancel()
		}
	})

	if errors.Is(err, context.DeadlineExceeded) {
		s.interactor.Printf("%s решает пропустить ход!", p.User().FullName())
	}
	cancel()
}

func (s *Session) makeEnemyAction(ctx context.Context, e *Enemy, es *EnemySquad) {
	s.interactor.Printf("Ходит %s.", e.Name)

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}

	// choose random player and attack him
	p := s.players.ChooseRandomAlive()
	p.MakeDamage(e.Attack)

	if p.Health == 0 {
		s.interactor.Printf("%s атакует %s. %s убит!", e.Name, p.User().FullName(), p.User().FullName())
	} else {
		s.interactor.Printf("%s атакует %s.", e.Name, p.User().FullName())
	}

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}
}

func (s *Session) exploreTreasureRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Вы находите комнату с сокровищами.")
}

func (s *Session) exploreTrapRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Что-то подсказывает вам, что тут не все так безобидно, как кажется на первый взгляд.")
}

func (s *Session) exploreRestRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Вы находите комнату, в которой можно перевести дух и обговорить дальнейшие планы.")
}

func (s *Session) receivePauseAware(ctx context.Context, f func(UserMessage)) error {
	return s.interactor.Receive(ctx, func(umsg UserMessage) {
		if !s.pauser.IsPaused() {
			f(umsg)
		}
	})
}

func (s *Session) timeoutFor(ctx context.Context, d time.Duration) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	timer := utils.NewTimer(d, s.pauser)

	go func() {
		defer cancel()
		<-timer.WaitCompleted(ctx)
	}()

	return ctx, cancel
}

func (s *Session) makePauseFor(ctx context.Context, d time.Duration) error {
	if ctx.Err() == nil {
		<-utils.NewTimer(d, s.pauser).WaitCompleted(ctx)
	}

	return ctx.Err()
}
