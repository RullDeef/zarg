package model

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"zarg/lib/model/enemy"
	enemySquad "zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player"
	"zarg/lib/model/player/squad"
	"zarg/lib/model/reorder_referee"
	"zarg/lib/model/trap"
	"zarg/lib/model/weapon"
	"zarg/lib/model/weapon/showcase"
	"zarg/lib/utils"
)

type Session struct {
	interactor Interactor
	players    *squad.PlayerSquad
	onDone     func()
	cancelFunc func()
	pauser     *utils.Pauser
}

func NewSession(i Interactor, onDone func()) *Session {
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
	if s.players.Len() > 1 && !s.determinePlayersOrder(ctx) {
		s.interactor.Printf("Что ж, если вы не можете определиться с очередностью, командной работы точно не будет. Поход отменён.")
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

	s.receiveWithAlert(ctx, 30*time.Second, func(umsg UserMessage, _ func()) {
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
		return weapon.RandomWeapon(0, 6, 2)
	})

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. "
	ask += "Выберите ваше оружие среди представленных:\n"
	ask += weaponShowcase.WeaponsInfo()
	ask += "И поторопитесь, через 30 секунд выдвигаемся!"
	s.interactor.Printf(ask)

	canceled := s.receiveWithAlert(ctx, 30*time.Second, func(umsg UserMessage, cancel func()) {
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
	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg UserMessage, cancel func()) {
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

func (s *Session) explore(ctx context.Context, fm *FloorMaze) {
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

func (s *Session) exploreEnemiesRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Вы не одни... На вас напали!")

	enemySquad := enemySquad.New(2+rand.Intn(2), func() I.Enemy {
		return enemy.Random(11, 3, func(e *enemy.Enemy) {
			p := s.players.ChooseRandomAlive()
			healthBefore := p.Health()
			p.Damage(e.AttackPower())

			if p.Alive() {
				s.interactor.Printf("%s атакует %s. (HP:%d->%d)", e.Name, p.FullName(), healthBefore, p.Health)
				if s.makePauseFor(ctx, time.Second) != nil {
					return
				}
			} else {
				s.interactor.Printf("%s атакует %s. %s убит!", e.Name, p.FullName(), p.FullName())
				if s.makePauseFor(ctx, 5*time.Second) != nil {
					return
				}
			}
		})
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

func (s *Session) makePlayerAction(ctx context.Context, p I.Player, es *enemySquad.EnemySquad) {
	s.interactor.Printf("Ходит %s.", p.FullName())

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}

	optsInfo := "Варианты действия:\n"
	opts := map[int]func(){}
	i := 1
	es.ForEachAlive(func(e I.Enemy) {
		h1, h2 := e.Health(), e.Health()-p.Weapon().Attack()
		if h2 > 0 {
			optsInfo += fmt.Sprintf("%d) Атаковать %s (HP: %d->%d)\n", i, e.Name(), h1, h2)
			opts[i] = func() {
				s.interactor.Printf("%s атакует %s", p.FullName(), e.Name())
				e.Damage(p.Weapon().Attack())
			}
		} else {
			optsInfo += fmt.Sprintf("%d) Убить %s\n", i, e.Name())
			opts[i] = func() {
				s.interactor.Printf("%s убивает %s", p.FullName(), e.Name())
				e.Damage(p.Weapon().Attack())
			}
		}
		i += 1
	})
	s.interactor.Printf(optsInfo)

	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg UserMessage, cancel func()) {
		opt, err := strconv.Atoi(umsg.Message)
		if umsg.User.ID() == p.ID() && err == nil {
			if action := opts[opt]; action != nil {
				action()
				cancel()
			}
		}
	}, 45*time.Second, "Ещё 15 секунд, чтобы сделать выбор!")

	if !canceled {
		s.interactor.Printf("%s решает пропустить ход!", p.FullName())
	}
}

func (s *Session) makeEnemyAction(ctx context.Context, e I.Enemy, es *enemySquad.EnemySquad) {
	s.interactor.Printf("Ходит %s.", e.Name)
	if s.makePauseFor(ctx, 3*time.Second) != nil {
		return
	}

	e.Attack()
}

func (s *Session) exploreTreasureRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Вы находите комнату с сокровищами.")
	s.interactor.Printf("Но пока что здесь пусто... Разработчик еще не добавил ништяков...")
}

func (s *Session) exploreTrapRoom(ctx context.Context, fm *FloorMaze) {
	s.interactor.Printf("Что-то подсказывает вам, что тут не все так безобидно, как кажется на первый взгляд.")

	if s.makePauseFor(ctx, 4*time.Second) != nil {
		return
	}

	probMap := utils.NewPropMap()

	probMap.Add(trap.New("Гигантская стрела вылетела прямо из стены!", trap.DamageRandom, 15), 4)
	probMap.Add(trap.New("Острые шипы выступают перед вашими ногами!", trap.DamageFirst, 22), 4)
	probMap.Add(trap.New("С потолка сваливается огромный камень!", trap.DamageRandom, 11), 5)
	probMap.Add(trap.New("Из темноты вылетает стая летучих мышей!", trap.DamageEveryone, 9), 5)
	probMap.Add(trap.New("Вы попадаете под душ из кислоты!", trap.DamageEveryone, 13), 2)
	probMap.Add(trap.New("Пол проваливается, и кто-то оказывается в лаве!", trap.DamageRandom, 999), 1)

	t := probMap.Choose().(*trap.Trap)
	s.interactor.Printf(t.Description())

	if s.makePauseFor(ctx, 2*time.Second) != nil {
		return
	}

	damaged := t.Activate(s.players)
	info := ""

	var killedNames []string
	for _, p := range damaged {
		if !p.Alive() {
			killedNames = append(killedNames, p.FullName())
		} else {
			info += fmt.Sprintf("%s (HP:->%d)\n", p.FullName(), p.Health())
		}
	}

	if killedNames != nil {
		if len(killedNames) == 1 {
			info += killedNames[0] + " убит!"
		} else {
			info += fmt.Sprintf("%s убиты!", strings.Join(killedNames, ", "))
		}
	}

	s.interactor.Printf(info)
}

func (s *Session) exploreRestRoom(ctx context.Context, fm *FloorMaze) {
	info := "Вы находите комнату, в которой можно перевести дух и обговорить дальнейшие планы.\n"
	info += "Голосуйте \"в путь\" за то, чтобы продолжить поход, и \"строй\" за то, чтобы изменить очередность."
	s.interactor.Printf(info)

	if s.makePauseFor(ctx, 2*time.Second) != nil {
		return
	}

	s.players.ForEachAlive(func(p I.Player) {
		p.Heal(50)
	})
	s.interactor.Printf("+50HP всем игрокам.")

	continueCounter := 0
	reorderCounter := 0
	reordering := false

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.receivePauseAware(ctx, func(umsg UserMessage) {
		p := s.players.GetByID(umsg.User.ID())
		if p == nil || reordering {
			return
		}

		msg := strings.ToLower(strings.TrimSpace(umsg.Message))
		if msg == "в путь" {
			continueCounter += 1
			if 2*continueCounter > s.players.LenAlive() {
				cancel()
			}
		} else {
			continueCounter = 0
		}

		if msg == "строй" {
			reorderCounter += 1
			if 2*reorderCounter > s.players.LenAlive() {
				reorderCounter = 0
				reordering = true
				go func(ctx context.Context) {
					if !s.determinePlayersOrder(ctx) {
						s.interactor.Printf("Очередность не изменена.")
					}
					reordering = false
				}(ctx)
			}
		} else {
			reorderCounter = 0
		}
	})

	s.interactor.Printf("Поход продолжается!")
}

// returns true if was canceled
func (s *Session) receiveWithAlert(ctx context.Context, d time.Duration, f func(umsg UserMessage, cancel func()), alertTime time.Duration, alertMsg string) bool {
	alarm := utils.AfterFunc(ctx, alertTime, s.pauser, func() {
		s.interactor.Printf(alertMsg)
	})
	defer alarm.Stop()
	return s.receiveWithTimeout(ctx, d, f)
}

// returns true if was canceled
func (s *Session) receiveWithTimeout(ctx context.Context, d time.Duration, f func(umsg UserMessage, cancel func())) bool {
	ctx, cancel := s.timeoutFor(ctx, d)
	canceled := false
	s.receivePauseAware(ctx, func(umsg UserMessage) {
		f(umsg, func() {
			canceled = true
			cancel()
		})
	})
	return canceled
}

func (s *Session) receivePauseAware(ctx context.Context, f func(UserMessage)) error {
	return s.interactor.Receive(ctx, func(umsg UserMessage) {
		if !s.pauser.IsPaused() {
			f(umsg)
		}
	})
}

func (s *Session) timeoutFor(parent context.Context, d time.Duration) (context.Context, func()) {
	ctx, cancel := context.WithCancel(parent)
	timer := utils.NewTimer(d, s.pauser)

	go func(parent context.Context) {
		defer cancel()
		<-timer.WaitCompleted(parent)
	}(parent)

	return ctx, cancel
}

func (s *Session) makePauseFor(ctx context.Context, d time.Duration) error {
	if ctx.Err() == nil {
		<-utils.NewTimer(d, s.pauser).WaitCompleted(ctx)
	}

	return ctx.Err()
}
